package main

import (
    // "time"
    "log"
    "os"
    "fmt"
    "github.com/chris-wood/thyrc/client"
)

func main() {
    file, err := os.OpenFile("thyrc.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer file.Close()
    log.SetOutput(file)

    fmt.Fprintln(os.Stderr, "Setup the log")

    args := os.Args[1:]
    if len(args) < 1 {
        fmt.Println("usage: go run client.go <server:port>")
        log.Fatalf("usage")
    }

    // err = runSession(args[0])
    // if err != nil {
    //     fmt.Fprintln(os.Stderr, "Failed.")
    //     fmt.Fprintln(os.Stderr, "Error: " + string(err.Error()))
    // }
}
