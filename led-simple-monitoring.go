package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var pinsIndex = []int{14, 15, 23, 18, 4, 24} //pattern is red, green, red...
var Leds [3]LED                              //array to hold all leds
//add websites to monitor here
var sites = [3]string{
	"https://alonsoarteaga.com",
	"https://alonsoarteaga.me",
	"http://google.com/example",
}

//LED has two rpio.pin types, green and blue, respectively
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
	for i := 0; i < len(Leds); i++ {
		Leds[i].green = rpio.Pin(pinsIndex[counter])
		Leds[i].red = rpio.Pin(pinsIndex[counter+1])
		counter = counter + 2
	}

	go StartChecks() // main

	//listen for interrupt and teardown (turn off leds)
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...\n")
		/** turn off leds on exit */
		for i := 0; i < len(Leds); i++ {
			Leds[i].green.Low()
			Leds[i].red.Low()
		}

		close(cleanupDone)
	}()
	<-cleanupDone
}

func StartChecks() {
	for counter, site := range sites {
		fmt.Printf("Checking site %v\n", site)
		GetReturnCode(site, Leds[counter])
	}
}

func GetReturnCode(site string, led LED) {
	//site one
	resp, err := http.Get(site)
	if resp.StatusCode >= 400 {
		// handle error, RED on
		led.green.Low()
		led.red.High()
		log.Printf("site %s unreachable, error is %v\n", site, err)
		defer resp.Body.Close()
		return
	}
	defer resp.Body.Close()

	log.Printf("site %s returned %v OK\n", site, resp.StatusCode)
	led.red.Low()
	led.green.High()
}
