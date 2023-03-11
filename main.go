package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

type Request struct {
	Recipients 	string `json:"recipients"`  
	Subject 	string `json:"subject"` 
    Message 	string `json:"message"`
}

func main() {
	// if in dev load config from .env
	env := os.Getenv("ENV")
	if env == "dev" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error loading config: ", err.Error())
			os.Exit(1)
		}
	}
	
	
	// create new fiber app, user cors, and logger
    app := fiber.New()
	app.Use(cors.New())
    app.Use(logger.New())

	app.Get("/", func (c *fiber.Ctx) error {
		return c.SendString("Serivce is online!")
	})
	
	// api endpoint handler
    app.Post("/api/send", func (c *fiber.Ctx) error {
		// parse request body
		request := new(Request)

		if err := c.BodyParser(request); err != nil {
			fmt.Println("Error parsing request body: ",err)
			return c.SendStatus(400)
		}

		// call function to handle the email sending
        err := sendEmail(request.Subject, request.Message, request.Recipients)
        if err != nil {
			fmt.Println("Error sending email: ", err.Error())
            return c.SendStatus(400)
        }
        return c.SendStatus(200)
    })

	// listen on port that comes from .env
	log.Fatal(app.Listen("0.0.0.0:" + os.Getenv("APP_PORT")))
}

// this function takes a message body and a list of email recipients
// sets up the smtp, auth and message and tries to send an email
// it returns an error
func sendEmail(emailSubject string, messageBody string, recipients string) error {
	// set up from/to and app password
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	// if list is empty, get email recipient from .env
	var recipientList []string
	if(len(recipients) == 0) {
		recipientList = strings.Split(os.Getenv("EMAIL_RECIPIENT"), ",")
	} else {
		recipientList = strings.Split(recipients, ",")	
	}

	// smtp server setup
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	address := host + ":" + port

	// is subject is empty, get it from .env
	var subject string
	if(len(emailSubject) == 0) {
		subject = os.Getenv("EMAIL_SUBJECT")
	} else {
		subject = fmt.Sprintf("Subject: %v\n", emailSubject)
	}

	message := []byte(subject + messageBody)

	// smtp auth
	auth := smtp.PlainAuth("", from, password, host)

	// send email
	err := smtp.SendMail(address, auth, from, recipientList, message)
	if err != nil {
		return err
	}

    return nil
}