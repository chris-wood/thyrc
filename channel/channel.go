package channel

import (
    // "time"
    // "net"
    // "fmt"
    // "github.com/chris-wood/thyrc/message"
)

type Channel struct {
    Name string
    inputChannel chan string
    outputChannel chan string
}

// New creates a new instance of the Server object.
func New(name string, input, output chan string) *Channel {
    return &Channel{Name: name, inputChannel: input, outputChannel: output}
}
