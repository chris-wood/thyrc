package main

import (
    "time"
    "log"
    "net"
    "os"
    "fmt"
    "strings"
    "github.com/jroimartin/gocui"
)

type Command struct {
    command string
    parameters []string
}

func handleCommand(command Command, connection net.Conn) (bool) {
    if command.command == "quit" {
        return false 
    } else {
        fmt.Println(command.command)
        return true
    }
}

func parseCommandString(commandString string) (Command) {
    var command Command
    stringFields := strings.Fields(commandString)
    if len(stringFields) > 0 {
        command.command = stringFields[0]
        if len(stringFields) > 1 {
            command.parameters = stringFields[1:]    
        } else{
            command.parameters = nil
        }
    } else {
        // error?
    }
    
    return command
}

func stopSession(connection net.Conn) {
    err := connection.Close()
    if err != nil {
        log.Fatal(err)
    }
}

func serverReadAndWrite(channelFromServer chan string, channelToServer chan string, connection net.Conn) {
    reply := make([]byte, 1024)
    stayAlive := true

    for ; stayAlive ; {
        stayAlive = true // TODO: read from concurrent channel here
        _, err := connection.Read(reply)
        if err != nil {
            fmt.Println("Write to server failed:", err.Error())
            return
        }

        // channelFromServer <- strings.TrimSpace(string(reply))
        response := strings.TrimSpace(string(reply))
        log.Println(response)
        fmt.Println(response)
        select {
            case msgToSend, ok := <-channelToServer:
                if ok {
                    log.Println("Sending: " + msgToSend)
                    rawBytes := []byte(msgToSend)
                    connection.Write(rawBytes)
                } else {
                    time.Sleep(time.Second)
                }
            default:
                continue
            }

        // // Shove the output to the main view
        // view, err := gui.View("main")
        // if err != nil {
        //     log.Fatal("Could not recover handle to the main view.")
        // }
    }
}

func ircHandler(channelFromServer chan string, channelToServer chan string) {
    alive := true
    msgToSend := ""
    connected := false

    for ; alive ; {

        if !connected {
            channelToServer <- "PASS none\n"
            channelToServer <- "NICK random\n"
            channelToServer <- "USER rawrrawr blah blah blah\n"
            connected = true // don't connect again
        } else {
            // Read and display the server response
            response := <-channelFromServer
            fmt.Println(response)

            // Read from stdin
            fmt.Scanln(msgToSend)
            channelToServer <- msgToSend

            // Reset the message to send
            msgToSend = ""
        }
    }

    // stringReply := strings.TrimSpace(string(reply))
    // if len(stringReply) > 0 {
    //     // fmt.Println(stringReply)
    //     // view.Clear()
    //     fmt.Fprintln(view, stringReply) 
    //     gui.Flush()
    // }
}

func runSession(gui *gocui.Gui, serverAddress string) (error) {
    connection, err := net.Dial("tcp", serverAddress)
    if err != nil {
        return err;
    }

    channelFromServer := make(chan string)
    channelToServer := make(chan string)
    killChannel := make(chan int)

    go serverReadAndWrite(channelFromServer, channelToServer, connection)
    go ircHandler(channelFromServer, channelToServer)

    <-killChannel // block until killed
    stopSession(connection)

    return nil
}

func nextView(g *gocui.Gui, v *gocui.View) (error) {
    if v == nil || v.Name() == "rightside" {
        return g.SetCurrentView("main")
    } else if v == nil || v.Name() == "main" {
        return g.SetCurrentView("input")
    } else {
        return g.SetCurrentView("rightside")
    }
}

func cursorDown(g *gocui.Gui, v *gocui.View) (error) {
    if v != nil {
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy+1); err != nil {
            ox, oy := v.Origin()
            if err := v.SetOrigin(ox, oy+1); err != nil {
                return err
            }
        }
    }
    return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) (error) {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-1); err != nil {
                return err
            }
        }
    }
    return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) (error) {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx-1, cy); err != nil && ox > 0 {
            if err := v.SetOrigin(ox-1, oy); err != nil {
                return err
            }
        }
    }
    return nil
}

func cursorRight(g *gocui.Gui, v *gocui.View) (error) {
    if v != nil {
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx+1, cy); err != nil {
            ox, oy := v.Origin()
            if err := v.SetOrigin(ox+1, oy); err != nil {
                return err
            }
        }
    }
    return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
    var l string
    var err error

    _, cy := v.Cursor()
    if l, err = v.Line(cy); err != nil {
        l = ""
    }

    maxX, maxY := g.Size()
    if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
        if err != gocui.ErrorUnkView {
            return err
        }
        fmt.Fprintln(v, l)
        if err := g.SetCurrentView("msg"); err != nil {
            return err
        }
    }
    return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
    if err := g.DeleteView("msg"); err != nil {
        return err
    }
    if err := g.SetCurrentView("side"); err != nil {
        return err
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.Quit
}

func keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("rightside", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, cursorLeft); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, cursorRight); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        return err
    }
    // if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
    //     return err
    // }
    if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
        return err
    }

    return nil
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    // right side (channel list) 
    if v, err := g.SetView("rightside", maxX - 30, 0, maxX - 1, maxY - 3); err != nil {
        if err != gocui.ErrorUnkView {
            return err
        }
        v.Highlight = true
        fmt.Fprintln(v, "Channel list")
    }

    // input field 
    if v, err := g.SetView("input", 0, maxY - 3, maxX - 1, maxY - 1); err != nil {
        if err != gocui.ErrorUnkView {
            fmt.Println("Error: could not set `main` view")
            return err
        }

        v.Editable = true
        v.Wrap = false
        v.Highlight = true
    }

    // main side 
    if v, err := g.SetView("main", 10, 0, maxX - 20, maxY - 3); err != nil {
        if err != gocui.ErrorUnkView {
            fmt.Println("Error: could not set `main` view")
            return err
        }

        v.Editable = false
        v.Wrap = true
        fmt.Fprintln(v, "Here we go...")

        if err := g.SetCurrentView("main"); err != nil {
            fmt.Println("couldn't set view to main")
            return err
        }
    }

    return nil
}

func main() {
    file, err := os.OpenFile("client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer file.Close()
    log.SetOutput(file)

    fmt.Fprintln(os.Stderr, "Setup the log")

    args := os.Args[1:]
    if len(args) < 1 {
        fmt.Println("usage: go run client.go <server:port>")
        log.Fatalf("usage")
    }

    fmt.Fprintln(os.Stderr, "Got the arguments, setting up the GUI")

    // gui := gocui.NewGui()
    // fmt.Fprintln(os.Stderr, "Initializing.")
    // err = gui.Init();
    // fmt.Fprintln(os.Stderr, "Checking for failure.")
    // if err != nil {
    //     fmt.Fprintln(os.Stderr, "Error during initialization: %s", err)
    //     log.Panicln(err)
    // }
    // defer gui.Close()

    // fmt.Fprintln(os.Stderr, "Setup the GUI")

    // gui.SetLayout(layout)
    // if err := keybindings(gui); err != nil {
    //     log.Panicln(err)
    // }
    // gui.SelBgColor = gocui.ColorGreen
    // gui.SelFgColor = gocui.ColorBlack
    // gui.ShowCursor = true

    err = runSession(nil, args[0])
    // if err != nil {
    //     fmt.Fprintln(os.Stderr, "Failed.")
    //     fmt.Fprintln(os.Stderr, "Error: " + string(err.Error()))
    // } else {
    //     err = gui.MainLoop()
    //     if err != nil && err != gocui.Quit {
    //         log.Panicln(err)
    //     }
    // }
}
