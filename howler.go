package main

import (
	"fmt"
	"github.com/bgveenstra/slacker"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

// the "howler" command
func main() {
	// @TODO - add better CLI
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

var verbose = false

func debugLog(label string, message string) {
	if verbose {
		log.Printf("%s: %s", label, message)
	}
}

// wrap slacker.PostSlackMessage, passing environment variable
func slack(message string) error {
	return slacker.PostSlackMessage(message, os.Getenv("HOWLER_SLACK_WEBHOOK_URL"))
}

func WatchDirForever(dir string) error {
	// @TODO watch types should come from flag or config - some arg
	// @TODO configure messages elsewhere for organization?
	eventWatchTypes := make(map[string]string)
	eventWatchTypes["CREATE"] = "created"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				debugLog("event", event.String())

				// @TODO - move all or part of message generation to helper function
				verb, isWatchType := eventWatchTypes[event.Op.String()]
				if isWatchType {
					message := fmt.Sprintf("File %s: %s\n", verb, event.Name)
					debugLog("message", message)
					err := slack(message)
					if err != nil {
						log.Fatal(err)
					}
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	// @TODO - can move this message generation elsewhere
	watchDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Watching %s", watchDir)
	debugLog("watch message", message)
	slack(message)

	// @TODO - replace with waitgroup
	done := make(chan bool)
	<-done

	return nil
}
