package controller

import "github.com/pragmataW/pragmaTraffic/auth/models"

type IAuthService interface {
	RegisterUser(user models.ReqUser, mailServer string, fromMail string, crypt models.Encryptor) error
	VerifyCode(mail string, code int, crypt models.Encryptor) error
	LoginUser(mail string, password string, crypt models.Encryptor, key string) (string, error)
	ChangePasswordService(username string, mail string, newPasswd string, fromMail string, mailServer string, crypt models.Encryptor) error
}

type AuthController struct {
	Src        IAuthService
	MailServer string
	FromMail   string
	Crypt      models.Encryptor
	JwtKet     string
}

type VerifyReqBody struct {
	Mail string `json:"mail"`
	Code int    `json:"code"`
}

type LoginReqBody struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type ChangePasswdBody struct {
	Mail      string `json:"mail"`
	Username  string `json:"username"`
	NewPasswd string `json:"newpasswd"`
}
