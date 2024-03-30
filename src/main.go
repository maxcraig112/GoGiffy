package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	token          = flag.String("token", "", "Bot access token")
	AppID          = flag.String("app", "", "Application ID")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session
var ctx context.Context

func init() { flag.Parse() }

func init() {
	ctx = context.Background()

	var err error
	token, err := GetToken("token.txt")
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

// func init() {
// 	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
// 			h(s, i)
// 		}
// 	})
// }

var (
	project_ID = "gogiffy"
	err        error
	// integerOptionMinValue          = 1.0
	// dmPermission                   = false
	// defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "giffy",
			Description: "Command for returning useful information about the bot",
		},
		{
			Name:        "search",
			Description: "Command for search for gifs with a particular tag",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tags",
					Description: "tags that you want to search by to find a gif",
					Required:    true,
				},
			},
		},
		{
			Name:        "stats",
			Description: "Command for returning interesting stats about the bot",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"giffy": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
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
		"search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			option := optionMap["tags"]
			var tags = strings.Split(option.StringValue(), ",")

			//remove spaces on either side
			for i, str := range tags {
				tags[i] = strings.TrimSpace(str)

				if len(str) > 0 && str[len(str)-1] == 's' {
					// Add a new string without 's'
					tags = append(tags, str[:len(str)-1])
				} else {
					// Add the original string with an s
					tags = append(tags, str+"s")
				}
			}
			// Access options in the order provided by the user.
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Currently searching for gifs with the tags `" + strings.Join(tags, ", ") + "`, This may take up 10 seconds",
				},
			})

			// Or convert the slice into a map

			bqUIDS, err := GetUIDFromTags(tags)
			if err != nil {
				panic(err)
			}
			if len(bqUIDS) != 0 {
				index := 0
				_, err = s.InteractionResponseEdit(i.Interaction, CreateSearchWebhookEdit(tags, bqUIDS, index))
				if err != nil {
					fmt.Println(err.Error())
				}
				bigquerySearch := BigquerySearch{
					message_id:  i.Interaction.ID,
					query_time:  time.Now(),
					token:       i.Token,
					tags:        tags,
					index:       index,
					bucket_uids: bqUIDS,
				}
				err = AddBigQuerySearch(bigquerySearch)
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				content := "Sorry! there are no gifs with the tags `" + strings.Join(tags, ", ") + "`, please try with other search terms."
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &content,
				})
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		},
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			stats, err := GetUrlStats()
			if err != nil {
				panic(err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{{
						Title: "Giffy Stats",
						Description: `Giffy bot dynamically stores all unique gifs it comes across. Here are some interesting statistics!
						`,
						Image: &discordgo.MessageEmbedImage{URL: "https://media1.tenor.com/m/oVakhxonutgAAAAC/christian-bale-american-psycho.gif"},
						Author: &discordgo.MessageEmbedAuthor{
							URL:     "https://github.com/maxcraig112",
							Name:    "Max.imilian",
							IconURL: "https://media.discordapp.net/attachments/846175975560839178/950204054754177124/cool_obama.jpg?ex=660a0e7c&is=65f7997c&hm=45d7bc572e2922799019b15e73210daf8dc98fdd6b8e85471e96e510e120dba2&",
						},

						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Total gifs stored",
								Value: strconv.Itoa(stats.total_gifs),
							},
							{
								Name:  "Gifs without text",
								Value: strconv.Itoa(stats.gifs_without_text),
							},
							{
								Name:  "Gifs with text",
								Value: strconv.Itoa(stats.gifs_with_text),
							},
							{
								Name:  "Number of unique tags",
								Value: strconv.Itoa(stats.unique_tags_count),
							},
						},
					}},
				},
			})
		},
	}

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fd_previous": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bigquerySearch, err := GetBigQuerySearch(i.Message.Interaction.ID)
			if err != nil {
				fmt.Println(err)
			}
			//decrement the index
			numberOfResults := len(bigquerySearch.bucket_uids)
			bigquerySearch.index = (bigquerySearch.index - 1 + numberOfResults) % numberOfResults
			i.Interaction.ID = i.Message.Interaction.ID
			i.Interaction.Token = bigquerySearch.token
			_, err = s.InteractionResponseEdit(i.Interaction, CreateSearchWebhookEdit(bigquerySearch.tags, bigquerySearch.bucket_uids, bigquerySearch.index))
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			AddBigQuerySearch(bigquerySearch)
		},
		"fd_next": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bigquerySearch, err := GetBigQuerySearch(i.Message.Interaction.ID)
			if err != nil {
				fmt.Println(err)
			}
			//decrement the index
			numberOfResults := len(bigquerySearch.bucket_uids)
			bigquerySearch.index = (bigquerySearch.index + 1) % numberOfResults
			i.Interaction.ID = i.Message.Interaction.ID
			i.Interaction.Token = bigquerySearch.token
			_, err = s.InteractionResponseEdit(i.Interaction, CreateSearchWebhookEdit(bigquerySearch.tags, bigquerySearch.bucket_uids, bigquerySearch.index))
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			// fmt.Println(bigquerySearch.bucket_uids)
			// fmt.Println(GetPublicURLFromUID(bigquerySearch.bucket_uids[bigquerySearch.index]))
			AddBigQuerySearch(bigquerySearch)

		},
		"fd_delete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Println("Original Message ID " + i.Message.Interaction.ID)

		},
	}
)

func processFile(fileName string) error {
	errorList := make([]string, 0)
	processedStrings := make([]string, 0)

	// Open the file
	file, err := os.Open(fileName) // Replace "your_file.txt" with the path to your text file
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		processedStrings = append(processedStrings, line)
		fmt.Println(strconv.Itoa(len(processedStrings)) + "/10")
		// Process each string and handle errors
		if _, err := ProcessUrls([]string{line}); err != nil {
			errorList = append(errorList, line)
		}
	}

	// Check for errors in scanning the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return err
	}

	// Save the errors to a file
	errorFile, err := os.Create("errors.txt")
	if err != nil {
		fmt.Println("Error creating error file:", err)
		return err
	}
	defer errorFile.Close()

	// Write errors to the file
	writer := bufio.NewWriter(errorFile)
	for _, err := range errorList {
		_, err := writer.WriteString(err + "\n")
		if err != nil {
			fmt.Println("Error writing to error file:", err)
			return err
		}
	}
	// Flush the buffer to ensure all data is written to the file
	writer.Flush()

	return nil
}
func main() {
	// err := processFile("archivedgifs.txt")
	// if err != nil {
	// 	panic(err)
	// }

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:

			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	if err != nil {
		panic(err)
	}

	s.AddHandler(messageCreate)

	s.Identify.Intents = discordgo.IntentsAll

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	allCmds, _ := s.ApplicationCommands(*AppID, *GuildID)
	for _, v := range allCmds {
		fmt.Println(v.Name)
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Adding commands...")
	for _, v := range commands {
		s.ApplicationCommandCreate(*AppID, *GuildID, v)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")

	allCmds, _ = s.ApplicationCommands(*AppID, *GuildID)
	for _, v := range allCmds {
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
	//Gifs can come in 2 forms, through an attachment, or through an embed
	var urls []string
	fmt.Println(m.Attachments)
	fmt.Println(m.Embeds)
	if len(m.Attachments) > 0 {
		for _, attachment := range m.Attachments {
			var url = attachment.URL
			if UrlIsGif(url) {
				urls = append(urls, attachment.URL)
			}
		}
	}

	if len(m.Embeds) > 0 {
		for _, embed := range m.Embeds {
			var url = embed.URL
			if UrlIsGif(url) {
				urls = append(urls, url)
			}
		}
	}

	if len(urls) > 0 {
		fmt.Println(fmt.Sprint(len(urls)) + " gifs found")
		_, err := ProcessUrls(urls)
		if err != nil {
			panic(err)
		}

		for _, url := range urls {
			BigqueryURL, err := GetBigQueryURL(url)
			if err != nil {
				panic(err)
			}

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(`Gif Stored
			Url: %s
			Contains Caption: %s
			Text: %s
			UUID: %s
			Bucket URL: %s`, BigqueryURL.url, strconv.FormatBool(BigqueryURL.contains_caption), BigqueryURL.text, BigqueryURL.bucket_uid, BigqueryURL.GetPublicURL()))
		}
	}

	//If the message contains attachments
	if (len(m.Attachments)) > 0 {
		fmt.Println(m.Attachments)
	}
}
