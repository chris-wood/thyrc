package client

import (
    "fmt"
    "github.com/chris-wood/thyrc/message"
    "github.com/chris-wood/thyrc/ui"
)

type Client struct {
    nick string
    pass string
    ui ui.ThyrcUI
}

// New creates a new instance of the Client object.
func New(nick, pass string, ui ui.ThyrcUI) *Client {
	return &Client{nick: nick, pass: pass, ui: ui}
}

func (c *Client) connect(inputChannel, outputChannel chan *message.Message) {
    passMessage := message.Parse("PASS " + c.pass)
    nickMessage := message.Parse("NICK " + c.nick)
    userMessage := message.Parse("USER blahblah blah blah blah")

    inputChannel <- passMessage
    passResponse := <-outputChannel
    fmt.Println(passResponse)

    inputChannel <- nickMessage
    nickResponse := <-outputChannel
    fmt.Println(nickResponse)

    inputChannel <- userMessage
    userResponse := <-outputChannel
    fmt.Println(userResponse)

    fmt.Println("Sent connection parameters")
}

func (c *Client) Run(inputChannel, outputChannel chan *message.Message) {
    c.connect(inputChannel, outputChannel)
    shutdownChannel := make(chan int)
    go c.handleServerMessages(outputChannel, shutdownChannel)
    c.handleClientMessages(inputChannel, shutdownChannel)
}

func (c *Client) handleServerMessages(outputChannel chan *message.Message, shutdownChannel chan int) {
    for {
        select {
        case serverMessage := <-outputChannel:
            fmt.Println(serverMessage)
            // TODO: this should be passed to the UI
        case _ = <-shutdownChannel:
            return
        }
    }
}

func (c *Client) handleClientMessages(inputChannel chan *message.Message, shutdownChannel chan int) {
    uiChannel := c.ui.GetInputChannel()
    for {
        fmt.Println("trying to read from the UI")
        msgString := <-uiChannel
        fmt.Println("Read " + msgString + " from the server!")

        if (msgString == "/quit") {
            shutdownChannel <- 1
            return;
        }

        fmt.Println("Sending '" + msgString + "' to server")

        msg := message.Parse(msgString)
        inputChannel <- msg
    }
}
