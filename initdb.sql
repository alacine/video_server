-- create database video_server;
-- create user 'video'@'%' identified by 'videoyes';
-- grant all privileges on video_server.* to 'video'@'%' identified by 'videoyes' with grant option;
set global time_zone = '+8:00';
set time_zone = '+8:00';
flush privileges;
alter database video_server default character set utf8mb4;

drop table if exists comments;
drop table if exists sessions;
drop table if exists users;
drop table if exists video_del_rec;
drop table if exists video_info;

create table comments (
	id varchar(64) not null,
	video_id varchar(64),
	author_id int(10),
	content text,
	post_time datetime default current_timestamp, primary key(id)
);

create table sessions (
	session_id tinytext not null,
	TTL tinytext,
    uid int
);
alter table sessions add primary key (session_id(60));

create table users (
	id int unsigned not null auto_increment,
	name varchar(40),
	pwd text not null,
	unique key (name),
	primary key (id)
);

create table video_del_rec (
	-- video_id varchar(64) not null,
    video_id int unsigned not null,
	primary key (video_id)
);

create table video_info (
	id int unsigned not null auto_increment,
	-- id varchar(64) not null,
	author_id int(10),
	title text,
    description text,
	create_time datetime default current_timestamp,
	primary key (id)
);
