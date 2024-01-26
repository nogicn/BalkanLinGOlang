package db

import (
	"database/sql"
	"io"
	"time"

	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	//_ "modernc.org/sqlite"

	"BalkanLinGO/models/activequestiondb"
	"BalkanLinGO/models/dictionarydb"
	"BalkanLinGO/models/dictionaryuserdb"
	"BalkanLinGO/models/languagedb"
	"BalkanLinGO/models/userdb"
	"BalkanLinGO/models/userworddb"
	"BalkanLinGO/models/worddb"
)

var DB *sql.DB

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	os.Remove(dst)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write([]byte{})
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func Init() {
	// create copy of a file
	err := copyFile("./db/testDB.sqlite3", "./db/database.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	//os.Setenv("SQLITE_LOG", "1")
	dsn := "file:./db/database.sqlite3?cache=shared" //&_journal_mode=WAL
	DB, err = sql.Open("sqlite3", dsn)
	/*loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	DB = sqldblogger.OpenDriver(dsn, DB.Driver(), loggerAdapter, sqldblogger.WithSQLQueryAsMessage(true), sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug), // default: LevelInfo
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug), sqldblogger.WithSQLArgsFieldname("sql_args"))
	DB.Ping()*/
	// write ahead logging
	//DB, err = sql.Open("sqlite3", "file::memory:?cache=shared&_journal_mode=WAL")
	//DB, err = sql.Open("mysql", "root:my-secret-pw@/surfpit")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(100)
	DB.SetConnMaxIdleTime(time.Minute * 3)
	// enable logs
	if os.Getenv("MIGRATE") == "true" {
		createAllTables()

		// Insert sample users
		sampleUsers := []string{"admin", "user"}
		sampleSurname := []string{"admin", "user"}
		sampleEmail := []string{"admin@balkanlingo.online", "user@balkanlingo.online"}
		samplePassword := []string{"123", "123"}
		sampleIsAdmin := []int{1, 0}

		for i, u := range sampleUsers {
			var userModel userdb.User
			userModel.Name = u
			userModel.Surname = sampleSurname[i]
			userModel.Email = sampleEmail[i]
			userModel.Password = samplePassword[i]
			userModel.IsAdmin = sampleIsAdmin[i]
			err := userdb.CreateUser(DB, &userModel)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

}

func createAllTables() {
	// Create tables
	if err := userdb.CreateUserTable(DB); err != nil {
		log.Fatal(err, "user.CreateUserTable")
	}
	if err := languagedb.CreateLanguageTable(DB); err != nil {
		log.Fatal(err, "language.CreateLanguageTable")
	}
	if err := dictionarydb.CreateDictionaryTable(DB); err != nil {
		log.Fatal(err, "dictionary.CreateDictionaryTable")
	}
	if err := dictionaryuserdb.CreateDictionaryUserTable(DB); err != nil {
		log.Fatal(err, "dictionaryuser.CreateDictionaryUserTable")
	}
	if err := worddb.CreateWordTable(DB); err != nil {
		log.Fatal(err, "word.CreateWordTable")
	}
	if err := userworddb.CreateUserWordTable(DB); err != nil {
		log.Fatal(err, "userword.CreateUserWordTable")
	}
	if err := activequestiondb.CreateActiveQuestionTable(DB); err != nil {
		log.Fatal(err, "activequestion.CreateActiveQuestionTable")
	}
}
