package app

import (
	"regexp"
	"strings"
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
	ignore, err := a.Store.CheckIgnore(e.m.Author.ID, e.m.ChannelID)
	if err != nil {
		return stopCascading, err
	}

	if ignore {
		return stopCascading, nil
	} else {
		return keepCascading, nil
	}
}

func pollRateLimits(a *App, e *eventData) (cascadeAction, error) {
	bypass, err := a.Store.CheckBypassRateLimit(e.m.Author.ID, e.m.ChannelID)
	if err != nil {
		return stopCascading, err
	}

	if bypass {
		return keepCascading, nil
	} else {
		e.rateLimited, err = a.Store.PollRateLimit(e.m.Author.ID, e.m.ChannelID)
		if err != nil {
			return stopCascading, err
		}

		return keepCascading, nil
	}
}
