create table if not exists test_comment
(
    id         bigserial                                                            not null
        constraint test_comment_pk
            primary key,
    create_at  timestamp with time zone default CURRENT_TIMESTAMP                   not null,
    delete_at  timestamp with time zone default to_timestamp((0)::double precision) not null,
    comment_id bigint                                                               not null,
    parent_id  bigint                                                               not null,
    org_id     bigint                                                               not null,
    user_id    bigint                                                               not null,
    bind_id    bigint                                                               not null,
    level_type integer                                                              not null,
    level      integer                                                              not null,
    title      varchar(300)                                                         not null,
    des        text                                                                 not null,
    des_files  bigint[]                                                             not null,
    params     jsonb                                                                not null
);

create unique index if not exists test_comment_id_uindex
    on test_comment (id);

