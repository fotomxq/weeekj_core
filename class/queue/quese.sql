create table if not exists test_queue
(
    id        bigserial                                          not null
        constraint test_queue_pk
            primary key,
    create_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    update_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    mod_id    bigint                                             not null,
    status    integer                                            not null,
    params    jsonb                                              not null
);

create unique index if not exists test_queue_id_uindex
    on test_queue (id);

create unique index if not exists test_queue_mod_id_uindex
    on test_queue (mod_id);

