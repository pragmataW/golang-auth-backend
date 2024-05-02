package service

import (
	"fmt"

	"github.com/pragmataW/pragmaTraffic/auth/models"
)

func (src AuthService) ChangePasswordService(username string, mail string, newPasswd string, fromMail string, mailServer string, crypt models.Encryptor) error {
	verificationCode, err := GenerateSecureRandomCode()
	if err != nil {
		return err
	}

	mailBody := MailBody{
		FromName:  "Coop. name",
		FromEmail: fromMail,
		ToName:    username,
		ToEmail:   mail,
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
				<h1>Change password code from !</h1>
				<p>Please verify your email address by entering the following verification code:</p>
				<p><strong>Verification Code: %d</strong></p>
				<p>This code is valid for the next 30 minutes. Please do not share this code with anyone.</p>
				<p>If you did not request this code, please ignore this email or contact us for support.</p>
				<p>Thank you for choosing Coop. name!</p>
				<p>Best regards,<br>Coop. name Team</p>
				</body>
				</html>`, verificationCode),
	}

	dbUsers, err := src.Repo.SelectUsers()
	if err != nil {
		return err
	}

	for _, u := range dbUsers{
		decryptedEmail, err := crypt.Decrypt(u.Email)
		fmt.Println(err)
		if err != nil{
			return err
		}

		if u.Username == username && decryptedEmail == mail && u.IsVerified{
			err = src.Repo.ChangeVerificationCode(u.Email, verificationCode)
			if err != nil{
				return err
			}
			encryptedPass, err := crypt.Encrypt(newPasswd)
			if err != nil{
				return err
			}
			err = src.Repo.ChangePassword(u.Email, encryptedPass)
			if err != nil{
				return err
			}
			err = PostRequest(mailServer, mailBody)
			if err != nil{
				return err
			}
			return nil
		}
	}
	return models.UserNotFoundError{}
}
