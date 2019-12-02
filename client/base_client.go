package main

import (
	"fmt"
	"net"
	"bufio"
	"strconv"
	"strings"
	"encoding/json"
)

type Dato struct {
	Id int  `json:"id"`
	Columna int  `json:"columna"`
	Valor string  `json:"valor"`
}

type Pred struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	L float64 `json:"L"`
}

func main(){
	con,_ := net.Dial("tcp", "192.168.0.22:8000")
	resp := ""
	r := bufio.NewReader(con)

	for{
		resp = ""
		fmt.Print("Desea agregar una nueva Muestra o predecir? (1: Agregar - 2: Predecir - 3: Visualziar datos): ")
		fmt.Scanf("%s\n", &resp)
		
		switch resp {
		case "1":
			fmt.Fprintf(con, "AGREGAR\n");
			
			msgId := ""

			for {
				msg, err := r.ReadString('\n')
				if err != nil {
					continue;
				}
				msgId = strings.TrimSpace(msg);
				break;
			}



			for {
				dato := Dato{};
				dato.Id, _ = strconv.Atoi(msgId)

				fmt.Print("Que dato desea ingresar (1: Dias transcurridos - 2: Color a* - 3: Color b* - 4: Color L* - 5: Salir): ")
				fmt.Scanf("%s\n", &resp)

				col, _ := strconv.Atoi(resp);
				dato.Columna = col;
				brake := false;

				var val string;
				switch col {
				case 1:
					fmt.Print("Ingrese los dias transcurridos: ");
					fmt.Scanf("%s\n", &val);
				case 2:
					fmt.Print("Ingrese el color a*: ");
					fmt.Scanf("%s\n", &val);
				case 3:
					fmt.Print("Ingrese el color b*: ");
					fmt.Scanf("%s\n", &val);
				case 4:
					fmt.Print("Ingrese el color c*: ");
					fmt.Scanf("%s\n", &val);
				case 5:
					fmt.Fprintf(con, "BREAK\n");
					brake = true;
				}

				if brake { 
					break ;
				}

				dato.Valor = val;

				jsonBytes, _ := json.Marshal(dato);
				jsonStr := string(jsonBytes) + "\n";

				//fmt.Println(jsonStr);
				fmt.Fprintf(con, jsonStr);
				
			}			
		case "2":
			fmt.Fprintf(con, "PREDECIR\n");
			predData := Pred{};
			for predData.A == 0 || predData.B == 0 || predData.L == 0{

				fmt.Print("Que dato desea ingresar (1: Color a* - 2: Color b* - 3: Color L*): ")
				fmt.Scanf("%s\n", &resp)

				col, _ := strconv.Atoi(resp);

				var val float64;
				switch col {	
				case 1:
					fmt.Print("Ingrese el color a*: ");
					fmt.Scanf("%f\n", &val);
					predData.A = val;
				case 2:
					fmt.Print("Ingrese el color b*: ");
					fmt.Scanf("%f\n", &val);
					predData.B = val;
				case 3:
					fmt.Print("Ingrese el color L*: ");
					fmt.Scanf("%f\n", &val);
					predData.L = val;
				}			
			}

			jsonBytes, _ := json.Marshal(predData);
			jsonStr := string(jsonBytes) + "\n";

			fmt.Fprintf(con, jsonStr);

			respChan := make(chan string);
			
			go func(){
				for{
					msgResp, err := r.ReadString('\n');
					if err != nil {
						continue;
					}
					msgResp = strings.TrimSpace(msgResp);
					respChan<- msgResp;
					break;
				}
			}()
				resp := <-respChan
			
				fmt.Print("La prediccion resulta ser: ");
				fmt.Println(resp);
			case "3":
				fmt.Fprintf(con, "VISUALIZAR\n");
				msgTable := ""

				for {
					msg, err := r.ReadString('z')
					if err != nil {
						continue;
					}
					msgTable = msg;
					break;
				}

				fmt.Println(msgTable);

		}
	}
}

func update(){}

func retry(){
	//El dato ingresado es erroneo, intentelo otra vez

}