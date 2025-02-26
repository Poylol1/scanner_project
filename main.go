package main

import (
	"fmt"
	"net/http"
	"os"

	//"os/exec"
	"github.com/otiai10/gosseract/v2"
	gomail "gopkg.in/mail.v2"
)

var GMAILKEY = os.Getenv("GMAILKEY")
var GPTKEY = os.Getenv("GPTKEY")

func main() {
	recieved, er := getTextFromImage("./tests/extendedTest.jpg")
	if er != nil {
		fmt.Println(er)
	}
	fmt.Println(recieved)
}

func messageToGPT(scanned string) (output string, er error) {
	req, er = http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", map[string]inteface{}{"model": "gpt-4o-mini", "Authorization": ""})
	// httpClient := http.Client{:
	// req, _ = := htthttpClient.
	return "TEMP", nil
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
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "mesapidemo@gmail.com", GMAILKEY)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

}
