package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From server: Hello I got your message "), addr)
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
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}
		fmt.Println(calculateInfos)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr)
	}
}

type CalculateInfos struct {
	OperationIndex  int
	OperationComand string
	FirstNumber     float32
	SecondNumber    float32
}
