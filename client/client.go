package client

import (
    "fmt"
    "github.com/chris-wood/thyrc/message"
)

type Client struct {
    inputChannel chan *message.Message
    outputChannel chan *message.Message
}

// New creates a new instance of the Client object.
func New() *Client {
	return &Client{}
}

func (c *Client) Connect(inputChannel, outputChannel chan *message.Message) {
    c.inputChannel = inputChannel
    c.outputChannel = outputChannel

    passMessage := message.Parse("PASS dsasdasdas")
    nickMessage := message.Parse("NICK asdasdasdad")
    userMessage := message.Parse("USER blahblah blah blah blah")

    inputChannel <- passMessage
    inputChannel <- nickMessage
    inputChannel <- userMessage

    fmt.Println("Sent connection parameters")
}

func (c *Client) Run() {
    for {
        msg := <-c.outputChannel
        fmt.Println(msg.Encode())
    }
}
