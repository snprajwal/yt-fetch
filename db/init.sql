DROP DATABASE IF EXISTS yt_fetch;
CREATE DATABASE yt_fetch;
\c yt_fetch;

CREATE TABLE video (
	id SERIAL PRIMARY KEY,
	slug VARCHAR,
	title VARCHAR,
	channel VARCHAR,
	description VARCHAR,
	thumbnail VARCHAR,
	published_at TIMESTAMP
);

CREATE INDEX published_at_idx ON video (published_at DESC);
