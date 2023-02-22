package app

import (
	"strings"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/options"
)

func mentionLean(a *App, e *eventData) (cascadeAction, error) {
	if !e.rateLimited {
		leanWord := strings.ReplaceAll(e.matchWord, "lean", "**LEAN**")
		// adjacent LEANs add too many formatting characters, get rid of the extra ones
		leanWord = strings.ReplaceAll(leanWord, "****", "")

		_, err := e.s.ChannelMessageSendReply(e.m.ChannelID, "**I LOVE** "+leanWord, e.m.Reference())
		if err != nil {
			return keepCascading, err
		} else {
			return stopCascading, nil
		}
	}

	return keepCascading, nil
}

func mentionLeanFallback(a *App, e *eventData) (cascadeAction, error) {
	if options.FallbackReaction != "" {
		err := e.s.MessageReactionAdd(e.m.ChannelID, e.m.ID, options.FallbackReaction)
		if err != nil {
			return stopCascading, err
		}
	}

	return stopCascading, nil
}
