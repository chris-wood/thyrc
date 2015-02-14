package main

import (
    "bufio"
    "log"
    "net"
    "os"
    "fmt"
)

type Command struct {
    command string
    parameters []string
}

func parseCommandString(commandString string) (Command) {
    var command Command
    return command
}

func stopSession(connection net.Conn) {
    err := connection.Close()
    if err != nil {
        log.Fatal(err)
    }
}

func prompt(connection net.Conn) (bool, error) {
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')

        var _ = parseCommandString(text)
    }
}

func startSession(serverAddress string) {
    fmt.Println("Connecting to " + serverAddress)
    connection, err := net.Dial("tcp", serverAddress)
    if err != nil {
        log.Fatal(err)
    }

    go prompt(connection)
}

func main() {
    args := os.Args[1:]
    log.SetFlags(log.Lshortfile)
    go startSession(args[0])    
}
