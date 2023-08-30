create table segments
(
	segment_id int primary key auto_increment,
    segment varchar(255) not null
);

create table reference
(
	ref_id int primary key auto_increment,
    user_id int not null,
    id_segment int not null,
    delete_date datetime,
    foreign key (segment) references segments (segment_id)
)
