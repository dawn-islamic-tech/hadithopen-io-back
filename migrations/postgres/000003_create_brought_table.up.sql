create table if not exists hadith.brought
(
    id      bigserial primary key,
    brought jsonb default '{}'::jsonb not null
);
