DROP TABLE IF EXISTS posts;

CREATE TABLE IF NOT EXISTS posts(
                                    id SERIAL PRIMARY KEY,
                                    title TEXT UNIQUE NOT NULL,
                                    content TEXT NOT NULL,
                                    pubdate TEXT NOT NULL,
                                    pubtime BIGINT NOT NULL,
                                    link TEXT NOT NULL
);
