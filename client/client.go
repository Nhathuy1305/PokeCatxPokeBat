package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func drawBoard(board [][]string) {

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	// Function to generate a horizontal line
	horizontalLine := func(length int) string {
		return "+" + strings.Repeat("---+", length)
	}

	for _, row := range board {
		// Print horizontal line before each row
		fmt.Println(horizontalLine(len(row)))

		// Print cell values or empty spaces
		for _, cell := range row {
			if cell == "" {
				fmt.Print("|   ")
			} else {
				fmt.Printf("| %s ", cell)
			}
		}
		fmt.Println("|")
	}
	// Print the final horizontal line after all rows
	fmt.Println(horizontalLine(len(board[0])))
}
func main() {
	rows, cols := 1000, 1000
	//Initialize a board on client side
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
	}

	//READ INPUT AS PORT
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}

	//CREATE UDP ADDRESS
	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	//CONNECT TO SERVER USING UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	//GREETING SERVER
	_, err = conn.Write([]byte("Hello, Server!"))
	checkError(err)

	//READ FROM SERVER
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	checkError(err)

	fmt.Println(string(buf[0:n]))

	//SEND NAME TO SERVER
	var username string
	fmt.Print("username: ")
	fmt.Scanln(&username)

	var password string
	fmt.Print("password: ")
	fmt.Scanln(&password)

	checkError(err)

	//send
	_, err = conn.Write([]byte(username + "//" + string(password)))
	checkError(err)

	// var message string
	// for {
	// 	fmt.Print("Enter your message: ")
	// 	fmt.Scanln(&message)
	// 	_, err = conn.Write([]byte(message))
	// 	checkError(err)

	// 	n, err := conn.Read(buf[0:])
	// 	checkError(err)

	// 	fmt.Println(string(buf[0:n]))
	// }

	arr := [][]string{
		{"?", "", "?", ""},
		{"", "?", "", "?"},
		{"?", "", "?", ""},
		{"?", "", "?", ""},
		{"", "?", "", "?"},
		{"?", "", "?", ""},
	}

	// Draw the array
	drawBoard(arr)
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
