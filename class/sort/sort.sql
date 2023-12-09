create table if not exists test_sort
(
    id        bigserial                                          not null
        constraint test_sort_pk
            primary key,
    create_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    update_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    bind_id   bigint                                             not null,
    mark      varchar(100)                                       not null,
    parent_id bigint                                             not null,
    sort      integer                                            not null,
    cover_file_id bigint                                         not null,
    des_files bigint[] not null,
    name      varchar(300)                                       not null,
    des       text                                               not null,
    params    jsonb                                              not null
);

create unique index if not exists test_sort_id_uindex
    on test_sort (id);

