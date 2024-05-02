package repo

import (
	"sync"

	"database/sql"

	_ "github.com/lib/pq"
)

var (
	once sync.Once
	db   *sql.DB
)

const (
	insertUser       = "INSERT INTO users (username, password, email, verification_code, sent_at, is_verified) VALUES ($1, $2, $3, $4, $5, $6)"
	selectUsers      = "SELECT username, email, password, is_verified, sent_at FROM users"
	getCode          = "SELECT verification_code FROM users WHERE email = $1"
	changeCode       = "UPDATE users SET verification_code=$1, sent_at=$2, is_verified=$3 WHERE email=$4"
	getIsVerified    = "SELECT is_verified FROM users WHERE email = $1"
	changeIsVerified = "UPDATE users SET is_verified=$1, sent_at=$2 WHERE email=$3"
	changePasswd     = "UPDATE users SET password=$1 WHERE email=$2"
)

type SqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	Sslmode  string
}

type Repo struct {
	Db *sql.DB
}
