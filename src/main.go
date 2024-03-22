package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	token          = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

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

func init() { flag.Parse() }

func init() {
	var err error
	token, err := getAuthenticationToken()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Got Authentication Token")

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Discord Client Created")
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "giffy",
			Description: "Command for returning useful information about the bot",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"giffy": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
					Embeds: []*discordgo.MessageEmbed{{
						Title: "Who am I?",
						Description: `I am giffy! A multipurpose discord bot designed to allow the manipulation, tagging, archiving and retrieval of gifs!

						**Help**
						If you need help with commands, type ` + "`/help`" + `
						
						**Searching**
						All unique gifs sent in a server I'm in will automatically be catagorised and archived!
						
						This means that if you ever have a caption gif that you're dying to use, but can't find, you can simply type ` + "`/search [tag]`" + `and I'll try my best to find it for you!`,
						Image: &discordgo.MessageEmbedImage{URL: "https://c.tenor.com/oylHwLtwhbsAAAAC/gif-jif.gif"},
						Author: &discordgo.MessageEmbedAuthor{
							URL:     "https://github.com/maxcraig112",
							Name:    "Max.imilian",
							IconURL: "https://media.discordapp.net/attachments/846175975560839178/950204054754177124/cool_obama.jpg?ex=660a0e7c&is=65f7997c&hm=45d7bc572e2922799019b15e73210daf8dc98fdd6b8e85471e96e510e120dba2&",
						},
					}},
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(messageCreate)

	s.Identify.Intents = discordgo.IntentsAll

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	fmt.Print(m.Attachments)
}
