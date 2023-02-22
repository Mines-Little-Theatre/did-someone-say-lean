package app

import (
	"regexp"
	"strings"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/options"
)

func ignoreBots(a *App, e *eventData) (cascadeAction, error) {
	if e.m.Author.Bot {
		return stopCascading, nil
	} else {
		return keepCascading, nil
	}
}

// string is converted to lower case before being matched
var leanRegexp = regexp.MustCompile(`(?:^|[^a-z])([a-z]*lean[a-z]*)(?:[^a-z]|$)`)

func findLeanWord(a *App, e *eventData) (cascadeAction, error) {
	if match := leanRegexp.FindStringSubmatch(strings.ToLower(e.m.Content)); match != nil {
		e.matchWord = match[1]
		return keepCascading, nil
	} else {
		return stopCascading, nil
	}
}

func processIgnores(a *App, e *eventData) (cascadeAction, error) {
	for _, u := range options.IgnoreUsers {
		if e.m.Author.ID == u {
			return stopCascading, nil
		}
	}

	return keepCascading, nil
}

func pollRateLimits(a *App, e *eventData) (cascadeAction, error) {
	bypass, err := a.State.CheckBypassRateLimit(e.m.Author.ID, e.m.ChannelID)
	if err != nil {
		return stopCascading, err
	}

	if bypass {
		return keepCascading, nil
	} else {
		e.rateLimited, err = a.State.PollRateLimit(e.m.Author.ID, e.m.ChannelID)
		if err != nil {
			return stopCascading, err
		}

		return keepCascading, nil
	}
}
