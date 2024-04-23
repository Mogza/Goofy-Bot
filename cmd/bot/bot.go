package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	dexApi "github/Mogza/Goofy-Bot/cmd/token"
	"github/Mogza/Goofy-Bot/cmd/utils"
	"log"
	"os"
	"os/signal"
	"strings"
)

var s *discordgo.Session

// Commands list
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

// Command handler function
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
		utils.CheckError(err, "Error while responding to `/token`")
	},
}

// Detect token handler function
var botHandlers = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
	"detectToken": func(s *discordgo.Session, m *discordgo.MessageCreate) {
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

			_, err := s.ChannelMessageSendEmbedsReply(m.ChannelID, responseContent, m.Reference())
			utils.CheckError(err, "Error while responding to `/token`")
		}
	},
}

func InitBot() {
	var err error

	// Loading .env
	err = godotenv.Load()
	utils.CheckError(err, "Error while loading the .env")

	// Initialize session
	s, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	utils.CheckError(err, "Error while creating the bot")

	// Handling commands
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Handling token detection in chat
	s.AddHandler(botHandlers["detectToken"])
}

func Run() {
	// Open Session and close at return
	err := s.Open()
	utils.CheckError(err, "Error while running the bot")
	defer func(s *discordgo.Session) {
		err := s.Close()
		utils.CheckError(err, "Error while closing the bot")
	}(s)

	// Loading commands
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		utils.CheckError(err, "Cannot create command")
		registeredCommands[i] = cmd
	}

	// Run the bot
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
