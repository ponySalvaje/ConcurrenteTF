package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	PROT = "tcp"
	HOST = "localhost:8005"
)

var nodes map[string]bool

func handle(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
	msg, _ := r.ReadString('\n')
	fmt.Print("Recibido: ", msg)
	nodes[msg] = true
}

func main() {
	nodes = make(map[string]bool)
	ln, _ := net.Listen(PROT, HOST)
	defer ln.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		con, _ := ln.Accept()
		go handle(con)
	}

}
