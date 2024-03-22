package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func getAuthenticationToken() (token string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(wd + `\\token.txt`)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {
	token, err := getAuthenticationToken()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Got Authentication Token")

	discordClient, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Discord Client Created")

	fmt.Println(discordClient)
}
