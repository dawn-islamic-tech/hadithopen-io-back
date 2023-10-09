create table if not exists hadith.translates
(
    id        bigserial primary key not null,
    lang      varchar(3)            not null,
    translate text                  not null,

    constraint hadith_translates_lang_fkey foreign key (lang) references hadith.lang (lang)
);
