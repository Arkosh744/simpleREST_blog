CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    email         varchar(255) not null unique,
    password      varchar(255) not null,
    registered_at timestamp    not null default now()
);

CREATE TABLE refresh_tokens (
                                id serial NOT NULL UNIQUE,
                                user_id int NOT NULL REFERENCES users (id),
                                token varchar(255) NOT NULL,
                                expires_at timestamp NOT NULL default now()
);

CREATE TABLE posts
(
    id        serial       not null unique,
    title     varchar(255) not null,
    body      varchar(255) not null,
    author_id integer      not null references users (id) default 1,
    createdAt timestamp    not null default now(),
    updatedAt timestamp    not null default now()
);

INSERT INTO public.users (id, name, email, password) VALUES (1, 'admin', 'go@golang.com', '53616c74792053616c74f7a9e24777ec23212c54d7a350bc5bea5477fdbb');
-- pass: aaaaaa