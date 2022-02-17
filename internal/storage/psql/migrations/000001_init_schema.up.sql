create table users
(
    id serial not null,
    username varchar(50) not null,
    password varchar(255) not null,
    created_at timestamp default current_timestamp
);

create unique index users_id_uindex
    on users (username);

alter table users
    add constraint users_pk
        primary key (id);

