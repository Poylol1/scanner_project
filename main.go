package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"

	//"os/exec"
	"github.com/otiai10/gosseract/v2"
	gomail "gopkg.in/mail.v2"
)

var GMAILKEY = os.Getenv("GMAILKEY")
var GPTKEY = os.Getenv("GPTKEY")

var LOW = gpio.Low
var HIGH = gpio.High

var t = time.NewTicker(50 * time.Microsecond)

func main() {
	_, er := host.Init()
	if er != nil {
		panic(er)
	}
	var motor_1_pulse = gpioreg.ByName("16")
	var motor_1_direction = gpioreg.ByName("20")
	var motor_1_enabled = gpioreg.ByName("21")
	//
	motor_1_enabled.Out(HIGH)
	motor_1_direction.Out(HIGH)
	for i := 0; i < 200; i++ {
		motor_1_pulse.Out(HIGH)
		<-t.C
		motor_1_pulse.Out(LOW)
	}
}

/*
ms - miliseconds.

direction - true is for right false for left.

ports - add the ports defineds by gpioreg.ByName("Number")
in order [direction,pulse,enabled]

If desire to go/move specific angle/distance change approach
*/
func rotate(ms int, direction bool, ports [3]gpio.PinIO) {
	pulse_per_rev := 200
	ports[2].Out(HIGH)
	ports[0].Out(LOW)
	if direction {
		ports[0].Out(HIGH)
	}
	for i := 0; i < ms*pulse_per_rev; i++ {
		ports[1].Out(HIGH)
		<-t.C
		ports[1].Out(LOW)
	}
	ports[2].Out(LOW)
}

func messageToGPT(message string) (output string, er error) {
	headers := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Please correct the following text. Try to avoid changing words. Correct the words, but not the sentence structure. If the text seems fine just return the text. Please correct all of them to the best of your capabilities",
			},
			{
				"role":    "user",
				"content": message,
			},
		},
	}
	data, er := json.Marshal(headers)
	if er != nil {
		fmt.Println(er)
		return "Failed Text Parsing JSON", er
	}

	req, er := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	if er != nil {
		fmt.Println(er)
		return "Failed To Parse Request", er
	}
	req.Header.Set("Authorization", ("Bearer " + GPTKEY))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, er := client.Do(req)
	if er != nil {
		fmt.Println(er)
		return "Client request error", er
	}

	defer req.Body.Close()

	bodyResp, er := io.ReadAll(resp.Body)
	if er != nil {
		fmt.Println(er)
		return "Failed reading response", er
	}
	type GPTResponse struct {
		Choices []struct {
			Message struct {
				Content string //"json: content"
			} //"json: message"
		} //"json: choices"
	}
	var jsonMap GPTResponse
	if er := json.Unmarshal(bodyResp, &jsonMap); er != nil {
		fmt.Println(er)
		return "JSON Decryption Failed", er
	}
	turner := jsonMap.Choices[0].Message.Content
	return turner, nil
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

// RASPBERRY PI PIN GPIO LAYOUT
//   1 . . 2
//   3 . . 4
//   5 . . 6
//   7 . . 8           - UART0 TX
//   9 . . 10
//  11 . . 12
//  13 . . 14
//  15 . . 16
//  17 . . 18
//  19 . . 20 SPI MOSI
//  21 . . 22 SPI MISO - SPI CS 0
//  23 . . 24 SPI CLK  - SPI CS 1
//  25 . . 26
//  27 . . 28
//  29 . . 30
//  31 . . 32
//  33 . . 34
//  35 . . 36
//  37 . . 38
//  39 . . 40
//
// 4. VCC Fan
// 6. GND Fan
// 8. Control Fan
// 34. Motor 1 GND
// 36: Motor 1 GPIO 16 Pulse - 200 P/rev
// 38: Motor 1 GPIO 20 Direction - 0 right  - 1 left ?
// 40: Motor 1 GPIO 21 Enable - 0 disabled -  1 enabled
//
//
//
