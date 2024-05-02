package service

import (
	"fmt"
	"time"

	"github.com/pragmataW/pragmaTraffic/auth/models"
)

func (s AuthService) RegisterUser(user models.ReqUser, mailServer string, fromMail string, crypt models.Encryptor) error {
	if len(user.Username) < 3 {
		return models.UsernameLengthError{}
	}

	if len(user.Password) < 8 {
		return models.PasswordLengthError{}
	}

	if len(user.Email) == 0 {
		return models.MailEmptyError{}
	}

	dbUsers, err := s.Repo.SelectUsers()
	if err != nil {
		return err
	}

	verificationCode, err := GenerateSecureRandomCode()
	if err != nil {
		return err
	}

	mailBody := MailBody{
		FromName:  "Coop. name",
		FromEmail: fromMail,
		ToName:    user.Username,
		ToEmail:   user.Email,
		Subject:   "Verification Code",
		Html: fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Document</title>
				</head>
				<body>
				<h1>Code from Coop. name!</h1>
				<p>Please verify your email address by entering the following verification code:</p>
				<p><strong>Verification Code: %d</strong></p>
				<p>This code is valid for the next 30 minutes. Please do not share this code with anyone.</p>
				<p>If you did not request this code, please ignore this email or contact us for support.</p>
				<p>Thank you for choosing Coop. name!</p>
				<p>Best regards,<br>Coop. name Team</p>
				</body>
				</html>`, verificationCode),
	}

	for _, u := range dbUsers {
		decryptedEmail, err := crypt.Decrypt(u.Email)
		fmt.Println(err)
		if err != nil{
			return err
		}

		if u.Username == user.Username && decryptedEmail == user.Email { 
			if !u.IsVerified {
				if time.Since(u.CreatedAt).Seconds() <= 30 * 60 {
					return models.CodeNotExpiredError{}
				} else {
					err = PostRequest(mailServer, mailBody)
					if err != nil {
						return models.MailCouldNotSent{}
					}
					s.Repo.ChangeVerificationCode(u.Email, verificationCode)
					return models.CodeResent{}
				}
			} else {
				return models.UserAlreadyExistsError{}
			}
		}

		if u.Username == user.Username {
			return models.UsernameAlreadyExistsError{}
		}
		if decryptedEmail == user.Email {
			return models.MailAlreadyExistsError{}
		}
	}

	err = PostRequest(mailServer, mailBody)
	if err != nil {
		return models.MailCouldNotSent{}
	}

	cryptedEmail, err := crypt.Encrypt(user.Email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	cryptedPassword, err := crypt.Encrypt(user.Password)
	if err != nil {
		fmt.Println(err)
		return err
	}

	usr := models.User{
		ReqUser: models.ReqUser{
			Username: user.Username,
			Email:    cryptedEmail,
			Password: cryptedPassword,
		},
		VerifyCode: verificationCode,
		CreatedAt:  time.Now(),
		IsVerified: false,
	}

	err = s.Repo.InsertUser(usr)
	if err != nil {
		return err
	}
	return nil
}