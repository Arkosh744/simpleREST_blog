CREATE TABLE posts
(
    id        serial       not null unique,
    title     varchar(255) not null,
    body      varchar(255) not null,
    createdAt timestamp    not null default now(),
    updatedAt timestamp    not null default now()
);