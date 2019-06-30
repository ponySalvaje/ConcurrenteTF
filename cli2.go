package main

import (
	"fmt"
	"net"
)

func main() {
	nodos:= make( map [string ] bool)
	nodos["8005"]=true 
	cli2("hola mundo",nodos)
}

 func cli2( dir string, nodos map [string] bool){

	
	for	nodo := range nodos { 
		host:=fmt.Sprintf("localhost:%s",nodo)
		fmt.Println(host)
		con, _ := net.Dial("tcp", host)
		fmt.Fprintln(con,dir)
		con.Close()
	}
 }