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