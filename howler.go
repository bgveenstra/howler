package main  

import (
	// "fmt";
    "log"
    "github.com/fsnotify/fsnotify"
    "github.com/bgveenstra/slacker"
    "fmt"
    "os"
    "path/filepath"
)
// https://github.com/fsnotify/fsnotify/blob/master/example_test.go
func main() {
	log.Println("begin main")
	watchDirArg := os.Args[1]

	WatchDirForever(watchDirArg)
	log.Println("end main")
}


func WatchDirForever(dir string){

	eventWatchTypes := make(map[string] bool)
	// watch types should come from flag or config - some arg
	eventWatchTypes["CREATE"] = true


	watchDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}	
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	} 

	// defer statement pushes a function call onto a list
	// list of saved calls is executed after the surrounding function returns
	// good for clean-up actions
	// here, closing the watcher after this WatchDirForever function is finished
	// Close removes all watches and closes the events channel. 
	defer watcher.Close()



	log.Println("before")
	// invoking anonymous function in a goroutine
	// executes alongside the rest of this function's code
	go func() {
		log.Println("goroutine!")
		for { // while true
			select {
				// read from watcher.Events (channel of Events)
				case event := <-watcher.Events:
					// if there is an event to be had, log it
					log.Println("event name:", event.Name)
					log.Println("event op:", event.Op)
					if eventWatchTypes[event.Op.String()] {
						slacker.PostSlackMessage(fmt.Sprintf("Holwer Update \n%s %s", event.Op, event.Name))
					}
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
	// which can be defined outside this WathDirForever function
	// can it be defined inside this function?
	// go doWhatever("arg")

	log.Println("after")


	// add a file (directory, in this case) for watcher to watch (non-recursively)
	// is there any guarentee that this will happen before the coroutine starts?
	// or is order only okay because events will only start getting added to the channel
	// once the watcher has at least one directory to watch?
	err = watcher.Add(dir)
	if err != nil {
		log.Println("erorr")
		log.Fatal(err)
	} else {
		log.Println("Watching", watchDir)
		slacker.PostSlackMessage(fmt.Sprintf("Watching %s", watchDir))
	}


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


	// read from done...
	// but when is it getting written to?
	// it's not. 
	// reads are blocking, so it will just wait to read from done
	// until there is something in it or i kill the program
	x := <-done
	log.Println("done never able to read!", x)
}