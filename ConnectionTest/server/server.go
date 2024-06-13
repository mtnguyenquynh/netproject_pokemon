package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)


func main() {
    fmt.Println("Server is listening on port 8080...")

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        return
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err.Error())
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    InitData()


    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    // Read client's name
    fmt.Fprintf(writer, "What is your name?\n")
    writer.Flush()
    name, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading client's name:", err.Error())
        return
    }
    name = strings.TrimSpace(name)

    // Initialize user input
    userInput := &UserInput{username: name}

    // Choose name based on client's UserInput
    ChooseName(userInput, writer)
	fmt.Println(userInput)
	fmt.Println("User", userInput.username, "has joined.")


  // Prepare Cynthia's team (example with hardcoded values)
    venusaur := NewPokemon("Venusaur", true)
    charmeleon := NewPokemon("Charmeleon", true)
    wartortle := NewPokemon("Wartortle", true)
    blastoise := NewPokemon("Blastoise", true)
    caterpie := NewPokemon("Caterpie", true)
    bulbasaur := NewPokemon("Bulbasaur", true) 

    cynthiasTeam := []*Pokemon{venusaur, charmeleon, wartortle, blastoise, caterpie, bulbasaur}

    //Simulate battle with Cynthia
    cynthiasInput := &UserInput{"Cynthia", "", venusaur, cynthiasTeam, "", true, false}

    // Send team selection prompt
    fmt.Fprintf(writer, "Type a number and press ENTER to choose an option\n")
    writer.Flush()
  // Receive team choice
  choiceStr, err := reader.ReadString('\n')
  if err != nil {
      fmt.Println("Error reading choice:", err.Error())
      return
  }
  choiceStr = strings.TrimSpace(choiceStr)



  // Process team choice
ChooseTeam(userInput, choiceStr, reader, writer)




    fmt.Fprintf(writer, "\n[[ BATTLE ]] Starting a battle\n")
    writer.Flush()
    Battle(userInput, cynthiasInput)

    fmt.Fprintf(writer, "\nBattle concluded. Closing connection.\n")
    writer.Flush()
}
