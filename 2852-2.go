package main

import (
	"fmt"
	"math/rand"
	"time"
	"net"
	"log"
	"encoding/json"
)

type pokemon struct {
	nombre string
	latitud int
	longitud int
}

type jugador struct {
	nombre string
	latitud int
	longitud int
}

type tMsg struct {
	Code int
	Msg string
}

var localAddr string
var otherAddrs []string
var pokemons []pokemon

func main(){

	fmt.Print("Dirección local: ")
	fmt.Scanln(&localAddr)
	
	otherPlayerAddr := ""

	fmt.Print("Dirección de otro jugador: ")
	fmt.Scanln(&otherPlayerAddr)

	if otherPlayerAddr == "" {	//Primera conexion
		
	} else {
		otherAddrs = append(otherAddrs, otherPlayerAddr)
		//Update en el nuevo
		sendMsg(otherPlayerAddr, tMsg {0, localAddr})
	}
	go poKreate()
	server()
}

func randomPokeSpawn() pokemon{
	rand.Seed(time.Now().UTC().UnixNano())
	randX := rand.Intn(100)
	randY := rand.Intn(100)
	newPokemon := pokemon {"Chu", randX, randY}
	return newPokemon
}

func server() {
	if ln, err := net.Listen("tcp", localAddr); err != nil {
		log.Panicln(err.Error())
	} else {
		defer ln.Close()
		fmt.Println(localAddr, "listening")
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Panicln(err.Error())
			} else {
				go handle(conn)
			}
		}
	}
}

func sendMsg(addr string, msg tMsg) {
	if conn, err := net.Dial("tcp", addr); err != nil {
		log.Println(err.Error())
	} else {
		defer conn.Close()
		enc := json.NewEncoder(conn)
		fmt.Println("Sending", msg, "to", addr)
		enc.Encode(&msg)
	}
}

func addToNet(newAddr string) {
	for _, addr := range otherAddrs {
		sendMsg(addr, tMsg {1, newAddr})
	}
}

func poKreate(){
	for {
		var tPassed int64;
		initTime := time.Now().UTC().UnixNano();
		for {
			tPassed = - initTime + time.Now().UTC().UnixNano();
			if tPassed > 60000000000 { //Cada 1 minuto
				fmt.Println("Creating pokemon!")
				newPokemon := poKreate(); //Creacion de pokemon
				//Buscar los 2 mas cercanos
				for i, addrs := range otherAddrs{
					//Mandar solicitud de confirmacion
					sendMsg(addrs, tMsg {2, newPokemon}){}
					if i > 2 {
						break;
					}
				}
				break;
			}
		}
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	//fmt.Println(conn.RemoteAddr(), "accepted")
	var msg tMsg
	dec := json.NewDecoder(conn)
	if err := dec.Decode(&msg); err != nil {
		log.Println(err.Error())
	} else {
		fmt.Println("Got", msg)
		switch msg.Code {
		case 0: //Agregar en existente

			accepted := ""

			fmt.Print("El jugador ", msg.Msg, " se quiere conectar mediante ti. Aceptar? (y/n): ")
			fmt.Scanln(&accepted)

			if accepted == "y" {
				addToNet(msg.Msg)
				fmt.Println(msg.Msg, " se ha conectado mediante ti!")
				otherAddrs = append(otherAddrs, msg.Msg)
				//fmt.Println(otherAddrs)
			} 

		case 1: //Agregado en red
			otherAddrs = append(otherAddrs, msg.Msg)

			fmt.Print("What´s going on folks? Ha llegado un nuevo entrenador! Los entrenadores en la red ahora son los siguientes: ")
			fmt.Println(otherAddrs)
		case 2: //Agregar nuevo pokemon
			accepted := ""

			fmt.Print("Se desea agregar el pokemon ", msg.Msg, ". Aceptar? (y/n): ")
			fmt.Scanln(&accepted)

			if accepted == "y" {
				addToNet(msg.Msg)
				fmt.Println(msg.Msg, " se creo el Pokemon")
				pokemon = append(pokemon, msg.Msg)
			} 
		}
	}
}