create table if not exists hadith.versions
(
    id         bigserial                not null primary key,
    brought_id bigint                   not null,
    is_default boolean default false    not null,
    original   text    default ''::text not null,
    version    jsonb                    null
);
