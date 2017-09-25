// env GOOS=linux GOARCH=arm GOARM=5 go build
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/dddpaul/golang-evdev/evdev"
)

var keys = "X^1234567890XXXXqwertzuiopXXXXasdfghjklXXXXXyxcvbnmXXXXXXXXXXXXXXXXXXXXXXX"
var devName = "/dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd"

var server = "localhost:8080"
var uri = "/api/cardreader?id="

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
		panic(e)
	}
}

func readRfid(fd *evdev.InputDevice) string {
	// reads rfid in the format:
	// 0000802843
	// 0000802843
	var events []evdev.InputEvent

	// var rfid bytes.Buffer

	var rfid string
	// buf := make([]byte, 1024)
	var err error

	for {
		var done = false
		var key byte
		// bl := 0

		events, err = fd.Read()
		check(err)
		// var len = len(rfid)
		for i := range events {

			log.Trace(events[i].String())

			switch events[i].Type {
			case evdev.EV_KEY: //  0x01
				log.Trace(events[i].String())
				if events[i].Code == 28 {
					done = true
				}

				if !done {
					if events[i].Value == 1 {
						key = byte((events[i].Code - 1) % 10)
						log.Debugf("found EV_KEY: %d", key)
						rfid += strconv.Itoa(int(key))
						log.Debugf("concat RFID: %s", string(rfid))
					}
				}
			case evdev.EV_SYN: // 0x00
				fallthrough
			case evdev.EV_REL: // 0x02
				fallthrough
			case evdev.EV_ABS: // 0x03
				fallthrough
			case evdev.EV_MSC: // 0x04
				fallthrough
			case evdev.EV_SW: // 0x05
				fallthrough
			case evdev.EV_LED: // 0x11
				fallthrough
			case evdev.EV_SND: // 0x12
				fallthrough
			case evdev.EV_REP: // 0x14
				fallthrough
			case evdev.EV_FF: // 0x15
				fallthrough
			case evdev.EV_PWR: // 0x16
				fallthrough
			case evdev.EV_FF_STATUS: //  0x17
				fallthrough
			case evdev.EV_MAX: // 0x1f
				fallthrough
			default:
				// log.Debug("Type not found")
			}
		}

		if done {
			break
		}
	}
	return string(rfid)
}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func sendToServer(rfid string) error {
	url := "http://" + server + uri + rfid //DevSkim: ignore DS137138 until 2017-10-25
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("get:\n", keepLines(string(body), 3))
	return err
}

func main() {
	var dev *evdev.InputDevice
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	log.Info("syscid RFID start")
	dev, err := evdev.Open(devName)
	check(err)
	// defer evdev.

	for {
		rfid := readRfid(dev)
		if len(rfid) > 0 {
			log.Info("rfid detected: " + rfid)
			err = sendToServer(rfid)
			check(err)
		}

	}
}
