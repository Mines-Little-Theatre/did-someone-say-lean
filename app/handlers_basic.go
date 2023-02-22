package app

import (
	"strings"
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
	fallbackReaction, err := a.Store.GetFallbackReaction()
	if err != nil {
		return stopCascading, err
	}

	if fallbackReaction != "" {
		err = e.s.MessageReactionAdd(e.m.ChannelID, e.m.ID, fallbackReaction)
		if err != nil {
			return stopCascading, err
		}
	}

	return stopCascading, nil
}
