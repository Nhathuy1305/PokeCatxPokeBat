package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

var POKEMONS []Pokemon

// var exchangeMap = make(map[string]string)
// var exchangeJSON []byte
var MAP = make(map[string]string)

var PLAYERS []Player //for authentication
var ROWS, COLS = 10, 18
var BOARD = make([][]string, ROWS)
var CONNECTIONS = make(map[string]net.Conn)
var despawnQueues []string
var POKEMON_LOCATIONS = make(map[string]string)
var PLAYER_LOCATIONS = make(map[string]string)

var pokeBalls_P1 []Pokemon
var pokeBalls_P2 []Pokemon

var P1 string
var P2 string

// ////////////////////////////////////////////////////////////////////////////////////
// Load players.json to PLAYERS
func loadPlayers(filename string) []Player {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	var players []Player
	err = json.Unmarshal(bytes, &players)
	if err != nil {
		return nil
	}

	return players
}
func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
func loadPokemons(filename string) []Pokemon {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	var pokemons []Pokemon
	err = json.Unmarshal(bytes, &pokemons)
	if err != nil {
		return nil
	}

	return pokemons
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
func generateRandomPokemons(num int) map[string]string {
	pokemonLocations := make(map[string]string)
	for range num {
		for {
			spawnX := rand.Intn(ROWS)
			spawnY := rand.Intn(COLS)
			if BOARD[spawnX][spawnY] == "" {
				pokemonID := POKEMONS[rand.Intn(len(POKEMONS))].ID

				BOARD[spawnX][spawnY] = pokemonID

				despawnQueues = append(despawnQueues, strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY))
				pokemonLocations[strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY)] = pokemonID
				POKEMON_LOCATIONS[strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY)] = pokemonID
				break
			}

		}

	}
	return pokemonLocations
}

// ////////////////////////////////////////////////////////////////////////////////////
func handlePokemons() {
	spawnTicker1min := time.NewTicker(1 * time.Minute)
	despawnTicker5min := time.NewTicker(5 * time.Minute)
	NUMBERTOPROCESS := 3
	for {
		select {
		case <-spawnTicker1min.C:

			newPokemonLocations, err := json.Marshal(generateRandomPokemons(NUMBERTOPROCESS))
			checkError(err)

			for _, tcpConn := range CONNECTIONS {

				tcpConn.Write([]byte(newPokemonLocations))

			}
		case <-despawnTicker5min.C:
			despawnedPokemonLocations := make(map[string]string)
			for i := 0; i < NUMBERTOPROCESS; i++ {
				despawnedPokemonLocations[despawnQueues[i]] = ""
			}
			//CUT OFF THE FIRST 3 ELEMENTS or 50
			despawnQueues = despawnQueues[NUMBERTOPROCESS:]
			sent, _ := json.Marshal(despawnedPokemonLocations)

			for _, tcpConn := range CONNECTIONS {
				tcpConn.Write([]byte(sent))
			}
		}
	}

}
func battle(player1 string, player2 string) {
	fmt.Println(player1, " vs ", player2)

}

var player1Turn = true

// HandleConnection handles incoming client connections
func HandleConnection(conn net.Conn) {

	defer conn.Close()

	// Create a reader to read data from the connection
	reader := bufio.NewReader(conn)
	for {
		battleStatus := false
		// Read data from the connection
		player_coord, err := reader.ReadString('\n')
		fmt.Println("message: ", player_coord)

		player_coord = strings.TrimSpace(player_coord)

		//DISCONNECTED
		if err != nil {
			for name, connection := range CONNECTIONS {
				//tim connection hien tai
				if connection == conn {
					// tim ten
					for location, player := range PLAYER_LOCATIONS {
						if player == name {
							delete(PLAYER_LOCATIONS, location)
						}
					}
					delete(CONNECTIONS, name)
					quit := make(map[string]string)
					quit[strings.TrimSpace(name)] = "quit"
					sentQuit, _ := json.Marshal(quit)
					for _, connection := range CONNECTIONS {

						connection.Write([]byte(sentQuit))
					}
					fmt.Println(strings.TrimSpace(name) + " is disconnected")
					break

				}

			}
			return

		}
		//battling message
		if strings.Split(player_coord, "-")[0] == "battle" {
			P := strings.Split(player_coord, "-")[1]

			main_message := strings.TrimSpace(strings.Split(player_coord, "-")[2])

			turn := make(map[string]string)
			//submit pokemons
			if isNumber(main_message) {
				fmt.Println("Submitting pokemon...")
				pokeId := main_message
				for pokeIndex := 0; pokeIndex < len(POKEMONS); pokeIndex++ {
					if POKEMONS[pokeIndex].ID == pokeId {
						if P == P1 {
							pokeBalls_P1 = append(pokeBalls_P1, POKEMONS[pokeIndex])
						} else if P == P2 {
							pokeBalls_P2 = append(pokeBalls_P2, POKEMONS[pokeIndex])
						}
						break
					}

				}

				if len(pokeBalls_P1) == 3 && len(pokeBalls_P2) == 3 {
					fmt.Println("Let's Battle!")

					speed_P1, _ := strconv.Atoi(pokeBalls_P1[0].Stats["Speed"])
					speed_P2, _ := strconv.Atoi(pokeBalls_P2[0].Stats["Speed"])
					wait := make(map[string]string)
					wait["battle"] = "wait"
					sentWait, _ := json.Marshal(wait)
					if speed_P1 >= speed_P2 {
						fmt.Println("P1 turn")
						//sent message to p1 to go first
						turn["battle"] = P1
						sentTurn, _ := json.Marshal(turn)
						CONNECTIONS[P1].Write([]byte(sentTurn))

						CONNECTIONS[P2].Write([]byte(sentWait))
					} else {
						fmt.Println("P2 turn")
						//sent message to p2 to go first
						turn["battle"] = P2
						sentTurn, _ := json.Marshal(turn)
						CONNECTIONS[P2].Write([]byte(sentTurn))

						CONNECTIONS[P1].Write([]byte(sentWait))

						player1Turn = false
					}
				}
				//not submitting means movement
			} else if len(strings.Split(main_message, "*")) == 2 {
				currentPokemon, _ := strconv.Atoi(strings.Split(main_message, "*")[0])
				action := strings.TrimSpace(strings.Split(main_message, "*")[1])
				if player1Turn {
					if action == "attack" {
						attackedPokemonIndex := rand.Intn(3)
						fmt.Println("attack")
						P2_HP, _ := strconv.Atoi(pokeBalls_P2[attackedPokemonIndex].Stats["HP"])

						spec_attk := rand.Intn(2)

						if spec_attk == 1 {
							P1_Damage, _ := strconv.Atoi(pokeBalls_P1[currentPokemon].Stats["Sp Atk"])
							P2_Defense, _ := strconv.Atoi(pokeBalls_P2[attackedPokemonIndex].Stats["Sp Def"])
							P2_HP -= (P1_Damage - P2_Defense)
						} else {
							P1_Damage, _ := strconv.Atoi(pokeBalls_P1[currentPokemon].Stats["Attack"])
							P2_Defense, _ := strconv.Atoi(pokeBalls_P2[attackedPokemonIndex].Stats["Defense"])
							P2_HP -= (P1_Damage - P2_Defense)
						}

						if P2_HP <= 0 {
							pokeBalls_P2 = append(pokeBalls_P2[:attackedPokemonIndex], pokeBalls_P2[attackedPokemonIndex+1:]...)
						}

						attackMsg := make(map[string]string)

						attackMsg["battle"] = "attacked-" + strconv.Itoa(P2_HP) + "-" + strconv.Itoa(attackedPokemonIndex)
						sentAttackMsg, _ := json.Marshal(attackMsg)

						CONNECTIONS[P2].Write([]byte(sentAttackMsg))
					}
				} else {
					if action == "attack" {
						attackedPokemonIndex := rand.Intn(3)
						fmt.Println("attack")
						P2_HP, _ := strconv.Atoi(pokeBalls_P1[attackedPokemonIndex].Stats["HP"])

						spec_attk := rand.Intn(1)

						if spec_attk == 1 {
							P1_Damage, _ := strconv.Atoi(pokeBalls_P2[currentPokemon].Stats["Sp Atk"])
							P2_Defense, _ := strconv.Atoi(pokeBalls_P1[attackedPokemonIndex].Stats["Sp Def"])
							P2_HP -= (P1_Damage - P2_Defense)
						} else {
							P1_Damage, _ := strconv.Atoi(pokeBalls_P2[currentPokemon].Stats["Attack"])
							P2_Defense, _ := strconv.Atoi(pokeBalls_P1[attackedPokemonIndex].Stats["Defense"])
							P2_HP -= (P1_Damage - P2_Defense)
						}

						if P2_HP <= 0 {
							pokeBalls_P1 = append(pokeBalls_P2[:attackedPokemonIndex], pokeBalls_P2[attackedPokemonIndex+1:]...)
						}

						attackMsg := make(map[string]string)

						attackMsg["battle"] = "attacked-" + strconv.Itoa(P2_HP) + "-" + strconv.Itoa(attackedPokemonIndex)
						sentAttackMsg, _ := json.Marshal(attackMsg)

						CONNECTIONS[P1].Write([]byte(sentAttackMsg))
					}
					if action == "switch" {

					}
				}

			}

		} else {
			// iterate through connections
			for name, connection := range CONNECTIONS {
				name = strings.TrimSpace(name)
				// find matching name
				if connection == conn {
					//find the name corresponding to connection address
					for playerLocation, player := range PLAYER_LOCATIONS {
						player = strings.TrimSpace(player)
						if player == name {

							//if player hits pokemon location
							if pokemonIndex, exists := POKEMON_LOCATIONS[strings.TrimSpace(player_coord)]; exists {
								fmt.Println("CATCHING...")
								//send pokemon id to the player through "conn"
								// catched: "index"
								catched := make(map[string]string)
								catched[strings.TrimSpace(name)] = pokemonIndex

								sentCatched, _ := json.Marshal(catched)
								conn.Write(sentCatched)
								//remove location of the pokemon from BOARD and POKEMON_LOCATIONS
								pokemonX, _ := strconv.Atoi(strings.Split(playerLocation, "-")[0])
								pokemonY, _ := strconv.Atoi(strings.Split(playerLocation, "-")[1])
								BOARD[pokemonX][pokemonY] = ""
								delete(POKEMON_LOCATIONS, strings.TrimSpace(player_coord))
								fmt.Println(POKEMON_LOCATIONS)
								for _, tcpConn := range CONNECTIONS {
									if tcpConn != conn {
										pokemonGone := make(map[string]string)
										pokemonGone[strings.TrimSpace(player_coord)] = ""
										sentPokemonGone, _ := json.Marshal(pokemonGone)
										tcpConn.Write([]byte(sentPokemonGone))
									}
								}
							} else if enemy_name, exists := PLAYER_LOCATIONS[strings.TrimSpace(player_coord)]; exists {
								enemy_name = strings.TrimSpace(enemy_name)
								//BATTLE
								fmt.Println("battle")
								//send battle noti to player
								battleInfo := make(map[string]string)
								battleInfo["battle"] = enemy_name
								sentBattleInfo, _ := json.Marshal(battleInfo)
								conn.Write([]byte(sentBattleInfo))

								//send to the enemy
								battledInfo := make(map[string]string)
								battledInfo["battle"] = name
								sentBattledInfo, _ := json.Marshal(battledInfo)
								CONNECTIONS[enemy_name].Write([]byte(sentBattledInfo))

								battleStatus = true
								pokeBalls_P1 = []Pokemon{}
								pokeBalls_P2 = []Pokemon{}
								P1 = name
								P2 = enemy_name
								// start new routine for
								// go battle(name, enemy_name)

								fmt.Println(P1, " vs ", P2)
							}

							//remove previous location to make animation
							if !battleStatus {
								delete(PLAYER_LOCATIONS, playerLocation)
							}
						}
					}
					if !battleStatus {
						PLAYER_LOCATIONS[player_coord] = name
						fmt.Println(PLAYER_LOCATIONS)
						for _, tcpConn := range CONNECTIONS {
							// if conn != tcpConn {
							sentPLAYER_LOCATIONS, _ := json.Marshal(PLAYER_LOCATIONS)
							tcpConn.Write([]byte(sentPLAYER_LOCATIONS))

							// }
						}
					}
				}
			}

			// Echo the message back to the client

		}
	}
}
func handle(conn net.Conn) {
	// Verfify username and password
	infoReader := bufio.NewReader(conn)

	// Get username
	username, err := infoReader.ReadString('\n')
	checkError(err)

	// Get password
	password, err := infoReader.ReadString('\n')
	checkError(err)

	if verifyPlayer(strings.TrimSpace(username), strings.TrimSpace(password), PLAYERS) {
		// Handle the connection in a new goroutine
		_, err := conn.Write([]byte("successful"))
		initialPokemons := ""
		for i := range 3 {
			initialPokemons += strconv.Itoa(rand.Intn(len(POKEMONS)))
			if i < 2 {
				initialPokemons += "-"
			}
		}
		conn.Write([]byte(initialPokemons))
		checkError(err)
		time.Sleep(22 * time.Second)
		CONNECTIONS[strings.TrimSpace(username)] = conn
		fmt.Println(CONNECTIONS)

		//send pokemon locations
		sentPOKEMON_LOCATIONS, _ := json.Marshal(POKEMON_LOCATIONS)
		conn.Write([]byte(sentPOKEMON_LOCATIONS))
		for {
			playerX := rand.Intn(ROWS)
			playerY := rand.Intn(COLS)
			if BOARD[playerX][playerY] == "" {
				BOARD[playerX][playerY] = username
				PLAYER_LOCATIONS[strconv.Itoa(playerX)+"-"+strconv.Itoa(playerY)] = username

				// ini_coord := make(map[string]string)
				// ini_coord[strconv.Itoa(playerX)+"-"+strconv.Itoa(playerY)] = username

				// //send initial player's location
				// sentInital_Coord, _ := json.Marshal(ini_coord)
				// conn.Write([]byte(sentInital_Coord))
				break
			}
		}

		//send player locations
		sentPLAYER_LOCATIONS, _ := json.Marshal(PLAYER_LOCATIONS)
		for _, connection := range CONNECTIONS {
			connection.Write([]byte(sentPLAYER_LOCATIONS))
		}

		go HandleConnection(conn)
	} else {
		conn.Write([]byte("failed"))
	}
}
func main() {
	// Set up BOARD
	for i := range BOARD {
		BOARD[i] = make([]string, COLS)
	}
	// Load pokemon.json to array POKEMONS
	POKEMONS = loadPokemons("pokedex.json")

	//Load players.json to array PLAYERS
	PLAYERS = loadPlayers("players.json")
	generateRandomPokemons(3)
	fmt.Println(POKEMON_LOCATIONS)
	go handlePokemons()

	// Start listening for incoming connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	/// spawn pokem
	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handle(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
