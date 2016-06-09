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

    inputChannel chan *message.Message
    outputChannel chan *message.Message
}

// New creates a new instance of the Client object.
func New(nick, pass string, ui ui.ThyrcUI) *Client {
	return &Client{nick: nick, pass: pass, ui: ui}
}

func (c *Client) Connect(inputChannel, outputChannel chan *message.Message) {
    c.inputChannel = inputChannel
    c.outputChannel = outputChannel

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

func (c *Client) Run() {
    for {
        msg := <-c.outputChannel
        fmt.Println(msg)
        // TODO: pass this to the right channel
    }
}
