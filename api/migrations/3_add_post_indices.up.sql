ALTER TABLE posts 
    ADD CONSTRAINT slug_unique UNIQUE (slug);

CREATE INDEX created_idx ON posts (created);
