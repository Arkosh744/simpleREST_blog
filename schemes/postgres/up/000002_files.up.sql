CREATE TABLE files
(
    id        serial       not null primary key,
    name      varchar(255) not null,
    author_id integer      not null references users (id) default 1,
    comments  varchar(255)                                default 'No comments',
    createdAt timestamp    not null                       default now()
);