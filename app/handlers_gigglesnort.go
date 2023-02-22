package app

import (
	"log"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/options"
)

func mentionGigglesnort(a *App, e *eventData) (cascadeAction, error) {
	if options.Gigglesnort != nil {
		response, ok := options.Gigglesnort[e.matchWord]
		if ok {
			e.gigglesnort = true
			if !e.rateLimited {
				_, err := e.s.ChannelMessageSendReply(e.m.ChannelID, response, e.m.Reference())
				if err != nil {
					return keepCascading, err
				}

				return stopCascading, nil
			} else {
				return keepCascading, nil
			}
		}
	}

	return keepCascading, nil
}

func mentionGigglesnortFallback(a *App, e *eventData) (cascadeAction, error) {
	if e.gigglesnort {
		if options.FallbackReaction != "" {
			err := e.s.MessageReactionAdd(e.m.ChannelID, e.m.ID, options.FallbackReaction)
			if err != nil {
				log.Println("mentionGigglesnortFallback :", err)
			}
		}

		if options.GigglesnortFallbackReaction != "" {
			err := e.s.MessageReactionAdd(e.m.ChannelID, e.m.ID, options.GigglesnortFallbackReaction)
			if err != nil {
				log.Println("mentionGigglesnortFallback :", err)
			}
		}

		return stopCascading, nil
	}

	return keepCascading, nil
}
