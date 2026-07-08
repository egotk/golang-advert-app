ALTER TABLE advertapp.adverts
ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('russian', title), 'A') ||
    setweight(to_tsvector('russian', description), 'B')
) STORED;

CREATE INDEX idx_adverts_fts 
ON advertapp.adverts 
USING GIN(search_vector);