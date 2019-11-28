package main

import (
	"fmt"
	"net"
)

func main(){

	con, _ := net.Dial("tcp", "localhost:8000");
	defer con.Close()
	fmt.Println("Ingrese valor: ")
	var msg string
	fmt.Scanf("%s", &msg)
	fmt.Fprintf(con, msg);
    
}