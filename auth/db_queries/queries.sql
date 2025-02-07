/*
*   PostgreSQL DB
*   
*/

/* create db */
CREATE DATABASE messdfs;

/* create tables */
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS friends (
    friendship_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    friend VARCHAR(255) NOT NULL, 
    CONSTRAINT fk_users FOREIGN KEY(username) REFERENCES users(username) ON DELETE SET NULL,
    CONSTRAINT fk_users2 FOREIGN KEY(friend) REFERENCES users(username) ON DELETE SET NULL 
);

CREATE TABLE IF NOT EXISTS directories (
    user_dirs_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    directory VARCHAR(255) NOT NULL UNIQUE,
    CONSTRAINT fk_users FOREIGN KEY(username) REFERENCES users(username) ON DELETE SET NULL
);

/* data insertion */
INSERT INTO users (username, password) VALUES ('paul foo', '1234paulfoo');

INSERT INTO friends (username, friend) VALUES ('paul foo', 'erika bar');

INSERT INTO directories (username, directory) VALUES ('elia renzoni', 'js');