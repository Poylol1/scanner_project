package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	gomail "gopkg.in/mail.v2"
	"io"
	"net/http"
	"os"
	"os/exec"
	"periph.io/x/conn/v3/gpio"
	"time"
)

var GMAILKEY = os.Getenv("GMAILKEY")
var GPTKEY = os.Getenv("GPTKEY")
var HOME = os.Getenv("HOME")

var LOW = gpio.Low
var HIGH = gpio.High

var t = time.NewTicker(50 * time.Microsecond)

/*
degrees - 360th of a rotation.

direction - true is for right false for left.

ports - add the ports defineds by gpioreg.ByName("Number")
in order [direction,pulse,enabled]

If desire to go/move specific angle/distance change approach
*/
func rotate(degrees int, direction bool, ports [3]gpio.PinIO) {
	pulse_per_rev := 200 // TODO FIX THISj
	ports[2].Out(LOW)
	ports[0].Out(LOW)
	if direction {
		ports[0].Out(HIGH)
	}
	fmt.Print("I entered the for loop")
	for i := 0; i < degrees*pulse_per_rev/360; i++ {
		//	fmt.Print("I passed here")
		ports[1].Out(HIGH)
		<-t.C
		ports[1].Out(LOW)
		<-t.C
	}
	fmt.Print("I leaved the for loop")
	ports[0].Out(LOW)
	ports[1].Out(LOW)
	ports[2].Out(HIGH)
}

func turnOffPorts(ports []gpio.PinIO) {
	ports[1].Out(HIGH)
	for i := 1; i < len(ports); i++ {
		ports[i].Out(LOW)
	}
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

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "mesapidemo@gmail.com", GMAILKEY)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

}

func getPicture() (imagePath string, er error) {
	cmd := exec.Command("echo \"${python " + HOME + "/Projects/scanner_project/pythonWrapper/camera.py\"}")
	out, er := cmd.Output()
	if er != nil {
		fmt.Println("Could not run: ", er)
		return "Error", er
	}
	return string(out), nil
}
func processText(text string) []string {
	turner := make([]string, 0, 4)
	subjectText := []int{0, 0}
	toText := []int{0, 0}
	bodyText := []int{0, 0}
	fromText := []int{0, 0}

out:
	for i := range text {
		if text[i:i+2] == "To" {
			toText[0] = i
		}
		if text[i:i+4] == "From" {
			toText[1] = i - 1
			fromText[0] = i
		}

		if text[i:i+7] == "Subject" {
			fromText[1] = i - 1
			subjectText[0] = i

		}

		if text[i:i+4] == "Body" {
			subjectText[1] = i - 1
			bodyText[0] = i
			break out
		}
	}
	from := text[fromText[0]:fromText[1]]
	to := text[toText[0]:toText[1]]
	subject := text[subjectText[0]:subjectText[1]]
	body := text[bodyText[0]:bodyText[1]]
	turner[0] = from
	turner[1] = to
	turner[2] = subject
	turner[3] = body
	return turner
}
