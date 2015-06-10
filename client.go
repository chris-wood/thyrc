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

type Command struct {
    command string
    parameters []string
}

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
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

        // Shove the output to the main view
        view, err := gui.View("main")
        if err != nil {
            log.Fatal("Could not recover handle to the main view.")
        }

        stringReply := strings.TrimSpace(string(reply))
        if len(stringReply) > 0 {
            fmt.Fprint(view, stringReply) 
        }

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

			connected = true // don't connect again
		}
    }
}

func startSession(gui *gocui.Gui, serverAddress string) (error) {
    connection, err := net.Dial("tcp", serverAddress)
    if err != nil {
        return err;
    }

    go readFromServer(gui, connection)
    // go prompt(gui, connection)

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
    if v, err := g.SetView("main", 0, 0, maxX - 30, maxY - 3); err != nil {
        if err != gocui.ErrorUnkView {
            fmt.Println("Error: could not set `main` view")
            return err
        }

        v.Editable = false
        v.Wrap = false

        if err := g.SetCurrentView("main"); err != nil {
            fmt.Println("couldn't set view to main")
            return err
        }
    }

    

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

    err := startSession(gui, args[0])
    if err != nil {
        fmt.Println("Failed.")
        fmt.Println("Error: " + string(err.Error()))
    } else {
        err = gui.MainLoop()
        if err != nil && err != gocui.Quit {
            log.Panicln(err)
        }
    }
}
