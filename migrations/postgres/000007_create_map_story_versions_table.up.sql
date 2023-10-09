create table if not exists hadith.map_story_versions
(
    story_id   bigint not null,
    version_id bigint not null,

    constraint hadith_map_story_versions_version_id_uniq unique (version_id),
    constraint hadith_map_story_versions_story_id_fkey foreign key (story_id) references hadith.stories (id),
    constraint hadith_map_story_versions_version_id_fkey foreign key (version_id) references hadith.versions (id)
);
