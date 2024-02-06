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