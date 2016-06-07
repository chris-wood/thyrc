package message

import (
    // "time"
    "fmt"
    "strings"
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

func parseMessageParameters(message string) ([]string, error) {
    trimmedMessage := strings.TrimSpace(strings.Replace(message, "  ", " ", -1))
    return strings.Split(trimmedMessage, " "), nil
}

func parseMessageCommand(message string) (string, error) {
    length := len(message)
    for i := 0; i < length; i++ {
        if message[i] == ' ' {
            return string(message[0:i]), nil
        }
    }
    return "", MessageError{"Could not parse a command"}
}

func parseMessagePrefix(message string) string {
    return ""
}

func parseMessage(message string) *Message {
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

    parameters, paramError := parseMessageParameters(message[len(msg.Command) + 1:])
    if paramError == nil {
        msg.Parameters = parameters
    } else {
        fmt.Println("Error parsing parameters")
    }

    return msg
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

func (m *Message) Encode() string {
    result := m.Command
    for _, param := range(m.Parameters) {
        result += " " + param
    }
    result += "\n"

    return result;
}
