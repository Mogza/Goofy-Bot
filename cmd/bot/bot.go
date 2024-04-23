package bot

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	dexApi "github/Mogza/Goofy-Bot/cmd/token"
	"log"
	"os"
	"os/signal"
	"strings"
)

var s *discordgo.Session

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "token",
		Description: "Get the Token informations",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "addr",
				Description: "Token address",
				Required:    true,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"token": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var tokenAddress string

		options := i.ApplicationCommandData().Options
		for _, opt := range options {
			tokenAddress = opt.StringValue()
		}

		message, imageUrl := dexApi.GetToken(tokenAddress)

		responseContent := []*discordgo.MessageEmbed{
			{
				Title:       "Token Overview :",
				Description: message,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: imageUrl,
				},
			},
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: responseContent,
			},
		})
		checkError(err, "Error while responding to `/token`")
	},
}

func checkError(e error, message string) {
	if e != nil {
		log.Fatalln(message, ":", e)
	}
}

func InitBot() {
	var err error

	flag.Parse()

	// Initialize session
	s, err = discordgo.New("Bot " + "MTIzMTc1NTkzNDIzNDE4MTY2NQ.GQtFzB.VRQQrjz4fV05a2AB_cSG505PDWKWm1iyLyJmq4")
	checkError(err, "Error while creating the bot")

	// Handling commands
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if len(m.Content) > 30 && strings.Count(m.Content, " ") == 0 {
			var tokenAddress string

			tokenAddress = m.Content

			message, imageUrl := dexApi.GetToken(tokenAddress)

			responseContent := []*discordgo.MessageEmbed{
				{
					Title:       "Token Overview :",
					Description: message,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: imageUrl,
					},
				},
			}

			_, err = s.ChannelMessageSendEmbedsReply(m.ChannelID, responseContent, m.Reference())
			checkError(err, "Error while responding to `/token`")
		}
	})
}

func Run() {
	err := s.Open()
	checkError(err, "Error while running the bot")
	defer func(s *discordgo.Session) {
		err := s.Close()
		checkError(err, "Error while closing the bot")
	}(s)

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		checkError(err, "Cannot create command")
		registeredCommands[i] = cmd
	}

	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
