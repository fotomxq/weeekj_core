create table if not exists test_tag
(
    id        bigserial                                          not null
        constraint test_tag_pk
            primary key,
    create_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    update_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    bind_id   bigint                                             not null,
    name      varchar(300)                                       not null
);

create unique index if not exists test_tag_id_uindex
    on test_tag (id);

