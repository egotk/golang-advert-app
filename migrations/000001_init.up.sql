CREATE SCHEMA advertapp;

CREATE TABLE advertapp.users(
    id                 SERIAL                PRIMARY KEY,
    version            INT          NOT NULL DEFAULT 1,
    email              VARCHAR(255) NOT NULL UNIQUE,
    full_name          VARCHAR(100) NOT NULL,
    phone_number       VARCHAR(20)  NOT NULL UNIQUE,
    password_hash      VARCHAR(100) NOT NULL,
    role               VARCHAR(100) NOT NULL DEFAULT 'user',
    failed_login_count INT          NOT NULL DEFAULT 0,
    locked_until       TIMESTAMP,
    created_at         TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at         TIMESTAMP    NOT NULL DEFAULT now(),
    image_path         VARCHAR(255),

    CONSTRAINT version_positive         CHECK (version > 0),
    CONSTRAINT email_regex_valid        CHECK (email ~ '^[^\s@]+@[^\s@]+\.[^\s@]+$'),
    CONSTRAINT email_len_valid          CHECK (char_length(email) BETWEEN 3 and 255),
    CONSTRAINT email_lowercase          CHECK (email = lower(email)),
    CONSTRAINT full_name_len_valid      CHECK (char_length(full_name) BETWEEN 3 AND 100),
    CONSTRAINT phone_number             CHECK (phone_number ~ '^\+[1-9]\d{1,14}$'),
    CONSTRAINT phone_number_len         CHECK (char_length(phone_number) BETWEEN 4 AND 20),
    CONSTRAINT password_hash_len_valid  CHECK (char_length(password_hash) BETWEEN 1 AND 100),
    CONSTRAINT user_role_len_valid      CHECK (char_length(role) BETWEEN 1 AND 100),
    CONSTRAINT user_role_valid          CHECK (role IN ('user', 'admin')),
    CONSTRAINT failed_login_count_valid CHECK (failed_login_count BETWEEN 0 AND 5),
    CONSTRAINT locked_after_created     CHECK (locked_until > created_at),
    CONSTRAINT updated_after_created    CHECK (updated_at >= created_at)
);