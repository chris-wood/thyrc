package main

import (
    "io"
    "log"
    "net"
    "os"
)

func main() {
    log.SetFlags(log.Lshortfile)
    con, err := net.Dial("tcp", ":7")
    if err != nil {
        log.Fatal(err)
    }

    _, err = io.Copy(con, os.Stdin)
    if err != nil {
        log.Fatal(err)
    }
    if tcpcon, ok := con.(*net.TCPConn); ok {
        tcpcon.CloseWrite()
    }

    _, err = io.Copy(os.Stdout, con)
    if err != nil {
        log.Fatal(err)
    }

    err = con.Close()
    if err != nil {
        log.Fatal(err)
    }
}