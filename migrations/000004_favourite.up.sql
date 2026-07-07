CREATE TABLE advertapp.favourites(
    advert_id  INTEGER     NOT NULL       REFERENCES advertapp.adverts(id) ON DELETE CASCADE,
    user_id    INTEGER     NOT NULL       REFERENCES advertapp.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL       DEFAULT now(),

    PRIMARY KEY(advert_id, user_id)
);

ALTER TABLE advertapp.adverts ADD COLUMN fav_count INT NOT NULL DEFAULT 0;