package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler tries to handle a message, returning true if no more handlers should run
type Handler func(*discordgo.Session, *discordgo.MessageCreate) bool

// handlerCascade is the order in which to run handlers until one returns true
var handlerCascade = [...]Handler{
	mentionLean,
}

// string is converted to lower case before being matched
var leanRegexp = regexp.MustCompile(`(?:^|[^a-z])([a-z]*lean[a-z]*)(?:[^a-z]|$)`)

func mentionLean(s *discordgo.Session, e *discordgo.MessageCreate) bool {
	if match := leanRegexp.FindStringSubmatch(strings.ToLower(e.Content)); match != nil {
		// injection attack considerations:
		// the regexp only matches on word characters,
		// so a malicious user can't:
		// - make the bot mention anyone
		// - mess with the markdown
		// - make the bot post a link
		// so this is probably safe.

		leanWord := strings.ReplaceAll(match[1], "lean", "**LEAN**")
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

	return false
}
