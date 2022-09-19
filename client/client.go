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
	var firstNumber float64
	var secondNumber float64
	var valueWasSaved bool = false
	var keepInFor bool = true
	var keepInLoopAfterReceiveResponse bool = true
	var saveValue string
	var msgReceiver MsgReceiver

	fmt.Println("Bem vindo a calculadora UDP")

	p := make([]byte, 2048)
	for keepInFor {
		conn, err := net.Dial("udp", "127.0.0.1:1234")
		if err != nil {
			fmt.Printf("Some error %v", err)
			return
		}
		comandIndex, comandString = operationMenu()
		if comandIndex == -1 && comandString == "exit" {
			conn.Close()
			return
		}
		if !valueWasSaved {
			firstNumber, secondNumber = numberReceiver()
		} else {
			secondNumber = numberReceiverWithParams(firstNumber)
		}

		jsonData, err := jsonBuilder(comandIndex,
			comandString, firstNumber, secondNumber)

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

		n, err := bufio.NewReader(conn).Read(p)
		if err == nil {

			err = json.Unmarshal(p[:n], &msgReceiver)
			if err != nil {
				fmt.Printf("could not Unmarshal json: %s\n", err)
				return
			}

			fmt.Printf("%s\n", msgReceiver.Msg)

			for keepInLoopAfterReceiveResponse {

				if msgReceiver.Result != "error" {
					fmt.Println("Wanna save the value to use on the next operation? Y/N")
					fmt.Scan(&saveValue)
				}

				if strings.ToLower(saveValue) == "y" && msgReceiver.Result != "error" {
					valueWasSaved = true
					firstNumber, err = strconv.ParseFloat(msgReceiver.Result, 32)
					if err != nil {
						fmt.Printf("Some error %v\n", err)
					}
					secondNumber = 0
					keepInLoopAfterReceiveResponse = false

				} else if strings.ToLower(saveValue) == "n" || msgReceiver.Result == "error" {

					valueWasSaved = false
					firstNumber = 0
					secondNumber = 0
					keepInLoopAfterReceiveResponse = false

				} else {
					consoleClear()
					fmt.Println("An error occurred while reading the chosen operation")
				}
			}
			keepInLoopAfterReceiveResponse = true

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

	fmt.Println("Select the operation you want to perform:")

	for loop {

		fmt.Println("Type 1 or + or Soma to add")
		fmt.Println("Type 2 or - or subtracao to subtract")
		fmt.Println("Type 3 or / or divisao to divide")
		fmt.Println("Type 4 or * or multiplicacao to multiply")
		fmt.Println("Type 5 or sair to exit")

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
			} else if strings.ToLower(comandString) == "sair" {
				return -1, "exit"
			} else {
				consoleClear()
				fmt.Println("An error occurred while reading the chosen operation")
				continue
			}
			return comandEscolhido.EnumIndex(), comandEscolhido.String()
		} else {
			comand = comandReturn
		}
		if comand == 5 {
			return -1, "sair"
		}
		if !(comand < 1 || comand > 4) {
			comandEscolhido = Operacoes(comand)
			loop = false
			fmt.Println("The command chosen was:", comandEscolhido.EnumIndex(), ":", comandEscolhido.String())
		} else {
			consoleClear()
			fmt.Println("An error occurred while reading the chosen operation")
		}
	}

	return comandEscolhido.EnumIndex(), comandEscolhido.String()
}

func numberReceiver() (float64, float64) {
	var firstNumber float64
	var secondNumber float64

	fmt.Println("Enter the First number to use:")
	fmt.Scan(&firstNumber)
	fmt.Println("Enter the second number to use:")
	fmt.Scan(&secondNumber)
	return firstNumber, secondNumber
}

func numberReceiverWithParams(firstNumber float64) float64 {
	var secondNumber float64
	fmt.Println("Enter the second number to use:")
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

type MsgReceiver struct {
	Msg    string
	Result string
}

func jsonBuilder(comandIndex int, comandString string, firstNumber float64, secondNumber float64) ([]byte, error) {
	data := map[string]interface{}{
		"operationIndex":  comandIndex,
		"operationComand": comandString,
		"firstNumber":     firstNumber,
		"secondNumber":    secondNumber,
	}

	return json.Marshal(data)
}
