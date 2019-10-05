CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100)
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
    modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
    slug VARCHAR(50) NOT NULL,
    title VARCHAR(100) NOT NULL,
    body TEXT NOT NULL,
    picture TEXT NOT NULL,
    description TEXT NOT NULL,
    published BOOLEAN NOT NULL,
    author_id integer references users (id)
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE posttags (
    id SERIAL PRIMARY KEY,
    tag_id integer references tags (id),
    post_id integer references posts (id)
);