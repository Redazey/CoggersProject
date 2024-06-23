create table roles
(
    id   serial primary key,
    name text not null
);

create table users
(
    id         serial primary key,
    name       text not null,
    birthdate  date,
    photourl   text,
    email      text not null,
    password   text not null,
    roleId     integer not null 
        references roles(id),
    push       boolean
);

create table servers
(
    Id        serial primary key,
    Ip        text not null,
	Name      text not null,
	Version   text not null,
	MaxOnline integer not null,
	Online    integer not null
);
