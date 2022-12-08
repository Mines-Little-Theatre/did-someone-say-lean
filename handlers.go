package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler tries to handle a message, returning true if no more handlers should run
type Handler func(s *discordgo.Session, e *discordgo.MessageCreate, matchWord string) bool

// handlerCascade is the order in which to run handlers until one returns true
var handlerCascade = [...]Handler{
	processIgnores,
	mentionGigglesnort,
	mentionLean,
}

func mentionLean(s *discordgo.Session, e *discordgo.MessageCreate, matchWord string) bool {
	// injection attack considerations:
	// the regexp only matches on word characters,
	// so a malicious user can't:
	// - make the bot mention anyone
	// - mess with the markdown
	// - make the bot post a link
	// so this is probably safe.

	leanWord := strings.ReplaceAll(matchWord, "lean", "**LEAN**")
	// adjacent LEANs add too many formatting characters, get rid of the extraneous ones
	leanWord = strings.ReplaceAll(leanWord, "****", "")

	_, err := s.ChannelMessageSendReply(e.ChannelID,
		fmt.Sprintf("**I LOVE** %s", leanWord),
		e.Reference())
	if err != nil {
		log.Println("failed to send reply:", err)

		if fallbackReaction != "" {
			err = s.MessageReactionAdd(e.ChannelID, e.ID, fallbackReaction)
			if err != nil {
				log.Println("failed to react (fallback):", err)
			}
		}
	}

	return true
}

func processIgnores(s *discordgo.Session, e *discordgo.MessageCreate, matchWord string) bool {
	for _, u := range ignoreUsers {
		if e.Author.ID == u {
			return true
		}
	}

	return false
}

func mentionGigglesnort(s *discordgo.Session, e *discordgo.MessageCreate, matchWord string) bool {
	if gigglesnort != nil {
		response, ok := gigglesnort[matchWord]
		if ok {
			_, err := s.ChannelMessageSendReply(e.ChannelID, response, e.Reference())
			if err != nil {
				log.Println("failed to send gigglesnort reply:", err)

				if fallbackReaction != "" {
					err = s.MessageReactionAdd(e.ChannelID, e.ID, fallbackReaction)
					if err != nil {
						log.Println("failed to react (fallback):", err)
					}
				}
			}

			return true
		}
	}

	return false
}
