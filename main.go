package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func readEnv(key string) string {
	result, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("please set the %s environment variable", key)
	}
	return result
}

func main() {
	token := readEnv("LEAN_TOKEN")

	bot, err := discordgo.New(token)
	if err != nil {
		log.Fatalln("failed to create bot:", err)
	}

	bot.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentMessageContent
	bot.AddHandler(handleMessageCreate)

	err = bot.Open()
	if err != nil {
		log.Fatalln("failed to open connection:", err)
	}

	log.Println("Bot is running")

	// clean shutdown on CTRL+C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
	log.Println("Shutting down bot")
	bot.Close()
}

var leanRegexp = regexp.MustCompile(`(?:^|\W)(\w*lean\w*)(?:\W|$)`)

func handleMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot {
		return // don't respond to yourself or other bots
	}

	if match := leanRegexp.FindStringSubmatch(strings.ToLower(e.Content)); match != nil {
		// injection attack considerations:
		// the regexp only matches on word characters,
		// so a malicious user can't:
		// - make the bot mention anyone
		// - mess with the markdown
		// - make the bot post a link
		// so this is probably safe.

		leanWord := strings.ReplaceAll(match[1], "lean", "**LEAN**")
		// handle adjacent LEANs
		leanWord = strings.ReplaceAll(leanWord, "LEAN****LEAN", "LEANLEAN")

		_, err := s.ChannelMessageSendReply(e.ChannelID,
			fmt.Sprintf("**I LOVE** %s", leanWord),
			e.Reference())
		if err != nil {
			log.Println("failed to send reply:", err)
		}
	}
}
