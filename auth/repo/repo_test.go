package repo

import (
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pragmataW/pragmaTraffic/auth/models"
	"github.com/stretchr/testify/assert"
)

func TestSelectVerificationCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := Repo{Db: db}
	query := "SELECT verification_code FROM users WHERE email = \\$1"
	rows := sqlmock.NewRows([]string{"verification_code"}).AddRow(123456)

	mock.ExpectPrepare(query).ExpectQuery().WithArgs("example@example.com").WillReturnRows(rows)

	code, err := r.SelectVerificationCode("example@example.com")
	assert.NoError(t, err)
	assert.Equal(t, 123456, code)
}

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := Repo{Db: db}
	query := "INSERT INTO users \\(username, password, email, verification_code, sent_at, is_verified\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("johndoe", "secret", "john@example.com", 123456, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	err = r.InsertUser(models.User{
		ReqUser: models.ReqUser{
			Username: "johndoe",
			Password: "secret",
			Email:    "john@example.com",
		},
		VerifyCode: 123456,
		CreatedAt:  time.Now(),
		IsVerified: false,
	})
	assert.NoError(t, err)
}

func TestSelectUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := Repo{Db: db}

	// Statik bir zaman tanımlıyoruz
	fixedTime, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 12:00:00")

	rows := sqlmock.NewRows([]string{"username", "email", "password", "is_verified", "sent_at"}).
		AddRow("user1", "user1@example.com", "password123", true, fixedTime).
		AddRow("user2", "user2@example.com", "password456", false, fixedTime)

	mock.ExpectPrepare("^SELECT username, email, password, is_verified, sent_at FROM users$").
		ExpectQuery().
		WillReturnRows(rows)

	expectedUsers := []models.User{
		{
			ReqUser: models.ReqUser{
				Username: "user1",
				Email:    "user1@example.com",
				Password: "password123",
			},

			IsVerified: true,
			CreatedAt:  fixedTime,
		},
		{
			ReqUser: models.ReqUser{
				Username: "user2",
				Email:    "user2@example.com",
				Password: "password456",
			},

			IsVerified: false,
			CreatedAt:  fixedTime,
		},
	}

	users, err := repo.SelectUsers()
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChangeVerificationCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := Repo{Db: db}

	mail := "user@example.com"
	newCode := 987654

	query := "UPDATE users SET verification_code=\\$1, sent_at=\\$2, is_verified=\\$3 WHERE email=\\$4"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(newCode, sqlmock.AnyArg(), false, mail).WillReturnResult(sqlmock.NewResult(0, 1))

	err = r.ChangeVerificationCode(mail, newCode)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectIsVerified(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := Repo{Db: db}
	mail := "user@example.com"
	query := "SELECT is_verified FROM users WHERE email = \\$1"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(mail).WillReturnRows(sqlmock.NewRows([]string{"is_verified"}).AddRow(true))

	isVerified, err := r.SelectIsVerified(mail)
	assert.NoError(t, err)
	assert.True(t, isVerified)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChangeIsVerified(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := Repo{Db: db}
	mail := "user@example.com"
	isVerified := true

	query := "UPDATE users SET is_verified=\\$1 WHERE email=\\$2"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(isVerified, mail).WillReturnResult(sqlmock.NewResult(0, 1))

	err = r.ChangeIsVerified(mail, isVerified)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChangePassword(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := Repo{Db: db}
    mail := "user@example.com"
    newPass := "newPassword123"

    mock.ExpectPrepare("UPDATE users SET password=\\$1 WHERE email=\\$2").
        ExpectExec().
        WithArgs(newPass, mail).
        WillReturnResult(sqlmock.NewResult(0, 1))

    err = repo.ChangePassword(mail, newPass)

    assert.NoError(t, err)
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}