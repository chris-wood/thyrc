package message

import (
    // "time"
    "fmt"
)

type Message struct {
    Prefix string
    Command string
    Parameters []string
}

type MessageError struct {
    message string
}

func (e MessageError) Error() string {
    return e.message
}

func parseMessageParameters(message []byte) ([]string, error) {
    list := make([]string, 1)

    return list, nil
}

func parseMessageCommand(message []byte) (string, error) {
    length := len(message)
    for i := 0; i < length; i++ {
        if message[i] == ' ' {
            return string(message[0:i]), nil
        }
    }
    return "", MessageError{"Could not parse a command"}
}

func parseMessagePrefix(message []byte) string {
    return ""
}

func parseMessageAsBytes(message []byte) *Message {
    msg := new(Message)

    if message[0] == ':' {
        msg.Prefix = parseMessagePrefix(message[1:])
    }

    command, commandError := parseMessageCommand(message)
    if commandError == nil {
        msg.Command = command
    } else {
        fmt.Println("Error parsing command")
    }

    parameters, paramError := parseMessageParameters(message[len(msg.Command):])
    if paramError == nil {
        msg.Parameters = parameters
    } else {
        fmt.Println("Error parsing parameters")
    }

    return msg
}

func parseMessage(message string) *Message {
    return parseMessageAsBytes([]byte(message))
}

// Parse a message string and create a new instance of the Message object.
func Parse(messageString string) *Message {
	return parseMessage(messageString)
}

// New creates a new instance of the Message object.
func New(prefix string, command string, parameters []string) *Message {
	return &Message{
        Prefix: prefix,
        Command: command,
        Parameters: parameters,
    }
}
