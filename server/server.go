package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Pokemon struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Types []string          `json:"types"`
	Stats map[string]string `json:"stats"`
	Exp   string            `json:"exp"`
}

type Player struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	PokeBalls []Pokemon `json:"pokeBalls"`
}

func loadPlayers(filename string) ([]Player, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var players []Player
	err = json.Unmarshal(bytes, &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}

func savePlayers(filename string, players []Player) error {
	bytes, err := json.MarshalIndent(players, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
func verifyPlayer(username, password string, players []Player) bool {

	for _, user := range players {
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			return err == nil
		}
	}
	return false
}
func handleClient(conn *net.UDPConn, players []Player) {
	//BUFFER DECLARATION
	var buf [512]byte

	//RESOLVE UDP ADDRESS
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	//SEND DATE TIME TO USER
	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)

	//READNAME
	n, err := conn.Read(buf[0:])
	if err != nil {
		return
	}
	playerInfo := string(buf[0:n])

	// Split the input string by the delimiter "//"
	parts := strings.Split(playerInfo, "//")

	// Trim spaces, tabs, and newlines from each part
	username := strings.TrimSpace(parts[0])
	password := strings.TrimSpace(parts[1])
	if verifyPlayer(username, password, players) {
		fmt.Println(addr, "-", "\"", username, "\"", "-", "connected!")
	}

}

func main() {
	exchangeJSON := make(map[string]string)
	rows, cols := 100, 100
	//Initialize a board on client side
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
	}

	// Load pokemon.json
	file, err := os.Open("pokedex.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Decode the JSON data into an array of Pokemon structs
	var pokemons []Pokemon
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pokemons); err != nil {
		fmt.Println("Error:", err)
		return
	}
	//randomize 50 pokemons

	for range 50 {
		newPokemonId := rand.Intn(len(pokemons)) + 1
		spawnX := rand.Intn(rows)
		spawnY := rand.Intn(cols)
		board[spawnX][spawnY] = strconv.Itoa(newPokemonId)

		exchangeJSON[strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY)] = strconv.Itoa(newPokemonId)
	}

	fmt.Println(board)

	//LOAD players.json for verification
	players, err := loadPlayers("players.json")
	if err != nil {
		fmt.Println("Error loading users:", err)
		return
	}
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		handleClient(conn, players)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
