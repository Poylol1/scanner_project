package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	gomail "gopkg.in/mail.v2"
	"os"
)

func main() {
	recieved, er := getTextFromImage("./tests/testImage.png")
	if er != nil {
		fmt.Println(er)
	}
	sendMail("Pollo", "poylolt@gmail.com", recieved, "This message was read from an image of a document with OCR!!")
}
func getTextFromImage(imagePath string) (text string, er error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(imagePath)
	text, er = client.Text()
	return text, er
}
func sendMail(From string, To string, Subject string, msg string) {
	// message
	message := gomail.NewMessage()

	//                          This should be an env variable
	message.SetHeader("From", "mesapidemo@gmail.com")
	//                      This would be inputted by user
	message.SetHeader("To", To)
	message.SetHeader("Subject", "[ From: "+From+"] "+Subject)

	message.SetBody("text/plain", msg)

	// I should understand what does this line does I hope it is not much problem
	//                                               This password should be an environment variable
	password := os.Getenv("GMAILKEY")
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "mesapidemo@gmail.com", password)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

}
