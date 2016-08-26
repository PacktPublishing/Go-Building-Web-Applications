package main

import (
    "github.com/howeyc/fsnotify"
)

func main() {

    scriptDone := make(chan bool)
    dirSpy, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        for {
            select {
            case fileChange := <-dirSpy.Event:
                log.Println("Something happened to a file:", fileChange)
            case err := <-dirSpy.Error:
                log.Println("Error with fsnotify:", err)
            }
        }
    }()

    err = dirSpy.Watch("/mnt/sharedir")
    if err != nil {
      fmt.Println(err)
    }

    <-scriptDone


    dirSpy.Close()
}