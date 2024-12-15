CREATE EXTENSION citext;
CREATE DOMAIN d_email AS citext
    CHECK ( value ~
            '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );

CREATE TABLE client
(
    id              SERIAL PRIMARY KEY,
    first_name      VARCHAR(50) NOT NULL,
    last_name       VARCHAR(50) NOT NULL,
    username        VARCHAR(50) NOT NULL UNIQUE,
    hashed_password VARCHAR     NOT NULL,
    email           D_EMAIL     NOT NULL UNIQUE,
    is_verified     BOOLEAN     NOT NULL DEFAULT FALSE,
    code            VARCHAR(6),
    exp             TIMESTAMP            DEFAULT (now() AT TIME ZONE 'utc' + '1 minute'::interval)
);

CREATE INDEX idx_full_name ON client (first_name, last_name);
CREATE INDEX idx_username ON client (username);
CREATE INDEX idx_email ON client (email);