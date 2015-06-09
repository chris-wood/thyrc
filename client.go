package main

import (
    "bufio"
    "log"
    "net"
    "os"
    "fmt"
    "strings"
    "io"
    "io/ioutil"
    "github.com/jroimartin/gocui"
)

// go run client.go irc.freenode.net:6666

// startup protocol:
//PASS none
//NICK sorandom29      
//USER blah blah blah blah

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

func prompt(gui *gocui.Gui, connection net.Conn) (bool, error) {
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')

        command := parseCommandString(text)
        if !handleCommand(command, connection) {
            break
        }
    }
    return true, nil
}

func readFromServer(gui *gocui.Gui, connection net.Conn) {
    reply := make([]byte, 1024)
    stayAlive := true 

	connected := false

    for ; stayAlive ; {
        stayAlive = true // TODO: read from concurrent channel here
        _, err := connection.Read(reply)
        if err != nil {
            fmt.Println("Write to server failed:", err.Error())
            return
        }
        fmt.Println("$> " + string(reply))

        // Shove the output to the main view
        // view, err := gui.View("main")
        // fmt.Fprintf(view, "%s", string(reply))

		if !connected {
			// pass, nick, user
			// fmt.Println("Sending PASS...") 
			passCommand := []byte("PASS none\n")
			connection.Write(passCommand)

			// fmt.Println("Sending NICK...")
			nickCommand := []byte("NICK random\n")
			connection.Write(nickCommand)

			// fmt.Println("Sending USER...")
			userCommand := []byte("USER rawrrawr blah blah blah\n")
			connection.Write(userCommand)

			connected = true
		}
    }
}

func startSession(gui *gocui.Gui, serverAddress string) error {
    // fmt.Println("Connecting to " + serverAddress + "...")
    connection, err := net.Dial("tcp", serverAddress)
    if err != nil {
        return err;
    }

    go readFromServer(gui, connection)
    prompt(gui, connection)

    // fmt.Println("IRC setup, returning to GUI connection")

    return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
    if v == nil || v.Name() == "side" {
        return g.SetCurrentView("main")
    }
    return g.SetCurrentView("side")
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
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

func cursorUp(g *gocui.Gui, v *gocui.View) error {
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

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
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

func cursorRight(g *gocui.Gui, v *gocui.View) error {
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
    if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
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
    if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
        return err
    }
    if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
        return err
    }

    if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
        return err
    }
    return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
    f, err := ioutil.TempFile("", "gocui_demo_")
    if err != nil {
        return err
    }
    defer f.Close()

    p := make([]byte, 5)
    v.Rewind()
    for {
        n, err := v.Read(p)
        if n > 0 {
            if _, err := f.Write(p[:n]); err != nil {
                return err
            }
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }
    return nil
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    // left side
    if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
        if err != gocui.ErrorUnkView {
            return err
        }
        v.Highlight = true
        fmt.Fprintln(v, "Item 1")
        fmt.Fprintln(v, "Item 2")
        fmt.Fprintln(v, "Item 3")
        fmt.Fprint(v, "\rWill be")
        fmt.Fprint(v, "deleted\rItem 4\nItem 5")
    }

    // main side
    if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
        if err != gocui.ErrorUnkView {
            return err
        }
        // b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")
        // if err != nil {
        //     panic(err)
        // }
        // fmt.Fprintf(v, "%s", b)
        v.Editable = false
        v.Wrap = true
        if err := g.SetCurrentView("main"); err != nil {
            return err
        }
    }

    // TODO: input text field

    return nil
}

func main() {
    args := os.Args[1:]
    if len(args) < 1 {
        fmt.Println("usage: go run client.go <server:port>")
        return 
    }
    log.SetFlags(log.Lshortfile)

    gui := gocui.NewGui()
    if err := gui.Init(); err != nil {
        log.Panicln(err)
    }
    defer gui.Close()

    gui.SetLayout(layout)
    if err := keybindings(gui); err != nil {
        log.Panicln(err)
    }
    gui.SelBgColor = gocui.ColorGreen
    gui.SelFgColor = gocui.ColorBlack
    gui.ShowCursor = true

    err := startSession(nil, args[0])
    if err != nil {
        fmt.Println("Failed.")
        fmt.Println("Error: " + string(err.Error()))
    } else {
        err := gui.MainLoop()
        if err != nil && err != gocui.Quit {
            log.Panicln(err)
        }
    }
}
