package service

import (
	"github.com/pragmataW/pragmaTraffic/auth/models"
)

type IAuthRepo interface {
	SelectVerificationCode(mail string) (int, error)
	SelectUsers() ([]models.User, error)
	InsertUser(usr models.User) error
	ChangeVerificationCode(mail string, code int) error
	ChangeIsVerified(mail string, isVerified bool) error
	SelectIsVerified(mail string) (bool, error)
	ChangePassword(mail string, newPass string) error
}

type AuthService struct {
	Repo       IAuthRepo
}

type MailBody struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ToName    string `json:"to_name"`
	ToEmail   string `json:"to_email"`
	Subject   string `json:"subject"`
	Html      string `json:"html"`
}
