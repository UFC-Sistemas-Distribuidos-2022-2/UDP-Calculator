package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var comandIndex int
	var comandString string
	var firstNumber float32
	var secondNumber float32

	var salvarNumber bool = false

	var keepInFor bool = true

	p := make([]byte, 2048)
	for keepInFor {
		conn, err := net.Dial("udp", "127.0.0.1:1234")
		if err != nil {
			fmt.Printf("Some error %v", err)
			return
		}

		fmt.Println("Bem vindo a calculadora UDP")

		comandIndex, comandString = operationMenu()
		if !salvarNumber {
			firstNumber, secondNumber = numberReceiver()
		} else {
			firstNumber = numberReceiverWithParams(secondNumber)
		}

		data := map[string]interface{}{
			"operationIndex":  comandIndex,
			"operationComand": comandString,
			"firstNumber":     firstNumber,
			"secondNumber":    secondNumber,
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		_, err = conn.Write(jsonData)
		if err != nil {
			fmt.Printf("Write erro: %s\n", err)
			return
		}
		//fmt.Fprintf(conn, jsonstring)

		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
		conn.Close()
	}
}

func operationMenu() (int, string) {
	var loop bool = true

	var comand int = -1
	var comandString string
	var comandEscolhido Operacoes

	fmt.Println("Selecione a operação que deseja realizar:")

	for loop {
		fmt.Println("Digite 1 ou + ou Soma para somar")
		fmt.Println("Digite 2 ou - ou subtracao para subtrair")
		fmt.Println("Digite 3 ou / ou divisao para dividir")
		fmt.Println("Digite 4 ou * ou multiplicacao para multiplicar")

		fmt.Scan(&comandString)
		comandReturn, err := strconv.Atoi(comandString)
		if err != nil {
			if comandString == "+" || strings.ToLower(comandString) == "soma" {
				comandEscolhido = Operacoes(1)
			} else if comandString == "-" || strings.ToLower(comandString) == "subtracao" {
				comandEscolhido = Operacoes(2)
			} else if comandString == "/" || strings.ToLower(comandString) == "divisao" {
				comandEscolhido = Operacoes(3)
			} else if comandString == "*" || strings.ToLower(comandString) == "multiplicacao" {
				comandEscolhido = Operacoes(4)
			} else {
				consoleClear()
				fmt.Println("Ocorreu um erro")
				continue
			}
			return comandEscolhido.EnumIndex(), comandEscolhido.String()
		} else {
			comand = comandReturn
		}
		//fmt.Scan(&comand)

		if !(comand < 1 || comand > 4) {
			comandEscolhido = Operacoes(comand)
			fmt.Println(comandEscolhido)
			loop = false
			fmt.Println("O comando escolhido foi:", comandEscolhido.EnumIndex(), ":", comandEscolhido.String())
		} else {
			consoleClear()
			fmt.Println("Ocorreu um erro")
		}
	}

	return comandEscolhido.EnumIndex(), comandEscolhido.String()
}

func numberReceiver() (float32, float32) {
	var firstNumber float32
	var secondNumber float32

	fmt.Println("Digite o Primeiro número a ser usado:")
	fmt.Scan(&firstNumber)
	fmt.Println("Digite o Segundo número a ser usado:")
	fmt.Scan(&secondNumber)
	return firstNumber, secondNumber
}

func numberReceiverWithParams(firstNumber float32) float32 {
	var secondNumber float32
	fmt.Println("Digite o Segundo número a ser usado:")
	fmt.Scan(&secondNumber)
	return secondNumber
}

func (o Operacoes) String() string {
	return [...]string{"Soma", "Subtração", "Divisão", "Multiplicação"}[o-1]
}

func (o Operacoes) EnumIndex() int {
	return int(o)
}

type Operacoes int

const (
	Soma Operacoes = iota + 1
	Subtração
	Divisão
	Multiplicação
)

func consoleClear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
