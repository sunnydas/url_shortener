package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
)

// NewDatabase will open a database connection and then try to migrate based on ENV
func NewDatabase() (*sql.DB, error) {

	db, err := CreateDatabase()
	if err != nil {
		return db, err
	}

	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "test" {
		err := migrateDatabase(db, false)
		if err != nil {
			return db, err
		}
	}

	return db, nil
}

func migrateDatabase(db *sql.DB, silence bool) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	absPath, _ := filepath.Abs("./")
	rootPath := strings.Split(absPath, "/url_shortener")[0]
	if rootPath == "/" {
		rootPath = ""
	}
	fullPath := fmt.Sprintf("%s/url_shortener/internal/db/migrations", rootPath)
	fmt.Printf("Searching for database migrations in %s \n", fullPath)
	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", fullPath),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if !silence {
		migration.Log = &MigrationLogger{}
		migration.Log.Printf("Applying database migrations")
	}
	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	version, _, err := migration.Version()
	if err != nil {
		return err
	}

	if !silence {
		migration.Log.Printf("Active database version: %d", version)
	}

	return nil
}

// CreateDatabase will open a connection to a database
func CreateDatabase() (*sql.DB, error) {
	// Load env variables
	serverName := os.Getenv("MYSQL_SERVER")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	db, err := sql.Open("mysql", user+":"+password+"@tcp("+serverName+")/"+dbName+"?charset=utf8mb4,utf8&parseTime=true&multiStatements=true")
	if err != nil {
		log.Fatal("Cannot open database connection. ", err)
	}

	// 10min lifetime
	db.SetConnMaxLifetime(time.Second * 600)

	return db, nil
}
