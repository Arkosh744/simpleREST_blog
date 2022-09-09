CREATE TABLE posts
(
    id        serial       not null unique,
    title     varchar(255) not null,
    body      varchar(255) not null,
    author_id integer      not null references users (id),
    createdAt timestamp    not null default now(),
    updatedAt timestamp    not null default now()
);

ALTER TABLE posts
    ADD author_id integer not null references users (id) default 1;