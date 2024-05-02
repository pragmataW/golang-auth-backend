package repo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pragmataW/pragmaTraffic/auth/models"
	_ "github.com/lib/pq"
)

func NewDb(cnf SqlConfig) *sql.DB {
	once.Do(func() {
		cnnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cnf.Host, cnf.Port, cnf.User, cnf.Password, cnf.DbName, cnf.Sslmode)
		var err error
		db, err = sql.Open("postgres", cnnStr)
		if err != nil {
			log.Fatal(err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}
	})

	return db
}

func (r Repo) InsertUser(usr models.User) error {
	query, err := r.Db.Prepare(insertUser)
	if err != nil {
		return err
	}
	defer query.Close()

	usr.CreatedAt = time.Now()
	usr.IsVerified = false

	_, err = query.Exec(usr.Username, usr.Password, usr.Email, usr.VerifyCode, usr.CreatedAt, usr.IsVerified)
	if err != nil {
		return err
	}

	return nil
}

func (r Repo) SelectUsers() ([]models.User, error) {
	query, err := r.Db.Prepare(selectUsers)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Username, &u.Email, &u.Password, &u.IsVerified, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r Repo) SelectVerificationCode(mail string) (int, error) {
	query, err := r.Db.Prepare(getCode)
	if err != nil {
		return 0, err
	}
	defer query.Close()

	var code int
	err = query.QueryRow(mail).Scan(&code)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func (r Repo) ChangeVerificationCode(mail string, code int) error {
	query, err := r.Db.Prepare(changeCode)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer query.Close()

	_, err = query.Exec(code, time.Now(), false, mail)
	if err != nil {
		return err
	}

	return nil
}

func (r Repo) SelectIsVerified(mail string) (bool, error) {
	query, err := r.Db.Prepare(getIsVerified)
	if err != nil {
		return false, err
	}
	defer query.Close()

	var isVerified bool
	err = query.QueryRow(mail).Scan(&isVerified)
	if err != nil{
		return false, err
	}
	return isVerified, nil
}

func (r Repo) ChangeIsVerified(mail string, isVerified bool) error {
    oldDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

    stmt, err := r.Db.Prepare(changeIsVerified)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(isVerified, oldDate, mail)
    if err != nil {
        return err
    }

    return nil
}

func (r Repo) ChangePassword(mail string, newPass string) error {
    stmt, err := r.Db.Prepare(changePasswd)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(newPass, mail)
    if err != nil {
        return err
    }

    return nil
}