package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

// Initiation connect to database
func Initiation() {

	var (
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASS")
		dbHost = os.Getenv("DB_HOST")
		dbPort = os.Getenv("DB_PORT")
		dbName = os.Getenv("DB_NAME")
	)

	tl, err := time.LoadLocation("Asia/Jakarta") // Load time location
	if err != nil {
		log.Fatalf("database: could not load time location: %s", err.Error())
	}
	mc := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Addr:                 fmt.Sprintf("%s:%s", dbHost, dbPort),
		Net:                  "tcp",
		DBName:               dbName,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  tl,
		ParseTime:            true,
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	log.Println("database: trying connect to database")

	createConnectionDB(mc.FormatDSN())
	log.Println("database: successfully connected to database")
}

// create database connection
func createConnectionDB(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("database: could not open connection to database: %s", err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("database: could not ping to database: %s", err.Error())
	}
	// Configure setting
	optimizeDBPerformance(db)
	// Check default table on database
	checkDefaultTable(db)
	// assign db pointer to private dbConn variable
	dbConn = db
}

// Optimize DB setting for better performance
// source: https://www.alexedwards.net/blog/configuring-sqldb
func optimizeDBPerformance(db *sql.DB) {
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Duration(300 * time.Second))
}

// check default table if not exists
// create table fruits
func checkDefaultTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS fruits (
			fruit_id			INT(6)			UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			fruit_name			VARCHAR(155)	NOT NULL,
			fruit_color			VARCHAR(55)		NOT NULL,
			fruit_image			VARCHAR(255)	DEFAULT '',
			fruit_created_at	TIMESTAMP		DEFAULT CURRENT_TIMESTAMP,
			fruit_updated_at	TIMESTAMP		DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) engine=innodb;
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("database: could not prepare query: %s", err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("database: could not exec query: %s", err.Error())
	}
}

// GetMariaDBconn get pointer database connection
func GetMariaDBconn() *sql.DB {
	return dbConn
}

// ReplaceDBConn this function for testing purpose only
// it will replace dbConn pointer to DBMock pointer from go-sqlmock
func ReplaceDBConn(DBMockConn *sql.DB) {
	dbConn = DBMockConn
}
