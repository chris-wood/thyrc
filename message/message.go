package message

import (
    // "time"
)

type Message struct {
    Prefix string
    Command string
    Parameters []string
}

func parseMessageParameters(message []byte) []string {
    list := make([]string, 1)

    return list
}

func parseMessageCommand(message []byte) string {
    return ""
}

func parseMessagePrefix(message []byte) string {
    return ""
}

func parseMessageAsBytes(message []byte) *Message {
    msg := new(Message)

    if message[0] == ':' {
        msg.Prefix = parseMessagePrefix(message[1:])
    }

    msg.Command = parseMessageCommand(message)
    msg.Parameters = parseMessageParameters(message[len(msg.Command):])

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
