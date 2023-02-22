package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/persist"
	"github.com/bwmarrin/discordgo"
)

type EventData struct {
	// The first word containing "lean" in the message
	MatchWord string
	// Whether rate limits apply to the event
	RateLimited bool
}

// Handler processes a message and its metadata, returning true if no more handlers should run
type Handler func(state persist.State, s *discordgo.Session, e *discordgo.MessageCreate, data *EventData) bool

const (
	keepCascading = false
	stopCascading = true
)

// handlerCascade is the order in which to run handlers until one returns true
var handlerCascade = [...]Handler{
	ignoreBots,
	findLeanWord,
	processIgnores,
	pollRateLimits,
	mentionGigglesnort,
	mentionLean,
}

func ignoreBots(state persist.State, s *discordgo.Session, e *discordgo.MessageCreate, data *EventData) bool {
	// don't respond to yourself or other bots
	if e.Author.Bot {
		return stopCascading
	} else {
		return keepCascading
	}
}

// string is converted to lower case before being matched
var leanRegexp = regexp.MustCompile(`(?:^|[^a-z])([a-z]*lean[a-z]*)(?:[^a-z]|$)`)

func findLeanWord(state persist.State, s *discordgo.Session, e *discordgo.MessageCreate, data *EventData) bool {
	if match := leanRegexp.FindStringSubmatch(strings.ToLower(e.Content)); match != nil {
		data.MatchWord = match[1]
		return keepCascading
	} else {
		return stopCascading
	}
}

func pollRateLimits(state persist.State, s *discordgo.Session, e *discordgo.MessageCreate, data *EventData) bool {
	bypass, err := state.CheckBypassRateLimit(e.Author.ID, e.ChannelID)
	if err != nil {
		log.Println(err)
		return stopCascading
	}

	if bypass {
		data.RateLimited = false
		return keepCascading
	} else {
		data.RateLimited, err = state.PollRateLimit(e.Author.ID, e.ChannelID)
		if err != nil {
			log.Println(err)
			return stopCascading
		}

		return keepCascading
	}
}

func mentionLean(state persist.State, s *discordgo.Session, e *discordgo.MessageCreate, data *EventData) bool {
	if !data.RateLimited {
		leanWord := strings.ReplaceAll(data.MatchWord, "lean", "**LEAN**")
		// adjacent LEANs add too many formatting characters, get rid of the extraneous ones
		leanWord = strings.ReplaceAll(leanWord, "****", "")

		_, err := s.ChannelMessageSendReply(e.ChannelID,
			fmt.Sprintf("**I LOVE** %s", leanWord),
			e.Reference())
		if err != nil {
			log.Println("failed to send reply:", err)
			mentionLeanFallback(s, e)
		}
	} else {
		mentionLeanFallback(s, e)
	}

	return stopCascading
}

func mentionLeanFallback(s *discordgo.Session, e *discordgo.MessageCreate) {
	if fallbackReaction != "" {
		err := s.MessageReactionAdd(e.ChannelID, e.ID, fallbackReaction)
		if err != nil {
			log.Println("failed to react (fallback):", err)
		}
	}
}

func processIgnores(state persist.State, s *discordgo.Session, e *discordgo.MessageCreate, data *EventData) bool {
	for _, u := range ignoreUsers {
		if e.Author.ID == u {
			return stopCascading
		}
	}

	return keepCascading
}

func mentionGigglesnort(state persist.State, s *discordgo.Session, e *discordgo.MessageCreate, data *EventData) bool {
	if gigglesnort != nil {
		response, ok := gigglesnort[data.MatchWord]
		if ok {
			if !data.RateLimited {
				_, err := s.ChannelMessageSendReply(e.ChannelID, response, e.Reference())
				if err != nil {
					log.Println("failed to send gigglesnort reply:", err)
					mentionGigglesnortFallback(s, e)
				}
			} else {
				mentionGigglesnortFallback(s, e)
			}

			return stopCascading
		}
	}

	return keepCascading
}

func mentionGigglesnortFallback(s *discordgo.Session, e *discordgo.MessageCreate) {
	if fallbackReaction != "" {
		err := s.MessageReactionAdd(e.ChannelID, e.ID, fallbackReaction)
		if err != nil {
			log.Println("failed to react (gigglesnort fallback 1):", err)
		}
	}

	err := s.MessageReactionAdd(e.ChannelID, e.ID, "❗")
	if err != nil {
		log.Println("failed to react (gigglesnort fallback 2):", err)
	}
}
