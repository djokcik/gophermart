CREATE TYPE "order_status" AS ENUM (
    'NEW',
    'REGISTERED',
    'PROCESSING',
    'PROCESSED',
    'INVALID'
);

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

create table orders
(
    id text not null,
    user_id int not null
        constraint orders_users_id_fk
            references users
            on update cascade on delete cascade,
    status order_status not null,
    uploaded_at timestamp default current_timestamp,
    accrual int default 0 not null
);

create unique index orders_id_uindex
    on orders (id);

alter table orders
    add constraint orders_pk
        primary key (id);