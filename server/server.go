package server

import (
    // "time"
    "net"
    // "fmt"
    "github.com/chris-wood/thyrc/channel"
    "github.com/chris-wood/thyrc/message"
)

type Server struct {
    address string
    connection net.Conn
    inputChannel chan *message.Message
    outputChannel chan *message.Message
    channels []*channel.Channel
}

// New creates a new instance of the Server object.
func New(address string) *Server {
    connection, err := net.Dial("tcp", address)
    if err != nil {
        return nil;
    }

	return &Server{address: address, connection: connection}
}

func (s *Server) MakeChannels() (chan *message.Message, chan *message.Message) {
    s.inputChannel = make(chan *message.Message)
    s.outputChannel = make(chan *message.Message)

    // Kick off the read and write functions
    go s.Write()
    go s.Read()

    return s.inputChannel, s.outputChannel
}

func (s *Server) Write() error {
    for {
        message := <-s.inputChannel
        encodedMessage := message.Encode()
        encodedMessageBytes := []byte(encodedMessage)

        _, err := s.connection.Write(encodedMessageBytes)
        if err != nil {
            // TODO: do what?
        }
    }
}

func (s *Server) Read() {
    reply := make([]byte, 1024)
    for {
        _, err := s.connection.Read(reply)
        if err != nil {
            // TODO: do what?
        }

        msg := message.Parse(string(reply))
        s.outputChannel <- msg
    }
}

func (s *Server) Join(channelName string) *channel.Channel {
    input := make(chan string)
    output := make(chan string)
    channel := channel.New(channelName, input, output)

    s.channels = append(s.channels, channel)

    return channel
}

func (s *Server) deleteChannel(channelIndex int) {
    s.channels = append(s.channels[:channelIndex], s.channels[channelIndex+1:]...)
}

func (s *Server) Part(channelName string) bool {
    for index, channel := range(s.channels) {
        if channel.Name == channelName {
            s.deleteChannel(index)
            return true
        }
    }
    return false
}

// _, err := connection.Read(reply)
// if err != nil {
//     fmt.Println("Read from server failed:", err.Error())
//     return
// }
//
// channelFromServer <-strings.TrimSpace(string(reply))
// response := strings.TrimSpace(string(reply))
// log.Println(response)
// fmt.Println(response)
//
// select {
//     case msgToSend, ok := <-channelToServer:
//         if ok {
//             log.Println("Sending: " + msgToSend)
//             rawBytes := []byte(msgToSend)
//             connection.Write(rawBytes)
//         } else {
//             // time.Sleep(time.Second)
//         }
//     default:
//         continue
//     }
