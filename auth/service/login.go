package service

import (
	"github.com/pragmataW/pragmaTraffic/auth/models"
)

func (s AuthService) LoginUser(mail string, password string, crypt models.Encryptor, key string) (string, error) {
	users, err := s.Repo.SelectUsers()
	if err != nil {
		return "", err
	}

	for _, u := range users {
		decryptedMail, err := crypt.Decrypt(u.Email)
		if err != nil {
			return "", err
		}
		decryptedPassword, err := crypt.Decrypt(u.Password)
		if err != nil {
			return "", err
		}

		if decryptedMail == mail && decryptedPassword == password {
			if u.IsVerified {
				jwt, err := GenerateJWT(u.Username, "auth-service", key)
				if err != nil {
					return "", err
				}
				return jwt, err
			} else {
				return "", models.NotVerifiedError{}
			}

		}
	}
	return "", models.UserNotFoundError{}
}
