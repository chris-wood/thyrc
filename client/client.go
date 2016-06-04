package client

import (
    // "time"
    "log"
    "net"
    "fmt"
    "strings"
)

type Client struct {
}

// New creates a new instance of the Client object.
func New() *Client {
	return &Client{}
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
            fmt.Println("Read from server failed:", err.Error())
            return
        }

        channelFromServer <-strings.TrimSpace(string(reply))
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
                    // time.Sleep(time.Second)
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
