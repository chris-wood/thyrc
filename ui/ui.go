package ui

import (
    "fmt"
    "log"
    "github.com/jroimartin/gocui"
    "github.com/chris-wood/thyrc/message"
)

type ThyrcUI interface {
    CreateChannel(channel string) chan *message.Message
    GetInputChannel() chan string
}

type ConsoleUI struct {
    gui *gocui.Gui
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(v, "Hello world!")
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func New() *ConsoleUI {
    g := gocui.NewGui()
    if err := g.Init(); err != nil {
        log.Panicln(err)
    }

    g.SetLayout(layout)
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        log.Panicln(err)
    }
    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Panicln(err)
    }

    return &ConsoleUI{gui: g}
}

func (ui *ConsoleUI) CreateChannel() chan *message.Message {
    channel := make(chan *message.Message)
    return channel
}

func (ui *ConsoleUI) GetInputChannel() chan string {
    channel := make(chan string)
    return channel
}
