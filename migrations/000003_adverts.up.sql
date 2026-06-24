CREATE TABLE advertapp.advert_categories(
    id        SERIAL                 PRIMARY KEY,
    parent_id INTEGER                REFERENCES advertapp.advert_categories(id) ON DELETE RESTRICT,
    name      VARCHAR(100) NOT NULL,

    CONSTRAINT name_len CHECK(length(trim(name)) BETWEEN 1 AND 100)
);

CREATE        INDEX idx_advert_categories_parent_id ON advertapp.advert_categories(parent_id);
CREATE UNIQUE INDEX unique_parent_name              ON advertapp.advert_categories(parent_id, name) WHERE parent_id IS NOT NULL;
CREATE UNIQUE INDEX unique_root_name                ON advertapp.advert_categories(name) WHERE parent_id IS NULL;

CREATE TABLE advertapp.adverts(
    id          SERIAL                  PRIMARY KEY,
    version     INT           NOT NULL  DEFAULT 1,
    user_id     INTEGER       NOT NULL  REFERENCES advertapp.users(id) ON DELETE CASCADE,
    title       VARCHAR(100)  NOT NULL,
    description VARCHAR(1500) NOT NULL,
    price       BIGINT        NOT NULL,
    category_id INTEGER       NOT NULL  REFERENCES advertapp.advert_categories(id) ON DELETE RESTRICT,
    status      VARCHAR(100)  NOT NULL  DEFAULT 'initial',
    views_count INTEGER       NOT NULL  DEFAULT 0,
    created_at  TIMESTAMPTZ   NOT NULL  DEFAULT now(),
    updated_at  TIMESTAMPTZ   NOT NULL  DEFAULT now(),

    CONSTRAINT version_positive      CHECK(version > 0),
    CONSTRAINT title_len             CHECK(length(trim(title)) BETWEEN 1 AND 100),
    CONSTRAINT description_len       CHECK(length(trim(description)) BETWEEN 1 AND 1500),
    CONSTRAINT price_valid           CHECK(price >= 0),
    CONSTRAINT views_count_positive  CHECK(views_count >= 0),
    CONSTRAINT updated_after_created CHECK(updated_at >= created_at),
    CONSTRAINT status_valid          CHECK(status IN ('initial', 'active', 'rejected', 'blocked', 'archived'))
);

CREATE INDEX idx_adverts_user_id ON advertapp.adverts(user_id);
CREATE INDEX idx_adverts_category_id ON advertapp.adverts(category_id);

CREATE TABLE advertapp.advert_images(
    id         SERIAL                PRIMARY KEY,
    advert_id  INTEGER      NOT NULL REFERENCES advertapp.adverts(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    position   INT          NOT NULL DEFAULT 0,
    path       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT position_positive         CHECK(position >= 0),
    CONSTRAINT unique_advert_id_position UNIQUE(advert_id, position),
    CONSTRAINT name_len                  CHECK(length(trim(name)) BETWEEN 1 AND 255),
    CONSTRAINT path_len                  CHECK(length(trim(path)) BETWEEN 1 AND 255)
);

CREATE INDEX idx_advert_images_advert_id ON advertapp.advert_images(advert_id);