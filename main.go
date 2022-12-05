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

func readEnvRequired(key string) string {
	result, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("please set the %s environment variable", key)
	}
	return result
}

type app struct {
	Token            string
	FallbackReaction string // optional
}

func main() {
	a := app{
		Token:            readEnvRequired("LEAN_TOKEN"),
		FallbackReaction: os.Getenv("LEAN_FALLBACK_REACTION"),
	}

	a.run()
}

func (a *app) run() {
	bot, err := discordgo.New(a.Token)
	if err != nil {
		log.Fatalln("failed to create bot:", err)
	}

	bot.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentMessageContent
	bot.AddHandler(a.handleMessageCreate)

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

// string is converted to lower case before being matched
var leanRegexp = regexp.MustCompile(`(?:^|[^a-z])([a-z]*lean[a-z]*)(?:[^a-z]|$)`)

func (a *app) handleMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
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
		// adjacent LEANS add too many formatting characters, get rid of the extraneous ones
		leanWord = strings.ReplaceAll(leanWord, "****", "")

		_, err := s.ChannelMessageSendReply(e.ChannelID,
			fmt.Sprintf("**I LOVE** %s", leanWord),
			e.Reference())
		if err != nil {
			log.Println("failed to send reply:", err)

			if a.FallbackReaction != "" {
				err = s.MessageReactionAdd(e.ChannelID, e.ID, a.FallbackReaction)
				if err != nil {
					log.Println("failed to react (fallback):", err)
				}
			}
		}
	}
}
