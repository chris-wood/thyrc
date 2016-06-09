package ui

import (
    "os"
    "bufio"
    "fmt"
    "github.com/chris-wood/thyrc/message"
)

type ThyrcUI interface {
    CreateWindow(channel string) chan *message.Message
    GetInputChannel() chan string
}

type TextUI struct {
    channels []chan *message.Message
    inputChannel chan string
}

func (ui TextUI) handleUserInput() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("> ")
        text, err := reader.ReadString('\n')
        if err == nil {
            ui.inputChannel <- text
        }
    }
}

func NewTextUI() ThyrcUI {
    textui := TextUI{inputChannel: make(chan string)}
    go textui.handleUserInput()
    return TextUI{}
}

func (ui TextUI) serveWindow(inputChannel chan *message.Message) {
    for {
        msg := <-inputChannel
        fmt.Println(msg)
    }
}

func (ui TextUI) CreateWindow(windowName string) chan *message.Message {
    channel := make(chan *message.Message)
    ui.channels = append(ui.channels, channel)
    return channel
}

func (ui TextUI) GetInputChannel() chan string {
    return ui.inputChannel
}
