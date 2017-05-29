package main

import (
	"os"
	"path/filepath"
  "log"
  "fmt"
  "github.com/fsnotify/fsnotify"
  "github.com/bgveenstra/slacker"
)


// handles the "howler" command
func main() {
	// @TODO - add better help/usage information (--help flag?)
	numArgs := len(os.Args)
	if numArgs < 2 {
		log.Fatal("Error - Too Few Arguments: howler requires a directory or file name\n usage: howler /tmp")
	} else if numArgs > 2 {
		log.Fatal("Error - Too Many Arguments: howler only accepts one directory or file name\n usage: howler /tmp")
	}
	watchDirArg := os.Args[1]
	err := WatchDirForever(watchDirArg)
	log.Fatal(err)
}


func WatchDirForever(dir string) error {

	// @TODO watch types should come from flag or config - some arg
	// @TODO configure messages elsewhere for organization?
	eventWatchTypes := make(map[string] string)
	eventWatchTypes["CREATE"] = "created"

	watchDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// @TODO - enable logging if "verbose" mode
				// log.Println("event:", event)
				// @TODO - move all or part of message generation to helper function
				verb, isWatchType := eventWatchTypes[event.Op.String()]
				if isWatchType {
					slacker.PostSlackMessage(fmt.Sprintf("File %s: %s\n", verb, event.Name))
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return err
	log.Println("Watching", watchDir)
	// @TODO - also move this message generation elsewhere
	slacker.PostSlackMessage(fmt.Sprintf("Watching %s", watchDir))

	// @TODO - replace with waitgroup
	done := make(chan bool)
	<-done

	return nil
}
