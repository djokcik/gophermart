DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        CREATE TYPE "order_status" AS ENUM (
            'NEW',
            'PROCESSING',
            'PROCESSED',
            'INVALID'
        );
    END IF;
END$$;

create table users
(
    id serial not null,
    username varchar(50) not null,
    password varchar(255) not null,
    created_at timestamp default current_timestamp,
    balance int default 0 not null
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

create table withdraw_log
(
    id serial not null,
    user_id int not null,
    sum int default 0,
    processed_at timestamp default current_timestamp,
    order_id text
);

create unique index withdraw_log_id_uindex
    on withdraw_log (id);

alter table withdraw_log
    add constraint withdraw_log_pk
        primary key (id);

