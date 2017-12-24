package alert

import (
	"testing"
	"time"
	"fmt"
)

func TestStart(t *testing.T) {
	tick := time.Tick(5e8)
	send := make(chan bool)

	var stop = make(map[string]chan bool)
	stop["s1"] = send

	go func() {
		time.Sleep(2e9)
		stop["s1"] <- true
		fmt.Println(" Stop !")
	}()
	var a = false
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-send:
			fmt.Println("Stop !")
			a = true
		}
		if a {
			fmt.Println("Done !")
			break
		}
	}
	time.Sleep(4e9)
	fmt.Println(" Done !")
}
