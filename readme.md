# Сервис динамического сегментирования пользователей
Сервис для сегментирования пользователей.

Используемые технологии:

- MySQL(в качестве базы данных)
- Docker
- gorilla/mux(веб библиотека)
- github.com/go-sql-driver/mysql(драйвер для работы с MySQL)
  
Для запуска следует в местах подключения к базе данных поменять значения на индивидуальные.

## Примеры

Примеры запросов:
  - [Добавление сегментов](#segment_add)
  - [Удаление сегмента](#segment_delete)
  - [Добавление пользователя в сегмент](#add_user_in_segment)
  - [Получение сегментов пользователя](#show_segments_of_user)

### Добавление сегмента <a name="segment_add"></a>
Пример запроса:

```curl
curl -XPOST http://localhost:8000/add_segment -H 'Content-Type:application/json' -d {\"segment_name\":\"AUDIO\"}
```

Пример ответа:

```json
{
  "answer": "success"
}
```
### Удаление сегмента <a name="segment_delete"></a>
Пример запроса:
```curl
curl -XPOST http://localhost:8000/delete_segment -H 'Content-Type:application/json' -d {\"segment_name\":\"AUDIO\"}
```

Пример ответа:

```json
{
  "answer": "success"
}
```
### Добавление пользователя в сегмент <a name="add_user_in_segment"></a>
Пример запроса:

```curl
curl -XPOST http://localhost:8000/add_user_in_segment -H 'Content-Type:application/json' -d {\"user\":4,\"segment_add\":[\"AUDIO\"],\"segment_delete\":[\"VIDEO\"]}
```

Пример ответа:

```json
{
  "answer": "success"
}
```
Стоит обратить внимание, что если попытаться добавить пользователю сегмент, который не был добавлен в список сегментов до этого, то сегмент пользователю не добавится, и в ответе будет предупреждение:

```curl
curl -XPOST http://localhost:8000/add_user_in_segment -H 'Content-Type:application/json' -d {\"user\":4,\"segment_add\":[\"AUDIO\",\"NEVER_ADDED_SEGMENT\"],\"segment_delete\":[\"AUDIO\"]}
```
Ответ:
```json
{
    "warning":"NEVER_ADDED_SEGMENT is not in segment list"
}
{
    "answer": "success"
}
```

Так же реализованна возможность временного добавления добавления сегмента пользователю:

```
curl -XPOST http://localhost:8000/add_user_in_segment -H 'Content-Type:application/json' -d {\"user\":4,\"segment_add\":[\"AUDIO\"],\"segment_delete\":[\"VIDEO\"],\"date_to_delete\":\"2023-08-29\"}
```
Ответ:
```json
{
    "warning":"NEVER_ADDED_SEGMENT is not in segment list"
}
{
    "answer": "success"
}
```

### Получение сегментов пользователя <a name="show_segments_of_user"></a>

Пример запроса:

```curl
curl -XGET http://localhost:8000/show_users_segments -H 'Content-Type:application/json' -d {\"user\":12}
```

Ответ:

```json
{
    "segments":["VIDEO","AUDIO"]
}
```
 
 ## Вопросы и проблемы с которыми столкнулся в ходе разработки

1. Как поступать когда при в списке сегментов, которые нужно добавить пользователю, находится сегмент, который ранее не был добавлен в список сегментов?
   >Решил что что сегмент пользователю добавлять не стоит, и отправлять предупреждение в формате JSON.
2. Каким образом реализовать удаление сегмента у пользователя, который был добавлен временно?
   >Решил что удаление будет производиться в момент вызова запроса получения сегментов пользователя. Будет производится проверка по полю с датой удаления сегмента.
3. Так же возникли проблемы с контейнеризацией пректа, так как сам только недавно начал изучать Docker. Так что для легкости копирования проекта решил сделать весь код пректа в одном файле.