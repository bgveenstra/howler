package main  

import (
	"fmt";
    "log";
    "github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	} 

	// defer statement pushes a function call onto a list
	// list of saved calls is executed after the surrounding function returns
	// good for clean-up actions
	// here, closing the watcher after this main function is finished
	// Close removes all watches and closes the events channel. 
	defer watcher.Close()

	// https://golang.org/pkg/builtin/#make
	// make allocates an object of type slice, map, or chan
	// takes in a TYPE - here, a chan (Channel) of bools
	// returns the object of type TYPE 
	// and size given by second arg; here, channel is unbuffered
	// because no buffer capacity argument given

	done := make(chan bool)
	
	// https://tour.golang.org/concurrency/2
	// Channel (send and receives block)
	// 	ch <- v    // Send v to channel ch.
	// v := <-ch  // Receive from ch, and
	//            // assign value to v.

	// done is an unbuffered channel of booleans

	// invoking anonymous function in a goroutine
	// executes alongside the rest of this function's code
	go func() {
		log.Println("goroutine!")
		for { // while true
			select {
				// read from watcher.Events (channel of Events)
				case event := <-watcher.Events:
					// if there is an event to be had, log it
					log.Println("event:", event)
				// read from watcher.Errors (channel of errors)
				case err := <-watcher.Errors:
					// if there is an erorr waiting for us, log that
					log.Println("error:", err)
			}
		}
		// goroutine never ends!
		log.Println("goroutine gone?!")
	}() // calling anonymous function it immediately - required?
	// alternative would be invoking named function 
	// which can be defined outside this main function
	// can it be defined inside this function?
	// go doWhatever("arg")


	// add a file (directory, in this case) for watcher to watch (non-recursively)
	// is there any guarentee that this will happen before the coroutine starts?
	// or is order only okay because events will only start getting added to the channel
	// once the watcher has at least one directory to watch?
	err = watcher.Add("/tmp/")
	if err != nil {
		log.Println("erorr")
		log.Fatal(err)
	} else {
		fmt.Printf("Watching %s\n", "/tmp/")
	}

	// read from done...
	// but when is it getting written to?
	// it's not. 
	// reads are blocking, so it will just wait to read from done
	// until there is something in it or i kill the program
	x := <-done
	log.Println("done never able to read!", x)
}