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


    CREATE TABLE IF NOT EXISTS active_question (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER ,
            word_id INTEGER,
            type INTEGER NOT NULL DEFAULT 1,
            FOREIGN KEY (user_id) REFERENCES user(id),
            FOREIGN KEY (word_id) REFERENCES word(id),
            UNIQUE (user_id, word_id)
        );