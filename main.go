package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/app"
	"github.com/Mines-Little-Theatre/did-someone-say-lean/persist"
	"github.com/Mines-Little-Theatre/did-someone-say-lean/utils"
	"github.com/bwmarrin/discordgo"
)

func main() {
	bot, err := discordgo.New(utils.ReadEnvRequired("LEAN_TOKEN"))
	if err != nil {
		log.Fatalln("failed to create bot:", err)
	}

	store, err := persist.Connect()
	if err != nil {
		log.Fatalln("failed to connect store:", err)
	}

	app := app.App{Store: store}

	bot.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentMessageContent
	bot.AddHandler(app.HandleMessageCreate)

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
