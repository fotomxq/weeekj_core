create table if not exists test_config_default
(
    id              bigserial                                          not null
        constraint test_config_default_pk
            primary key,
    mark            varchar(100)                                       not null,
    name            varchar(300)                                       not null,
    create_at       timestamp with time zone default CURRENT_TIMESTAMP not null,
    update_at       timestamp with time zone default CURRENT_TIMESTAMP not null,
    allow_public    boolean                                            not null,
    allow_self_view boolean                                            not null,
    allow_self_set  boolean                                            not null,
    value_type      integer                                            not null,
    value_check     text                                               not null,
    value_default   text                                               not null
);

create unique index if not exists test_config_default_id_uindex
    on test_config_default (id);

create unique index if not exists test_config_default_mark_uindex
    on test_config_default (mark);

create table if not exists test_config
(
    id        bigserial                                          not null
        constraint test_config_pk
            primary key,
    create_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    update_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    bind_id    bigint                                             not null,
    mark      varchar(100)                                       not null,
    val       text                                               not null
);

create unique index if not exists test_config_uindex
    on test_config (id);

