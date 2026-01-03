-- +goose Up
-- +goose StatementBegin
-- User table
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT 0,
    token TEXT DEFAULT NULL,
    UNIQUE (email)
);

-- Language table
CREATE TABLE IF NOT EXISTS language (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    shorthand TEXT NOT NULL UNIQUE,
    flag_icon TEXT NOT NULL
);

-- Dictionary table
CREATE TABLE IF NOT EXISTS dictionary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    language_id INTEGER NOT NULL,
    image_link TEXT NOT NULL,
    FOREIGN KEY (language_id) REFERENCES language(id)
);

-- Dictionary User table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS dictionary_user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    dictionary_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (dictionary_id) REFERENCES dictionary(id)
);

-- Word table
CREATE TABLE IF NOT EXISTS word (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    foreignWord TEXT NOT NULL,
    foreignDescription TEXT NOT NULL,
    nativeWord TEXT NOT NULL,
    nativeDescription TEXT NOT NULL,
    pronunciation TEXT NOT NULL,
    dictionary_id INTEGER NOT NULL,
    FOREIGN KEY (dictionary_id) REFERENCES dictionary(id),
    UNIQUE (foreignWord, foreignDescription, nativeWord, nativeDescription, dictionary_id)
);

-- User Word table (many-to-many relationship with additional data)
CREATE TABLE IF NOT EXISTS user_word (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    last_answered TEXT,
    delay INTEGER,
    active INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (word_id) REFERENCES word(id),
    FOREIGN KEY (user_id) REFERENCES user(id)
);

-- Index for user_word
CREATE INDEX IF NOT EXISTS user_word_id_index ON user_word (user_id);

-- Active Question table
CREATE TABLE IF NOT EXISTS active_question (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    word_id INTEGER,
    type INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (word_id) REFERENCES word(id),
    UNIQUE (user_id, word_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- User table
DROP TABLE IF EXISTS active_question;
DROP TABLE IF EXISTS user_word;
DROP TABLE IF EXISTS word;
DROP TABLE IF EXISTS dictionary_user;
DROP TABLE IF EXISTS dictionary;
DROP TABLE IF EXISTS language;
DROP TABLE IF EXISTS user;
-- +goose StatementEnd
