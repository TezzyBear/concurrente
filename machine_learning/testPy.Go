package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"os"
	"os/exec"
)

func main() {
	/*
	*	toPred = float64 (numero de entrada)
	*	para este script solo se esta recogiendo de input
	*	usar el codigo comentado justo debajo de esta parte en el programa blockchain
	*/
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresa numero E*: ")
	toPred, _ := gin.ReadString('\n')
	toPred = strings.TrimSpace(toPred)

	/*num := 1.94 //float64
	toPred := fmt.Sprintf("%f", num)*/


	//enviar comando de ejecucion python
	cmd := exec.Command("python", "LinearRegression.py", "makePrediction", toPred)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	//limpiar output
	strOut := string(out)
	strOut = strings.TrimSpace(strOut)
	strOut = strings.Replace(strOut, "[", "", -1)
	strOut = strings.Replace(strOut, "]", "", -1)
	split := strings.Split(strOut, ",")

	//prints para probar (quitar en el programa blockchain)
	fmt.Println(strOut)
	fmt.Println(split)

	//conversion de resultados (split[0] = regresion lineal, split[1] = regresion polinomial)
	linearResult, err := strconv.ParseFloat(split[0], 64)
	//print para probar
	if err == nil {
		fmt.Println(linearResult)
	}

	polynomialResult, err := strconv.ParseFloat(split[1], 64)
	//print para probar
	if err == nil {
		fmt.Println(polynomialResult)
	}
	//solo usar los valores de linearResult o polynomialResult
	//ambos tienen ventajas y desventajas
}