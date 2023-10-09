create table if not exists hadith.comments
(
    id         bigserial primary key not null,
    story_id   bigint                not null,
    brought_id bigint                not null,
    comment    jsonb                 not null default '{}',

    constraint hadith_comments_story_id_fkey foreign key (story_id) references hadith.stories (id),
    constraint hadith_comments_brought_id_fkey foreign key (brought_id) references hadith.brought (id)
);