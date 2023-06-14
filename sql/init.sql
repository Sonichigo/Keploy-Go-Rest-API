CREATE TABLE IF NOT EXISTS posts (
    id SERIAL ,
    title VarChar(150),
    content VarChar(2000),
    PRIMARY KEY(id)
);