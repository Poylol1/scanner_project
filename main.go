package main

import (
	"fmt"
	gomail "gopkg.in/mail.v2"
)

func main() {
	// message
	message := gomail.NewMessage()

	message.SetHeader("From", "mesapidemo@gmail.com")
	message.SetHeader("To", "poylolt@gmail.com")
	message.SetHeader("Subject", "[From: Name goes here ] This message was sent through a Raspberry Pi*")

	message.SetBody("text/plain", "Yey it worked!")

	// I should understand what does this line does I hope it is not much problem
	//                                               This password should be an environment variable
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "mesapidemo@gmail.com", "vqzwgroptpdisxwb")

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
