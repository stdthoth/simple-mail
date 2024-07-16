package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	emailPassword := os.Getenv("EMAIL_PASSWORD")
	yourMail := os.Getenv("MY_MAIL")
	recipient := "any email whatsoever"
	hostAddress := os.Getenv("HOST_ADDRESS")
	hostPort := os.Getenv("HOST_PORT")
	mailSubject := "Hey, I'm Just Checking On You."
	mailBody := "Hope you're doing okay! How are you doing today. "
	fullServerAddress := hostAddress + ":" + hostPort

	headerMap := make(map[string]string)
	headerMap["From"] = yourMail
	headerMap["To"] = recipient
	headerMap["Subject"] = mailSubject
	mailMessage := ""

	for k, v := range headerMap {
		mailMessage += fmt.Sprintf("%s: %s\\r", k, v)
	}

	mailMessage += "\\r" + mailBody

	authenticate := smtp.PlainAuth("", yourMail, emailPassword, hostAddress)
	tlsConfigurations := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         hostAddress,
	}

	conn, err := tls.Dial("tcp", fullServerAddress, tlsConfigurations)
	if err != nil {
		log.Panic(err)
	}

	newClient, err := smtp.NewClient(conn, hostAddress)
	if err != nil {
		log.Panic(err)

	}

	// Auth
	if err = newClient.Auth(authenticate); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = newClient.Mail(yourMail); err != nil {
		log.Panic(err)

	}

	if err = newClient.Rcpt(recipient); err != nil {
		log.Panic(err)
	}

	writer, err := newClient.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = writer.Write([]byte(mailMessage))
	if err != nil {
		log.Panic(err)
	}

	err = writer.Close()
	if err != nil {
		log.Panic(err)
	}

	err = newClient.Quit()
	if err != nil {
		fmt.Println("THERE WAS AN ERROR")
	}

	fmt.Println("Successful, the mail was sent!")

}
