DROP TABLE IF EXISTS results;

CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    image VARCHAR(255),
    chemistry INT,
    biology INT,
    maths INT,
    civic INT,
    english INT,
    aptitude INT,
    physics INT
);

CREATE INDEX IF NOT EXISTS idx_results_id ON results (id);
