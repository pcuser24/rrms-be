BEGIN;
create table users
(
    id         uuid                        default gen_random_uuid() not null
        constraint idx_user_user_id_primary primary key,
    email      varchar(45)                                           not null,
    group_id   uuid null,
    password   varchar(200) null,
    created_at timestamp(6) with time zone default now()             not null,
    updated_at timestamp(6) with time zone default now()             not null,
    created_by uuid null,
    updated_by uuid null
);

comment on table users is 'Bang user';

-- TODO Chinh them vao cac bang du lieu khoi tao day nhe

END;