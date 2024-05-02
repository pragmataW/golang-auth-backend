package controller

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pragmataW/pragmaTraffic/auth/models"
)

func NewController(src IAuthService, mailServer string, fromMail string, crypt models.Encryptor, jwtKey string) AuthController {
	return AuthController{
		Src:        src,
		MailServer: mailServer,
		FromMail:   fromMail,
		Crypt:      crypt,
		JwtKet: jwtKey,
	}
}

func (ctrl AuthController) RegisterController(c *fiber.Ctx) error {
	var req models.ReqUser

	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "wrong json format",
		})
	}

	err = ctrl.Src.RegisterUser(req, ctrl.MailServer, ctrl.FromMail, ctrl.Crypt)
	if err != nil {
		if _, ok := err.(models.UserAlreadyExistsError); ok {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"Message": "user already exists",
			})
		} else if _, ok := err.(models.MailCouldNotSent); ok {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"Message": "mail could not sent",
			})
		} else if _, ok := err.(models.CodeNotExpiredError); ok {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"Message": "code is already sent",
			})
		} else if _, ok := err.(models.UsernameLengthError); ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"Message": "username length must be greater than 3",
			})
		} else if _, ok := err.(models.PasswordLengthError); ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"Message": "password length must be greater than 8",
			})
		} else if _, ok := err.(models.MailEmptyError); ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"Message": "mail cannot be empty",
			})
		} else if _, ok := err.(models.UsernameAlreadyExistsError); ok {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"Message": "username already exists",
			})
		} else if _, ok := err.(models.MailAlreadyExistsError); ok {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"Message": "email already exists",
			})
		} else if _, ok := err.(models.CodeResent); ok {
			return c.Status(http.StatusOK).JSON(fiber.Map{
				"Message": "code re-sent",
			})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"Message": "internal server error",
			})
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"Message": "user created",
	})
}

func (ctrl AuthController) VerifyCodeController(c *fiber.Ctx) error {
	var req VerifyReqBody
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "wrong json format",
		})
	}

	err = ctrl.Src.VerifyCode(req.Mail, req.Code, ctrl.Crypt)
	if err != nil {
		if _, ok := err.(models.CodeExpiredError); ok {
			return c.Status(http.StatusGone).JSON(fiber.Map{
				"message": "code is expired",
			})
		} else if _, ok := err.(models.CodeNotFoundError); ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "wrong code",
			})
		} else if _, ok := err.(models.UserNotFoundError); ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "user not found",
			})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verified",
	})
}

func (ctrl AuthController) LoginController(c *fiber.Ctx) error {
	var req LoginReqBody
	err := c.BodyParser(&req)
	if err != nil{
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "wrong json format",
		})
	}

	jwt, err := ctrl.Src.LoginUser(req.Mail, req.Password, ctrl.Crypt, ctrl.JwtKet)
	if err != nil {
		if _, ok := err.(models.UserNotFoundError); ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "user not found",
			})
		}else if _, ok := err.(models.NotVerifiedError); ok {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": "user not verified",
			})
		}else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}

	c.Cookie(&fiber.Cookie{
		Name: "Authentication",
		Value: jwt,
		HTTPOnly: true,
		Expires: time.Now().Add(72 * time.Hour),
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "status ok",
	})
}

func (ctrl AuthController) LogoutController(c *fiber.Ctx) error{
	c.ClearCookie("Authentication")
    
    return c.Status(http.StatusOK).JSON(fiber.Map{
        "message": "You have been successfully logged out.",
    })
}

func (ctrl AuthController) ChangePasswdController(c *fiber.Ctx) error {
	var req ChangePasswdBody
	
	err := c.BodyParser(&req)
	if err != nil{
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "wrong json format",
		})
	}

	err = ctrl.Src.ChangePasswordService(req.Username, req.Mail, req.NewPasswd, ctrl.FromMail, ctrl.MailServer, ctrl.Crypt)
	if err != nil{
		if _, ok := err.(models.UserNotFoundError); ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "user not found",
			})
		}else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error " + err.Error(),
			})
		}
	}
	c.ClearCookie("Authentication")
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "code sent, verify your account",
	})
}