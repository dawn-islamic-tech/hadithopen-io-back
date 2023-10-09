create table if not exists hadith.stories
(
    id    bigserial primary key,
    title jsonb not null default '{}'
);
