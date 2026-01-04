package db

import (
	"context"
	"database/sql"
	"embed"
	"io"
	"path/filepath"

	"log"
	"os"

	//_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"

	_ "modernc.org/sqlite"

	repository "BalkanLinGO/internal/db/repository"

	"github.com/pressly/goose/v3"
)

//var DB *sql.DB
//var Queries *sqlcdb.Queries

//go:embed migrations/*.sql
var embedMigrations embed.FS

var (
	dburl      string
	dbInstance *service
)

type Service interface {
	Close() error
	GetRepositoryRW() *repository.Queries
	GetRepositoryRO() *repository.Queries
	Init()
}

type service struct {
	DBrw  *sql.DB
	Reprw *repository.Queries

	DBro  *sql.DB
	Repro *repository.Queries
}

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

func New(dburlOverride string) Service {
	dburl = dburlOverride
	doesExist := false

	if dburl == "" {
		log.Fatal("BLUEPRINT_DB_URL is not set, check your .env file")
	}

	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// check if file exists
	if _, err := os.Stat(dburl); err == nil {
		doesExist = true
	}

	// Ensure the directory exists
	dir := filepath.Dir(dburl)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}

	// Use mode=rwc so SQLite creates the DB file if it doesn't exist
	dbrw, err := sql.Open("sqlite", "file:"+dburl+"?mode=rwc&_journal_mode=WAL&_txlock=immediate")
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	// Verify we can connect / open the DB now
	if err := dbrw.Ping(); err != nil {
		log.Fatalf("failed to connect to sqlite db: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatalf("goose set dialect failed: %v", err)
	}

	dbrw.SetMaxOpenConns(1)
	dbrw.SetMaxIdleConns(1)

	dbro, err := sql.Open("sqlite", "file:"+dburl+"?mode=ro&_journal_mode=WAL")
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	// Verify we can connect / open the DB now
	if err := dbro.Ping(); err != nil {
		log.Fatalf("failed to connect to sqlite db: %v", err)
	}

	dbro.SetMaxIdleConns(1)

	dbInstance = &service{
		DBrw:  dbrw,
		Reprw: repository.New(dbrw),
		DBro:  dbro,
		Repro: repository.New(dbro),
	}

	if !doesExist {
		// Run migrations

		if err := goose.Up(dbrw, "migrations"); err != nil {
			log.Fatalf("goose up failed: %v", err)
		}
		dbInstance.Init()
	}

	return dbInstance
}

// Insert sample data
func (s *service) Init() {
	// users
	ctx := context.Background()
	sampleUsers := []string{"admin", "user"}
	sampleSurname := []string{"admin", "user"}
	sampleEmail := []string{"admin@balkanlingo.online", "user@balkanlingo.online"}
	samplePassword := []string{"123", "123"}
	sampleIsAdmin := []bool{true, false}

	for i, u := range sampleUsers {
		if sampleIsAdmin[i] {
			err := dbInstance.Reprw.CreateAdmin(ctx, repository.CreateAdminParams{
				Name:     u,
				Surname:  sampleSurname[i],
				Email:    sampleEmail[i],
				Password: samplePassword[i],
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := dbInstance.Reprw.CreateUser(ctx, repository.CreateUserParams{
				Name:     u,
				Surname:  sampleSurname[i],
				Email:    sampleEmail[i],
				Password: samplePassword[i],
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// dictionaries

}

func (s *service) Close() error {
	return s.DBrw.Close()
}

func (s *service) GetRepositoryRW() *repository.Queries {
	return s.Reprw
}
func (s *service) GetRepositoryRO() *repository.Queries {
	return s.Repro
}
