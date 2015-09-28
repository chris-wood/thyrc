package main

import (
    "time"
    "log"
    "net"
    "os"
    "fmt"
    "strings"
)

type Command struct {
    command string
    parameters []string
}

func handleCommand(command Command, connection net.Conn) (bool) {
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

func serverReadAndWrite(channelFromServer chan string, channelToServer chan string, connection net.Conn) {
    reply := make([]byte, 1024)
    stayAlive := true

    for ; stayAlive ; {
        stayAlive = true // TODO: read from concurrent channel here
        _, err := connection.Read(reply)
        if err != nil {
            fmt.Println("Write to server failed:", err.Error())
            return
        }

        // channelFromServer <- strings.TrimSpace(string(reply))
        response := strings.TrimSpace(string(reply))
        log.Println(response)
        fmt.Println(response)
        select {
            case msgToSend, ok := <-channelToServer:
                if ok {
                    log.Println("Sending: " + msgToSend)
                    rawBytes := []byte(msgToSend)
                    connection.Write(rawBytes)
                } else {
                    time.Sleep(time.Second)
                }
            default:
                continue
            }

        // // Shove the output to the main view
        // view, err := gui.View("main")
        // if err != nil {
        //     log.Fatal("Could not recover handle to the main view.")
        // }
    }
}

func ircHandler(channelFromServer chan string, channelToServer chan string) {
    alive := true
    msgToSend := ""
    connected := false

    for ; alive ; {

        if !connected {
            channelToServer <- "PASS none\n"
            channelToServer <- "NICK random\n"
            channelToServer <- "USER rawrrawr blah blah blah\n"
            connected = true // don't connect again
        } else {
            // Read and display the server response
            response := <-channelFromServer
            fmt.Println(response)

            // Read from stdin
            fmt.Scanln(msgToSend)
            channelToServer <- msgToSend

            // Reset the message to send
            msgToSend = ""
        }
    }

    // stringReply := strings.TrimSpace(string(reply))
    // if len(stringReply) > 0 {
    //     // fmt.Println(stringReply)
    //     // view.Clear()
    //     fmt.Fprintln(view, stringReply)
    //     gui.Flush()
    // }
}

func runSession(serverAddress string) (error) {
    connection, err := net.Dial("tcp", serverAddress)
    if err != nil {
        return err;
    }

    channelFromServer := make(chan string)
    channelToServer := make(chan string)
    killChannel := make(chan int)

    go serverReadAndWrite(channelFromServer, channelToServer, connection)
    go ircHandler(channelFromServer, channelToServer)

    <-killChannel // block until killed
    stopSession(connection)

    return nil
}

func main() {
    file, err := os.OpenFile("client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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

    err = runSession(args[0])
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed.")
        fmt.Fprintln(os.Stderr, "Error: " + string(err.Error()))
    }
}
