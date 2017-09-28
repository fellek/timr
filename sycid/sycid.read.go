// env GOOS=linux GOARCH=arm GOARM=7 go build
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jteeuwen/evdev"
)

var devName = "/dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd"

var uri = "/api/cardreader?id="

func check(e error) error {
	if e != nil {
		log.Fatal(e.Error())
	}
	return e
}

func readRfid(dev *evdev.Device) string {
	var rfid string

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	done := false
	for !done {
		var key byte

		select {
		case <-signals:
			return ""

		case evt := <-dev.Inbox:
			if evt.Type != evdev.EvKeys {
				continue // Not a key event.
			}

			if evt.Value == 0 {
				continue // Key is released -- we want key presses.
			}

			if evt.Code == 28 {
				done = true // end of transmission
			}

			if !done {
				key = byte((evt.Code - 1) % 10)
				log.Debugf("found EV_KEY: %d", key)
				rfid += strconv.Itoa(int(key))
				log.Debugf("concat RFID: %s", string(rfid))
			}
		}
	}
	return rfid
}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func sendToServer(rfid, server, port string) error {
	url := "http://" + server + ":" + port +
		uri + rfid //DevSkim: ignore DS137138 until 2017-10-25
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	lstr := fmt.Sprintf("body: [%s]", body)
	log.Debug(lstr)

	return err
}

func main() {
	argLogging := flag.String("log", "debug", "one of: debug, info, warn, error, fatal")
	argServer := flag.String("server", "localhost", "the name of the server, the rfid gets send to")
	argServerPort := flag.Int("port", 8000, "the portnumber of the 'server'")
	argUsage := flag.Bool("h", false, "writes some information about args")

	// Once all flags are declared, call flag.Parse() to execute the command-line parsing.
	flag.Parse()

	if *argUsage {
		flag.Usage()
		os.Exit(0)
	}
	log.SetHandler(cli.Default)
	log.SetLevelFromString(*argLogging)
	log.Info("syscid RFID start")
	log.Infof("starting with server %s, port %d, log %s. -h for help", *argServer, *argServerPort, *argLogging)
	dev, err := evdev.Open(devName)
	check(err)
	defer dev.Close()

	for {
		log.Info("Waiting for data")
		rfid := readRfid(dev)
		if len(rfid) > 0 {
			log.Info("rfid detected: " + rfid)
			err = sendToServer(rfid, *argServer, strconv.Itoa(*argServerPort))
			check(err)
		}
	}
}
