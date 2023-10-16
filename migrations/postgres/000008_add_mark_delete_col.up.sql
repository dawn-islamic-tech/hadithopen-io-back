alter table if exists hadith.stories
    add column if not exists mark_delete boolean not null default false
;
