package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort	string `mapstructure:"APP_PORT"`
	EmailFrom string `mapstructure:"EMAIL_FROM"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`
	EmailRecipient string `mapstructure:"EMAIL_RECIPIENT"`
	EmailSubject string `mapstructure:"EMAIL_SUBJECT"`
	SmtpHost string `mapstructure:"SMTP_HOST"`
	SmtpPort string `mapstructure:"SMTP_PORT"`
}

type Request struct {
	Recipients 	string `json:"recipients"`  
	Subject 	string `json:"subject"` 
    Message 	string `json:"message"`
}

func loadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return
}

func main() {
	// load config
	config, err := loadConfig(".env")
	if err != nil {
		fmt.Println("Error loading config: ", err.Error())
	}
	
	// create new fiber app, user cors, and logger
    app := fiber.New()
	app.Use(cors.New())
    app.Use(logger.New())
	
	// api endpoint handler
    app.Post("/api/send", func (c *fiber.Ctx) error {
		// parse request body
		request := new(Request)

		if err := c.BodyParser(request); err != nil {
			fmt.Println("Error parsing request body: ",err)
			return c.SendStatus(400)
		}

		// call function to handle the email sending
        err := sendEmail(config, request.Subject, request.Message, request.Recipients)
        if err != nil {
			fmt.Println("Error sending email: ", err.Error())
            return c.SendStatus(400)
        }
        return c.SendStatus(200)
    })

	// listen on port that comes from .env
	log.Fatal(app.Listen(config.AppPort))
}

// this function takes a message body and a list of email recipients
// sets up the smtp, auth and message and tries to send an email
// it returns an error
func sendEmail(config Config, emailSubject string, messageBody string, recipients string) error {
	// set up from/to and app password
	from := config.EmailFrom
	password := config.EmailPassword

	// if list is empty, get email recipient from .env
	var recipientList []string
	if(len(recipients) == 0) {
		recipientList = strings.Split(config.EmailRecipient, ",")
	} else {
		recipientList = strings.Split(recipients, ",")	
	}

	// smtp server setup
	host := config.SmtpHost
	port := config.SmtpPort
	address := host + ":" + port

	// is subject is empty, get it from .env
	var subject string
	if(len(emailSubject) == 0) {
		subject = config.EmailSubject
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