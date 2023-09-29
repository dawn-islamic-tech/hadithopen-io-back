create table if not exists hadith.lang
(
    lang varchar(3)  not null unique,
    name varchar(55) not null
);
