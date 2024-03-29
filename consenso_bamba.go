package main

import (
	"fmt"
	"net"
	"bufio"
)


var nodes int = 0
var inputs []string;

func main(){

	
	host := fmt.Sprintf("%s:8000", "localhost")
	ln, _ := net.Listen("tcp", host)

	go consensus()

	defer ln.Close()
	for {
		con, _ := ln.Accept()
		fmt.Println("Conectado!")
		nodes++
		go handle(con)
	}
}

func handle(con net.Conn){
	defer con.Close()
	r := bufio.NewReader(con)
	msg, _ := r.ReadString('\n')
	fmt.Println(msg)
	inputs = append(inputs, msg)
}

func consensus(){
	
	for len(inputs) == 0 || len(inputs) != nodes {
		
	}

	counterMap := make(map[string]int)

	for i:= 0; i < nodes; i++ {
		counterMap[inputs[i]]++
	}

	max := 0
	result := ""

	for i, val := range counterMap {
		if val > max {
			max = val
			result = i
		}
	}
	
	fmt.Println("La mayoria es: " + result + "!!!")
}