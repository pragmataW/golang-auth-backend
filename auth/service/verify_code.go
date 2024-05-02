package service

import (
	"fmt"
	"time"

	"github.com/pragmataW/pragmaTraffic/auth/models"
)

func (s AuthService) VerifyCode(mail string, code int, crypt models.Encryptor) error {
	users, err := s.Repo.SelectUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		decryptedEmail, err := crypt.Decrypt(u.Email)
		if err != nil {
			return err
		}

		if decryptedEmail == mail {
			dbCode, err := s.Repo.SelectVerificationCode(u.Email)
			if err != nil {
				fmt.Println(err)
				return err
			}

			if dbCode == code && time.Since(u.CreatedAt).Seconds() <= 30 * 60  {
				err = s.Repo.ChangeIsVerified(u.Email, true)
				if err != nil {
					return err
				}

				return nil
			}

			if dbCode != code{
				return models.CodeNotFoundError{}
			}
			if time.Since(u.CreatedAt).Seconds() > 30 * 60 {
				return models.CodeExpiredError{}
			}
		}
	}
	return models.UserNotFoundError{}
}
