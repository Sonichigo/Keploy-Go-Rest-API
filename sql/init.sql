CREATE TABLE IF NOT EXISTS posts (
    id char(8) NOT NULL,
    title char(150) NOT NULL,
    content char(2000) NOT NULL,
    PRIMARY KEY(id)
);