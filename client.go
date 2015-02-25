package main

import (
    "bufio"
    "log"
    "net"
    "os"
    "fmt"
    "strings"
//    "github.com/jroimartin/gocui"
)

//PASS none
//NICK sorandom29      
//USER blah blah blah blah

type Command struct {
    command string
    parameters []string
}

func handleCommand(command Command) (bool) {
    if command.command == "quit" {
        return false 
    } else {
        fmt.Println(command.command)
        return true
    }
}

func parseCommandString(commandString string) (Command) {
    var command Command
    stringFields := strings.Fields(commandString)
    if len(stringFields) > 0 {
        command.command = stringFields[0]
        if len(stringFields) > 1 {
            command.parameters = stringFields[1:]    
        } else{
            command.parameters = nil
        }
    } else {
        // error?
    }
    
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

        command := parseCommandString(text)
        if !handleCommand(command) {
            break
        }
    }
    return true, nil
}

func readFromServer(connection net.Conn) {
    reply := make([]byte, 1024)
    stayAlive := true 
    for ; stayAlive ; {
        stayAlive = true // TODO: read from concurrent channel here
        var _, err = connection.Read(reply)
        if err != nil {
            fmt.Println("Write to server failed:", err.Error())
            return
        }
        fmt.Println("$> " + string(reply))
    }
}

func startSession(serverAddress string) {
    fmt.Println("Connecting to " + serverAddress)
    connection, err := net.Dial("tcp", serverAddress)
    if err != nil {
        log.Fatal(err)
    }

    go readFromServer(connection)
    prompt(connection)
}

func main() {
    args := os.Args[1:]
    if len(args) < 1 {
        fmt.Println("usage: go run client.go <server:port>")
        return 
    }
    log.SetFlags(log.Lshortfile)
    startSession(args[0])    
}
