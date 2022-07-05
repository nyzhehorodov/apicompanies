-- Create initial database structure.

CREATE TABLE companies (
        id SERIAL PRIMARY KEY,
        name VARCHAR (64) NOT NULL,
        country VARCHAR (16) NOT NULL,
        website VARCHAR (256) NOT NULL,
        phone TIMESTAMP
);

---- create above / drop below ----

DROP TABLE companies CASCADE;
