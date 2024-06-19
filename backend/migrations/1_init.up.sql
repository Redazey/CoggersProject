create table users
(
    id         serial
        primary key,
    name       text not null,
    birthdate  date,
    photourl   text,
    email      text not null,
    password   text not null,
    roleId     int not null,
    push       boolean,
);

create table roles
(
    id   serial
        primary key,
    name text not null
);
