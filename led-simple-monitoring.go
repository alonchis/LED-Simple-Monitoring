package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	// BCM PINS on the board.
	pinsIndex = []int{14, 15, 23, 18, 4, 24} //LED1 green, LED1 red, LED2 green, ...

	//array to store all leds
	led [3]LED

	//Signal channel to listen for sys interrupt (CTRL+C)
	signalChan = make(chan os.Signal, 1)

	//add websites to monitor here
	sites = [3]string{
		"https://google.com",
		"https://wikipedia.org",
		"https://github.com/notfound",
	}

	period = time.Minute * 5 // changes how often to check status
)

//LED for two rpio.pin types, green and red, respectively
type LED struct {
	green rpio.Pin
	red   rpio.Pin
}

func main() {
	/** Open and map memory to access gpio, check for errors */
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	/** Unmap gpio memory when done */
	defer rpio.Close()

	//Init all leds
	// probably could be written better
	counter := 0
	for i := 0; i < len(led); i++ {
		led[i].green = rpio.Pin(pinsIndex[counter])
		led[i].red = rpio.Pin(pinsIndex[counter+1])
		counter = counter + 2
	}

	// Run go routines
	for i := 0; i < len(sites); i++ {
		go GetReturnCode(sites[i], led[i])
	}

	//listen for interrupt and teardown (turn off leds)
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...")
		/** turn off leds on exit */
		for i := 0; i < len(led); i++ {
			led[i].green.Low()
			led[i].red.Low()
		}

		close(cleanupDone)
	}()
	<-cleanupDone
}

func GetReturnCode(site string, led LED) {
	for {
		select {
		case <-signalChan: // exits on interrupt
			return
		default:
			resp, _ := http.Head(site)
			if resp.StatusCode >= 400 {
				// handle error, RED on
				led.green.Low()
				led.red.High()
				log.Printf("%s unreachable, error code %v\n", site, resp.StatusCode)
				time.Sleep(period)
				continue
			}

			log.Printf("%s returned %v OK\n", site, resp.StatusCode)
			led.red.Low()
			led.green.High()
			time.Sleep(period)
		}
	}
}
