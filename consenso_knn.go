package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"os"
	"encoding/csv"
	"io"
	"math"
	"strconv"
	"sort"
)

var dimensions int
var inputs []string

func main(){

	var nodes int
	var isServer string
	k := "";

	fmt.Print("Este programa sera el servidor (y/n): ")
	fmt.Scanf("%s\n", &isServer)

	if strings.ToLower(isServer) == "y" { //Server code

		fmt.Print("Ingrese la cantidad de nodos: ")
		fmt.Scanf("%d\n", &nodes)

		fmt.Print("Definir el valor de K: ")
		fmt.Scanf("%s\n", &k)

		ips := make([]string, nodes + 1) // Index 0 is server ip
		ipCounter := 0
		fmt.Print("La ip del servidor: ")
		ips[0] = getOutboundIP()
		fmt.Println(ips[0])
		nodesToEnd := nodes

		datasets := loadCsvAndPartition("wpbc.csv", nodes, true)

		newInstance := make([]float64, dimensions)

		for i := 0; i < dimensions; i++ {
			fmt.Print("Ingrese el valor " + string(i+49) + " de la nueva instancia: ")
			fmt.Scanf("%g\n", &newInstance[i])
		}

		host := fmt.Sprintf("%s:8000", ips[0])
		ln, _ := net.Listen("tcp", host)
		defer ln.Close()
		
		go getResult(&nodesToEnd)
		
		for {
			con, _ := ln.Accept()
			ipCounter++
			ips[ipCounter] = con.LocalAddr().String()
			fmt.Println(con.LocalAddr().String() + " se ha conectado.")
			go handle(k, newInstance, datasets[ipCounter-1], nodes, &nodesToEnd, con)
		}

	} else if strings.ToLower(isServer) == "n" { //Client code

		con,_ := net.Dial("tcp", "192.168.0.22:8000")
		r := bufio.NewReader(con)
		msgInst, _ := r.ReadString('\n')
		newInst := Vec{}
		tempStr := ""
		fmt.Println(msgInst)

		for _, ltr := range msgInst {
			if ltr == ',' {
				e, _ := strconv.ParseFloat(tempStr, 64)
				newInst.elements = append(newInst.elements, e )
				tempStr = ""
			} else{
				tempStr += string(ltr)
			}
		}

		newInst.size = len(newInst.elements)
		newInst.class = ""

		msgData, _ := r.ReadString('\n')
		vecs := format(msgData)
		
		kData, _ := r.ReadString('\n')
		
		kInt, _ := strconv.Atoi(strings.TrimSpace(kData))

		xd := kInt
		
		fmt.Fprintf(con, knn(xd, newInst, vecs))

		} else { //Did an oopsie	
			fmt.Print("Datos de servidor ingresado incorrrectamente, intentelo otra vez :(")
	}
}

func handle(k string, newInst []float64, ds string, n int, nToEnd* int, con net.Conn){
	defer con.Close()
	//Send dataset

	tempStr := ""
	for i := 0; i < len(newInst) ; i++ {
		tempStr += fmt.Sprintf("%g", newInst[i])
		tempStr += ","
	}

	tempStr += "\n"
	
	fmt.Fprintf(con, tempStr) //Send new inst
	 
	fmt.Fprintf(con, ds) //Send part of dataset

	fmt.Fprintf(con, k + "\n") //Send new inst

	r := bufio.NewReader(con)
	msg, _ := r.ReadString('\n') //Read result
	inputs = append(inputs, msg)
	*nToEnd--
	fmt.Println(msg) //Store result
}

func getOutboundIP() string {
    con, _ := net.Dial("udp", "8.8.8.8:80")
    defer con.Close()
	localAddr := con.LocalAddr().String()
    return strings.Split(localAddr, ":")[0]
}

func loadCsvAndPartition(filePath string, n int, hasHeader bool)  []string{

	ds := []string{}
	it := 0
    // Loading file.
    f, _ := os.Open(filePath)
    // Reader.
	r := csv.NewReader(f)
	record, _ := r.Read()
	dimensions = len(record) - 1
    for {
		record, err := r.Read()
		// Stop at EOF.		
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
		}		
		ds = append(ds,"");
        for _, val := range record {
			ds[it] += val;
			ds[it] += ",";
		}
		ds[it] += ";";
		it++
	}

	div := it/n+1
	res := it%n
	doneOnce := false
	j := 0

	dsF := make([]string, n)

	for i := 0; i < n; i++ {
		for _, e := range ds[j:j+div] {
			dsF[i] += e
		}
		res--
		j+=div
		fmt.Println(j)
		if res <= 0 {
			res = 1
			if doneOnce == false {
				doneOnce = true
				div--
			}
		}
		dsF[i] += "\n"
	}

	return dsF;
}

//KNN

type Vec struct{
	size int;
	elements []float64;
	class string;
}

func (v* Vec) distToVec(toV Vec) float64 { //Calculate distance to another Vec

	sums := []float64{};
	var sum float64 = 0;
	
	for i := 0; i < v.size; i++ {
		sums = append(sums, v.elements[i] - toV.elements[i]);
	}

	for i := 0; i < len(sums); i++ {
		sum += math.Pow(sums[i], 2);
	}
	
	return math.Sqrt(sum);
}


func getResult(nToEnd* int){
	
	for *nToEnd > 0 {
		
	}

	counterMap := make(map[string]int)

	for i:= 0; i < len(inputs); i++ {
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
	
	fmt.Println("La clase es: " + result + "!!!!!!!!!!")
}

func format(ds string) []Vec{
	v := []Vec{}
	newVec := Vec{}
	str := ""
	for i, ltr := range ds {
		if ltr == ','{
			if ds[i+1] == ';' {
				newVec.class = str
				newVec.size = len(newVec.class)
				v = append(v, newVec)
				newVec = Vec{}
				str = ""
			} else {
				e, _ := strconv.ParseFloat(str, 64)
				newVec.elements = append(newVec.elements, e)
				str = ""
			}
		}else {
			str += string(ltr);
		}
	}
	return v
}

func knn(k int, predV Vec, v []Vec) string {

	var class string; //Predicted class
	
	dists := []float64{};
	var dist float64;
	m := make(map[float64]Vec) //Dist to Vec map;

	for i := 0; i < len(v); i++ {
		dist = predV.distToVec(v[i]);
		dists = append(dists, dist);
		m[dist] = v[i];
	}

	sort.Float64s(dists);

	closest := make([]Vec, k);
	count := make(map[string]int);

	for i := 0; i < k; i++ {
		closest[i] = m[dists[i]];
		count[closest[i].class]++;
	}

	max := 0;

	for i, val := range count{
		if(max < val){
			max = val;
			class = i;
		}		
	}

	return class;
}