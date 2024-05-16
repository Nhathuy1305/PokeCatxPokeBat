### Directory and File Description

*   **cmd/**: This directory contains the entry points for the application.
    
    *   **server/**: Contains the main server application code.
        *   `main.go`: Entry point for starting the server.
    *   **client/**: Contains the client application code.
        *   `main.go`: Entry point for starting the client.
*   **internal/**: Contains the core application code.
    
    *   **battle/**: Handles the battle logic.
        *   `battle.go`: Main battle logic.
        *   `attack.go`: Handles attack logic.
        *   `player.go`: Handles player-related actions in battles.
    *   **catching/**: Handles the catching logic.
        *   `pokeworld.go`: Handles the pokeworld logic.
        *   `player.go`: Handles player-related actions in the pokeworld.
    *   **pokedex/**: Manages the pokedex and related functionalities.
        *   `crawler.go`: Web crawler to build the pokedex.
        *   `pokedex.go`: Handles pokedex data management.
        *   `types.go`: Defines types related to the pokedex.
    *   **pokemon/**: Manages pokemon data and related functionalities.
        *   `pokemon.go`: Handles pokemon data management.
        *   `types.go`: Defines types related to pokemon.
    *   **network/**: Manages network communications.
        *   `tcp_server.go`: Handles TCP server logic.
        *   `udp_server.go`: Handles UDP server logic.
        *   `client.go`: Handles client-side network logic.
*   **config/**: Configuration files and related code.
    
    *   `config.go`: Configuration management code.
    *   `config.yaml`: Configuration file.
*   **data/**: Contains data files.
    
    *   `pokedex.json`: The pokedex database.
    *   **player\_pokemons/**: Directory to store player-specific pokemon lists.
        *   `player_<id>.json`: JSON file containing a player's captured pokemon.
*   **scripts/**: Contains utility scripts.
    
    *   `generate_pokedex.go`: Script to generate the pokedex JSON file.
*   **pkg/**: Contains reusable utility packages.
    
    *   **utils/**: General utility functions.
        *   `json_utils.go`: Utility functions for working with JSON.
        *   `math_utils.go`: Utility functions for mathematical operations.
    *   **logger/**: Logging utility.
        *   `logger.go`: Logger utility functions.
*   **README.md**: Documentation for the project.