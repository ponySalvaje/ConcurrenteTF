package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

var nodos map[string]bool

func main() {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", "localhost:8001")
	for {
		connCliente1, _ := ln.Accept()
		go handle(connCliente1)
	}
}

func handle(conn net.Conn) {

	nodos = make(map[string]bool)
	nodos["8002"] = true
	nodos["8003"] = true
	nodos["8004"] = true

	//Recibir mensaje con IP del cliente1
	r := bufio.NewReader(conn)
	message, _ := r.ReadString('\n')
	fmt.Println("Mensaje recibido:", string(message))

	//Procesar -> agregar mensaje con dirección IP al mapa
	nodos[message] = true

	//Enviar el mapa al cliente1
	buf, _ := json.Marshal(nodos)
	fmt.Fprintf(conn, string(buf))

	//Enviar mensaje al cliente2
	fmt.Fprintf(conn, message)

}
