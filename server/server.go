package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, response string) {
	var msg string = "From server: Hello I got your message, your result is " + response

	data := map[string]interface{}{
		"msg":    msg,
		"result": response,
	}

	jsonData, err := json.Marshal(data)

	_, err = conn.WriteToUDP(jsonData, addr)

	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func main() {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		n, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		var calculateInfos CalculateInfos

		err = json.Unmarshal(p[:n], &calculateInfos)

		if err != nil {
			fmt.Printf("could not Unmarshal json: %s\n", err)
			return
		}
		res, err := calculate(calculateInfos)
		if err != nil {
			go sendResponse(ser, remoteaddr, "error")
		}
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr, fmt.Sprintf("%f", res))
	}
}

type CalculateInfos struct {
	OperationIndex  int
	OperationComand string
	FirstNumber     float64
	SecondNumber    float64
}

func calculate(info CalculateInfos) (float64, error) {
	var result float64
	switch info.OperationIndex {
	case 1:
		result = info.FirstNumber + info.SecondNumber
	case 2:
		result = info.FirstNumber - info.SecondNumber
	case 3:
		if info.SecondNumber == 0.0 {
			return 0.0, errors.New("error: you tried to divide by zero.")
		}
		result = info.FirstNumber / info.SecondNumber
	case 4:
		result = info.FirstNumber * info.SecondNumber
	}

	return result, nil

}
