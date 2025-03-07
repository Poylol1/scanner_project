package main

import (
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func main() {
	_, er := host.Init()
	if er != nil {
		panic(er)
	}
	motor_1_pulse := gpioreg.ByName("16")
	motor_1_direction := gpioreg.ByName("20")
	motor_1_enabled := gpioreg.ByName("21")
	//turnOffPorts([]gpio.PinIO{motor_1_enabled, motor_1_direction, motor_1_pulse})
	rotate(30, true, [3]gpio.PinIO{motor_1_direction, motor_1_pulse, motor_1_enabled})
	pic, er := getPicture()
	if er != nil {
		panic(er)
	}
	text, er := getTextFromImage(pic)
	if er != nil {
		panic(er)
	}
	out, er := messageToGPT(text)
	if er != nil {
		panic(er)
	}
	processed := processText(out)
	sendMail(processed[0], processed[1], processed[2], processed[3])
	rotate(30, false, [3]gpio.PinIO{motor_1_direction, motor_1_pulse, motor_1_enabled})
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
// 40: Motor 1 GPIO 21 Enable - 1 disabled-  0 enabled
// Camera Connection
// Set the SPI Driver if you desire the SPI screen connection
