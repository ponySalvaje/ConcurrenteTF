package main

import (
	"bufio"
	"fmt"
	"net"
	"encoding/json"
)


var nodos  map[string]bool

func Client_1(subdir, serverdir string ){
	host := fmt.Sprintf("%s:8001", serverdir)
	fmt.Println(host)
	con, _ := net.Dial("tcp",host)

	defer con.Close()
	fmt.Fprintln(con,subdir)

	r:= bufio.NewReader(con)
	msg, _ := r.ReadString('\n')

	fmt.Println(msg)

	rest := make(map[string]bool)

	_ =json.Unmarshal([]byte(msg), &rest )
	fmt.Println(rest)

	for nodo := range rest{

		nodos[nodo]= true
	}
	fmt.Println(nodos)
}

func main() {
	nodos = make(map[string]bool)
	nodos["7000"] = true
	Client_1("7000", "localhost"
}









