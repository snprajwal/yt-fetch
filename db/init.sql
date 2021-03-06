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

CREATE INDEX video_id_published_at_idx ON video (id DESC, published_at DESC);
CREATE INDEX video_title_description_idx ON video (title, description);
