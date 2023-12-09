create table if not exists core_config
(
    mark         varchar(100)                                       not null
        constraint core_config_pk
            primary key,
    create_at    timestamp with time zone default CURRENT_TIMESTAMP not null,
    update_at    timestamp with time zone default CURRENT_TIMESTAMP not null,
    allow_public boolean                                            not null,
    update_hash  varchar(50)                                        not null,
    name         varchar(300)                                       not null,
    group_mark   varchar(100)                                       not null,
    des          varchar(600)                                       not null,
    value_type   varchar(10)                                        not null,
    value        text                                               not null
);

create unique index if not exists core_config_mark_uindex
    on core_config (mark);

