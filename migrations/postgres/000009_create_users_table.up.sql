create table if not exists public.users
(
    id       bigserial primary key not null,
    login    varchar(55)        not null,
    email    varchar(255)       not null,
    pwd_hash text               not null,

    constraint users_login_uniq unique (login)
);
