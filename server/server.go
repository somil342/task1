package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strings"
)

func main() {
	port := os.Args[1]
	fmt.Println("Listening on port : ", port)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("port", err)
	}
	defer ln.Close()

	for {

		fmt.Println("waiting for conn")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept", err)
		}
		scanner := bufio.NewScanner(conn)

		for scanner.Scan() {
			fmt.Println("->")
			data := scanner.Text()
			strs := strings.Split(data, ",")
			for i, v := range strs {
				strs[i] = strings.Trim(v, " ")
			}
			sort.Strings(strs)
			io.WriteString(conn, strings.Join(strs, ",")+"\n")
		}
		conn.Close()
		fmt.Println("client closing conn")
	}

}
