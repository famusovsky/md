package db

import (
	"database/sql"
	"fmt"
	"os"
)

func OpenDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = getDsnFromEnv()
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getDsnFromEnv() string {
	dsn := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=%s host=%s",
		os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_HOST"))

	return dsn
}
