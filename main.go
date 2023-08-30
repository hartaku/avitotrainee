package main

import (
	"net/http"

	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"
)

type User struct {
	User_id int `json:"user"`
}

type Segment struct {
	Segment_name string `json:"segment_name"`
}

type Reference struct {
	User_id            int      `json:"user"`
	Segments_to_add    []string `json:"segment_add"`
	Segments_to_delete []string `json:"segment_delete"`
	Date_to_delete     string   `json:"date_to_delete"`
}

var Dbconnection string = "root:12345678@tcp(127.0.0.1:3306)/dinamic_segment"

func Add(w http.ResponseWriter, r *http.Request) {
	DB, err := sql.Open("mysql", Dbconnection)
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	body, _ := ioutil.ReadAll(r.Body)
	var users_segments Reference
	json.Unmarshal(body, &users_segments)
	for _, v := range users_segments.Segments_to_add {
		var exist int
		statementText := fmt.Sprintf(`SELECT count(*) from segments where segment="%s"`, v) //Проверяем, существует ли данный сегмент в списке сегментов
		statement, err := DB.Query(statementText)
		if err != nil {
			panic(err)
		}
		for statement.Next() {
			err = statement.Scan(&exist)
		}
		//если сегмент есть в списке сегментов, то мы можем добавить его пользователю
		if exist > 0 {
			if users_segments.Date_to_delete == "" {
				statementText := fmt.Sprintf(`insert into reference(user_id,id_segment) values(%d,(SELECT segment_id from segments where segment="%s"))`, users_segments.User_id, v)
				DB.Exec(statementText)
			} else {
				statementText := fmt.Sprintf(`insert into reference(user_id,id_segment,delete_date) values(%d,(SELECT segment_id from segments where segment="%s"),"%s")`, users_segments.User_id, v, users_segments.Date_to_delete)
				DB.Exec(statementText)
			}
			//если сегмента нет в списке сегментов, то мы не можем добавить его пользователю
		} else {
			json.NewEncoder(w).Encode(map[string]string{"warning": fmt.Sprintf("%s is not in segment list", v)})
		}
	}
	for _, v := range users_segments.Segments_to_delete {
		statementText := fmt.Sprintf(`DELETE FROM reference where id_segment= (SELECT segment_id from segments where segment="%s") and user_id=%d`, v, users_segments.User_id)
		log.Println(statementText)
		DB.Exec(statementText)
	}
	json.NewEncoder(w).Encode(map[string]string{"answer": "success"})
	log.Println(users_segments)
}

func Addseg(w http.ResponseWriter, r *http.Request) {
	DB, err := sql.Open("mysql", Dbconnection)
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var segment Segment
	json.Unmarshal(body, &segment)
	log.Print(segment)
	if segment.Segment_name != "" {
		DB.Exec(fmt.Sprintf(`INSERT INTO segments(segment) VALUES('%s')`, segment.Segment_name))
	}
	json.NewEncoder(w).Encode(map[string]string{"answer": "success"})
}

func Deleteseg(w http.ResponseWriter, r *http.Request) {
	DB, err := sql.Open("mysql", Dbconnection)
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var segment Segment
	json.Unmarshal(body, &segment)
	DB.Exec(fmt.Sprintf(`DELETE FROM reference WHERE id_segment=(select segment_id from segments where segment='%s')`, segment.Segment_name))
	DB.Exec(fmt.Sprintf(`DELETE FROM segments WHERE segment='%s'`, segment.Segment_name))
	json.NewEncoder(w).Encode(map[string]string{"answer": "success"})
}

func Show(w http.ResponseWriter, r *http.Request) {
	DB, err := sql.Open("mysql", Dbconnection)
	if err != nil {
		panic(err)
	}
	//Производим удаление записей сегментов пользователей,которые были добавлены временно
	DB.Exec(`delete from reference where delete_date<now()`)
	body, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(body, &user)
	var segments_of_user []string
	statement, err := DB.Query(`select distinct s.segment from reference r join segments s on r.id_segment=s.segment_id where user_id=?`, user.User_id)
	if err != nil {
		panic(err)
	}
	for statement.Next() {
		var temp_segment string
		err = statement.Scan(&temp_segment)
		if err != nil {
			panic(err)
		}
		segments_of_user = append(segments_of_user, temp_segment)
	}
	mp := map[string]interface{}{"segments": segments_of_user}
	json.NewEncoder(w).Encode(mp)
}

func main() {
	route := mux.NewRouter()
	route.HandleFunc("/add_user_in_segment", Add).Methods("POST")
	route.HandleFunc("/delete_segment", Deleteseg).Methods("POST")
	route.HandleFunc("/add_segment", Addseg).Methods("POST")
	route.HandleFunc("/show_users_segments", Show).Methods("GET")
	http.ListenAndServe(":8000", route)
}
