package ui

import (
    "log"
    "github.com/jroimartin/gocui"
    "github.com/chris-wood/thyrc/message"
)

type ThyrcUI interface {
    CreateWindow(channel string) chan *message.Message
    GetInputChannel() chan string
}

type ConsoleUI struct {
    gui *gocui.Gui
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    if _, err := g.SetView("main", 0, 0, maxX-1, maxY-5); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
    }

    if v, err := g.SetView("input", 0, maxY-5, maxX-1, maxY-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        if err := g.SetCurrentView("input"); err != nil {
            return err
        }
        v.Editable = true
        v.Wrap = true
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

func (ui *ConsoleUI) CreateWindow(windowName string) chan *message.Message {
    // TODO: create a new view here
    channel := make(chan *message.Message)
    return channel
}

func (ui *ConsoleUI) GetInputChannel() chan string {
    channel := make(chan string)
    return channel
}
