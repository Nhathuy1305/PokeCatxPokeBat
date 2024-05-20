Pokedex/<br>
├── cmd/<br>
│   ├── crawler/<br>
│   │   └── main.go          # Entry point to start the crawler.<br>
│   └── server/<br>
│       └── main.go          # Entry point to start the game server.<br>
├── pkg/<br>
│   ├── crawler/<br>
│   │   ├── crawler.go       # Logic to crawl and parse Pokémon data.<br>
│   │   └── parser.go        # Parses HTML and constructs Pokémon data.<br>
│   ├── battle/<br>
│   │   ├── battle.go        # Core battle mechanics and rules.<br>
│   │   └── calculator.go    # Calculates damage and battle outcomes.<br>
│   ├── storage/<br>
│   │   ├── storage.go       # Manages read/write operations for JSON files.<br>
│   │   └── file_lock.go     # Implements locking mechanisms for file operations.<br>
│   ├── pokemon/<br>
│   │   ├── pokemon.go       # Defines the Pokémon struct and methods.<br>
│   │   └── level.go         # Handles leveling and experience calculations.<br>
│   └── world/<br>
│       ├── world.go         # Manages world state, including player and Pokémon positions.<br>
│       └── spawn.go         # Handles the logic for Pokémon spawning and despawning.<br>
├── data/<br>
│   ├── pokedex.json         # Stores detailed information about each Pokémon.<br>
│   └── players.json         # Stores player data, including their Pokémon collections.<br>
├── config/<br>
│   └── config.go            # Contains configuration variables and constants.<br>
└── main.go                  # Main application entry point that initializes and starts components based on command-line arguments.<br>


#### `cmd/`

*   **crawler/main.go**: This is the main entry point for the web crawler module. It sets up and initiates the scraping process to populate the Pokémon database.
*   **server/main.go**: Starts the server that manages game sessions, player connections, and real-time interactions for battles and world exploration.

#### `pkg/`

*   **crawler/crawler.go**: Contains the primary functions to initiate web scraping, handle URL navigation, and manage data extraction.
*   **crawler/parser.go**: Dedicated to parsing the scraped HTML data and transforming it into structured Pokémon data that can be serialized into JSON.
*   **battle/battle.go**: Implements the turn-based combat logic including managing turns, determining move order based on Pokémon speed, and processing player actions.
*   **battle/calculator.go**: Responsible for computing battle-related calculations, such as attack damage, effects of special moves, and battle outcomes.
*   **storage/storage.go**: Provides functions to load and save data to/from JSON files, ensuring that data is consistently stored and retrieved.
*   **storage/file\_lock.go**: Implements a mechanism to lock files during read/write operations to prevent race conditions in a multi-user environment.
*   **pokemon/pokemon.go**: Defines the Pokémon structure with attributes and methods necessary for handling Pokémon behavior and states.
*   **pokemon/level.go**: Manages the experience and leveling system for Pokémon, including experience accumulation and level-up mechanics.
*   **world/world.go**: Manages the virtual game world's state, tracking player locations, movements, and interactions within the world.
*   **world/spawn.go**: Controls how and where Pokémon appear in the game world, including the timing and conditions of their spawning and despawning.

#### `data/`

*   **pokedex.json** and **players.json**: Serve as the primary storage for static and dynamic game data, respectively.

#### `config/`

*   **config.go**: Centralizes all configuration settings, making it easy to adjust game parameters such as spawn rates, world dimensions, and other gameplay variables.

#### `main.go`

*   This file acts as a dispatcher based on command-line arguments to either run the crawler or start the server, allowing for modular use of the application components.
