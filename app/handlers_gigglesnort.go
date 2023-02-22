package app

import (
	"log"
)

func mentionGigglesnort(a *App, e *eventData) (cascadeAction, error) {
	message, err := a.Store.GetGigglesnortMessage(e.matchWord)
	if err != nil {
		return stopCascading, err
	}

	if message != "" {
		e.gigglesnort = true
		if !e.rateLimited {
			_, err := e.s.ChannelMessageSendReply(e.m.ChannelID, message, e.m.Reference())
			if err != nil {
				return keepCascading, err
			}

			return stopCascading, nil
		} else {
			return keepCascading, nil
		}
	}

	return keepCascading, nil
}

func mentionGigglesnortFallback(a *App, e *eventData) (cascadeAction, error) {
	if e.gigglesnort {
		// TODO these should maybe be optimized into one query

		fallbackReaction, err := a.Store.GetFallbackReaction()
		if err != nil {
			return stopCascading, err
		}

		gigglesnortFallbackReaction, err := a.Store.GetGigglesnortFallbackReaction()
		if err != nil {
			return stopCascading, err
		}

		if fallbackReaction != "" {
			err = e.s.MessageReactionAdd(e.m.ChannelID, e.m.ID, fallbackReaction)
			if err != nil {
				log.Println("mentionGigglesnortFallback :", err)
			}
		}

		if gigglesnortFallbackReaction != "" {
			err := e.s.MessageReactionAdd(e.m.ChannelID, e.m.ID, gigglesnortFallbackReaction)
			if err != nil {
				log.Println("mentionGigglesnortFallback :", err)
			}
		}

		return stopCascading, nil
	}

	return keepCascading, nil
}
