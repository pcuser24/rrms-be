BEGIN;
create table IF NOT EXISTS users
(
    id         uuid                        default gen_random_uuid() not null
        constraint idx_user_user_id_primary primary key,
    email      varchar(45)                                           not null,
    group_id   uuid null,
    password   varchar(200) null,
    created_at timestamp(6) with time zone default now()             not null,
    updated_at timestamp(6) with time zone default now()             not null,
    created_by uuid null,
    updated_by uuid null,
    deleted_f   bool default false not null
);
COMMENT on table users is 'Bang user';
COMMENT  on COLUMN users.deleted_f is '1: deleted, 0: not deleted';

INSERT INTO users(email, password) values ('chinhpl+admin@gmail.com','123456');
INSERT INTO users(email, password) values ('chinhpl+admin1@gmail.com','123456');

-- TODO Chinh them vao cac bang du lieu khoi tao day nhe

END;