package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func readEnvRequired(key string) string {
	result, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("please set the %s environment variable", key)
	}
	return result
}

func main() {
	bot, err := discordgo.New(readEnvRequired("LEAN_TOKEN"))
	if err != nil {
		log.Fatalln("failed to create bot:", err)
	}

	readOptions()

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

func handleMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	if !e.Author.Bot { // don't respond to yourself or other bots

		for _, handler := range handlerCascade {
			if handler(s, e) {
				break
			}
		}
	}
}
