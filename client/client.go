package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type server struct {
	port string
	conn net.Conn
}

var servers = []*server{
	&server{"9000", nil},
	&server{"8080", nil},
	&server{"8000", nil},
}

var result map[int][]string = make(map[int][]string)
var ans []string

var ch chan int = make(chan int)

func main() {
	inputFile := os.Args[1]
	Solve(inputFile)
}

func Solve(inputFile string) (string, error) {

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	scanner := bufio.NewScanner(file)
	data := ""
	for scanner.Scan() {
		data += scanner.Text()
	}
	names := strings.Split(data, ",")

	active := 0
	for _, v := range servers {
		var err error
		v.conn, err = net.Dial("tcp", ":"+v.port)
		active++
		if err != nil {
			fmt.Println("server at addr ", v.port, " is not active")
			active--
		} else {
			fmt.Println("server at addr ", v.port, " is  active")
		}
	}
	fmt.Println("\n")

	dataset := len(names)/active + 1

	st := 0
	for i, v := range servers {

		if v.conn == nil {
			continue
		}

		en := st + dataset
		if en > len(names) {
			en = len(names)
		}

		input := strings.Join(names[st:en], ",")
		st = en
		go find_sorted(i, input)
	}

	for i := 0; i < active; i++ {
		<-ch
	}

	combine()
	str := strings.Join(ans, ",")
	str = strings.ReplaceAll(str, "\n", "")
	fmt.Println("\nResult:\n" + str)
	return str, nil
}

func find_sorted(server_num int, input string) {
	data := []byte(input + "\n")

	c := servers[server_num].conn
	_, err := c.Write(data)
	if err != nil {
		fmt.Println(err)
	}

	message, _ := bufio.NewReader(c).ReadString('\n')
	result[server_num] = strings.Split(message, ",")
	fmt.Println("from server at:", servers[server_num].port, "data:=", message)
	ch <- server_num
}

func combine() {

	if len(result) == 0 {
		return
	}

	ints := make([]int, 0)
	for i, _ := range result {
		ints = append(ints, i)
	}

	ans = result[ints[0]]

	for i := 1; i < len(ints); i++ {
		merge(ints[i])
	}

}

func merge(num int) {
	i := 0
	j := 0
	n := len(ans)
	m := len(result[num])
	res := []string{}

	for {
		if i >= n || j >= m {
			break
		}

		if ans[i] < result[num][j] {
			res = append(res, ans[i])
			i++
		} else {
			res = append(res, result[num][j])
			j++
		}
	}

	for {
		if i >= n {
			break
		}
		res = append(res, ans[i])
		i++
	}

	for {
		if j >= m {
			break
		}
		res = append(res, result[num][j])
		j++
	}

	ans = nil
	for _, v := range res {
		ans = append(ans, v)
	}
}
