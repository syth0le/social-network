CREATE TABLE IF NOT EXISTS user_table
(
    id              TEXT                     NOT NULL,
    username        TEXT UNIQUE              NOT NULL,
    hashed_password TEXT                     NOT NULL,
    first_name      TEXT                     NOT NULL,
    second_name     TEXT                     NOT NULL,
    sex             TEXT                     NOT NULL,
    birthdate       TIMESTAMP WITH TIME ZONE NOT NULL,
    biography       TEXT,
    city            TEXT,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at      TIMESTAMP WITH TIME ZONE,

    CONSTRAINT pk_user_table PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS token_table
(
    id         TEXT                     NOT NULL,
    user_id    TEXT                     NOT NULL,
    token      TEXT                     NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    alived_at  TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT pk_token_table PRIMARY KEY (id),
    CONSTRAINT fk_token_table_user_table FOREIGN KEY (user_id) REFERENCES user_table (id)
);