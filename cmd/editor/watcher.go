package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os/exec"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan error)

	log.Println("Wait....")
	go monitor(watcher, done)

	listFiles := make([]string, 2)
	listFiles[0] = "../../cmd/editor"
	listFiles[1] = "../static/templates/entry"

	for _, elm := range listFiles {
		err = watcher.Add(elm)
		if err != nil {
			return err
		}
	}

	return <-done
}

func monitor(watcher *fsnotify.Watcher, done chan error) {
	for {
		select {
		case event := <-watcher.Events:
			go notify(event)
		case err := <-watcher.Errors:
			done <- err
			return
		}
	}
}

var lock = false

func notify(event fsnotify.Event) {

	if lock {
		return
	}

	lock = true
	defer func() {
		lock = false
	}()

	if event.Op&fsnotify.Write == fsnotify.Write {
		log.Println(event.Name, event)
		log.Println("********************** deploy")
		command()
		log.Println("********************** end")
	}
	log.Println("Wait....")
	return
}

func command() {

	out, err := exec.Command("go", "run", "deploy.go").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("")
	fmt.Println("********************** command output")
	fmt.Println(string(out))
	fmt.Println("**********************")
}
