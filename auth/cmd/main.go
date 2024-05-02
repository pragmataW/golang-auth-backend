package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pragmataW/pragmaTraffic/auth/controller"
	"github.com/pragmataW/pragmaTraffic/auth/models"
	"github.com/pragmataW/pragmaTraffic/auth/repo"
	"github.com/pragmataW/pragmaTraffic/auth/service"
)

var (
	host       string
	port       int
	user       string
	password   string
	dbName     string
	sslMode    string
	mailServer string
	fromMail   string
	key        string
)

func main() {
	conf := repo.SqlConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbName,
		Sslmode:  sslMode,
	}

	r := repo.Repo{
		Db: repo.NewDb(conf),
	}

	src := service.AuthService{
		Repo: r,
	}

	crypt := models.NewEncryptor(key)

	ctrl := controller.NewController(src, mailServer, fromMail, *crypt, key)

	app := fiber.New()
	app.Post("/register", ctrl.RegisterController)
	app.Post("/verifyCode", ctrl.VerifyCodeController)
	app.Post("/login", ctrl.LoginController)
	app.Post("/logout", ctrl.LogoutController)
	app.Post("/changePassword", ctrl.ChangePasswdController)
	log.Fatal(app.Listen(":2000"))
}

func init() {
	host = os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, _ = strconv.Atoi(portStr)
	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbName = os.Getenv("POSTGRES_DB")
	sslMode = os.Getenv("POSTGRES_SSL")
	mailServer = os.Getenv("MAIL_SERVER")
	fromMail = os.Getenv("FROM_MAIL")
	key = os.Getenv("SECRET_KEY")
}
