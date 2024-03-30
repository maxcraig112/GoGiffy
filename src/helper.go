package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GetToken(filename string) (token string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(wd + `\\` + filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func CreateSearchWebhookEdit(tags []string, bqUIDS []string, index int) (webhookEdit *discordgo.WebhookEdit) {
	return &discordgo.WebhookEdit{
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Previous",
						Style:    discordgo.PrimaryButton,
						CustomID: "fd_previous",
					},
					discordgo.Button{
						Label:    "Next",
						Style:    discordgo.PrimaryButton,
						CustomID: "fd_next",
					},
					discordgo.Button{
						Label:    "Delete Gif",
						Style:    discordgo.DangerButton,
						CustomID: "fd_delete",
					},
				},
			},
		},
		Embeds: &[]*discordgo.MessageEmbed{{
			Title:       "Search by joe",
			Description: strconv.Itoa(index+1) + "/" + strconv.Itoa(len(bqUIDS)),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Tags",
					Value: strings.Join(tags, ", "),
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: GetPublicURLFromUID(bqUIDS[index]),
			},
		},
		},
	}
}

func sanitizeInput(input string) string {
	// Replace single quotes with two single quotes to escape them
	sanitized := strings.Replace(input, "'", "''", -1)
	return sanitized
}
