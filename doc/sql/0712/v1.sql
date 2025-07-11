create database basic_system;
use basic_system;

create table api
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    name       varchar(255) not null comment '名称',
    method     varchar(255) not null comment '请求方法',
    path       varchar(255) not null comment '请求路径',
    api_type   int          null comment '1 app  2 menu',
    creator    bigint       null comment '创建者',
    status     int          not null comment '状态（0 可用 1 不可用）',
    `describe` text         null comment '描述',
    created_at datetime(3)  null comment '创建时间',
    updated_at datetime(3)  null comment '更新时间',
    deleted_at datetime(3)  null comment '删除时间'
)
    comment 'API表';

create index idx_deleted_at
    on api (deleted_at);

create table app
(
    id          bigint unsigned auto_increment comment '主键ID'
        primary key,
    name        varchar(255) not null comment '名称',
    app_key     varchar(255) not null comment '密钥',
    entry       varchar(255) null comment '前端入口',
    secret      varchar(255) not null comment '密钥',
    status      int          not null comment '状态',
    sign_method varchar(255) not null comment '加密算法（如md5）',
    proxy       json         null,
    created_at  datetime(3)  null comment '创建时间',
    updated_at  datetime(3)  null comment '更新时间',
    deleted_at  datetime(3)  null comment '删除时间'
)
    comment 'APP表';

create index idx_deleted_at
    on app (deleted_at);

create table auth_relation
(
    role_id           bigint default 0 not null comment '角色ID',
    app_id            bigint default 0 not null comment '应用ID',
    menu_id           bigint default 0 not null comment '菜单ID',
    api_id            bigint default 0 not null comment 'API ID',
    user_id           bigint default 0 not null comment '用户ID',
    button_permission json             null,
    primary key (role_id, app_id, menu_id, api_id, user_id)
)
    comment '权限关联表';

create table department
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    name       varchar(255) not null comment '部门名称',
    code       varchar(255) not null comment '部门编码',
    leader     bigint       not null comment '部门领导（用户ID）',
    parent     bigint       not null comment '父级部门ID',
    created_at datetime(3)  null comment '创建时间',
    updated_at datetime(3)  null comment '更新时间',
    deleted_at datetime(3)  null comment '删除时间'
)
    comment '部门表';

create index idx_deleted_at
    on department (deleted_at);

create index idx_parent
    on department (parent);

create table menu
(
    id          bigint unsigned auto_increment comment '主键ID'
        primary key,
    name        varchar(255) not null comment '菜单名称',
    path        varchar(255) not null comment '菜单路径',
    status      int          not null comment '状态（0 可用 1 不可用）',
    parent      bigint       not null comment '父级菜单ID',
    menu_type   int          null comment '1目录 2菜单 ',
    menu_config json         null,
    sort        float        not null,
    buttons     json         null comment '按钮',
    created_at  datetime(3)  null comment '创建时间',
    updated_at  datetime(3)  null comment '更新时间',
    deleted_at  datetime(3)  null comment '删除时间'
)
    comment '菜单表';

create index idx_deleted_at
    on menu (deleted_at);

create index idx_parent
    on menu (parent);

create table role
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    status     int                                 not null comment '0 可用 1 不可用',
    name       varchar(255)                        not null comment '角色名称',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
)
    comment '角色表';

create index idx_deleted_at
    on role (deleted_at);

create table user
(
    id             bigint unsigned auto_increment comment '主键ID'
        primary key,
    updated_at     datetime(3)                         null comment '更新时间',
    deleted_at     datetime(3)                         null comment '删除时间',
    username       varchar(255)                        not null comment '用户名',
    password       varchar(255)                        not null comment '密码',
    nickname       varchar(255)                        not null comment '昵称',
    supper_manager int                                 not null comment '是否超级管理员（0 不是 1 是）',
    phone          varchar(255)                        null comment '手机号',
    avatar         varchar(255)                        null comment '头像',
    sex            int                                 not null comment '性别',
    email          varchar(255)                        null comment '邮箱',
    dept_id        int                                 null comment '部门',
    post           varchar(255)                        null comment '岗位',
    status         int                                 not null comment '状态（0 可用 1 不可用）',
    created_at     timestamp default CURRENT_TIMESTAMP null comment '创建时间'
)
    comment '用户表';

create index idx_deleted_at
    on user (deleted_at);

