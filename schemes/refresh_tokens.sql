CREATE TABLE refresh_tokens (
    id serial NOT NULL UNIQUE,
    user_id int NOT NULL REFERENCES users (id),
    token varchar(255) NOT NULL,
    expires_at timestamp NOT NULL default now()
);