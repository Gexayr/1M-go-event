CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    generated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    content TEXT NOT NULL
);
