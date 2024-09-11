-- drop database receipt;
-- create database receipt;
-- use receipt;

DROP TABLE IF EXISTS `unique_request`;
create table unique_request
(
    id          bigint unsigned auto_increment
        primary key,
    request_id  char(128)                           not null comment '对成功幂等',
    task_id     char(32)                            not null,
    create_time timestamp default CURRENT_TIMESTAMP not null,
    constraint uq_request_id
        unique (request_id)
)
    comment '防重表，必须，创建更新操作对成功幂等' collate = utf8mb4_general_ci;

DROP TABLE IF EXISTS `task`;
create table task
(
    id          char(32)                               not null
        primary key,
    request_id  char(128)                              not null comment '初始请求ID',
    type        varchar(128)                           not null comment '业务类型',
    state       varchar(128)                           not null comment '任务状态',
    version     int unsigned default 1                 not null,
    create_time timestamp    default CURRENT_TIMESTAMP not null,
    update_time timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint uq_request_id
        unique (request_id)
)
    comment '任务主表，必须，维护状态驱动执行' collate = utf8mb4_general_ci;

create index idx_state
    on task (state);


DROP TABLE IF EXISTS `data`;
create table data
(
    id               bigint unsigned auto_increment
        primary key,
    task_id          char(32)                                  not null,
    token            varchar(64)     default ''                not null,
    valid_from       bigint unsigned default 0                 not null,
    valid_to         bigint unsigned default 0                 not null,
    hash             varchar(128)    default ''                not null,
    from_address     varchar(128)    default ''                not null,
    to_address       varchar(128)    default ''                not null,
    asset            varchar(32)     default ''                not null,
    value            decimal(50, 15) default 0.000000000000000 not null,
    raw_data         JSON,
    transaction_time bigint unsigned default 0                 not null comment '业务时间',
    comment          varchar(128)    default ''                not null comment '备注说明'
)
    comment '业务字段表，必须，根据具体业务设计字段' collate = utf8mb4_general_ci;

create index idx_task_id
    on data (task_id);
