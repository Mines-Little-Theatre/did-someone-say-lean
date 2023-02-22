package app

import (
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/persist"
	"github.com/bwmarrin/discordgo"
)

// try to do stuff that doesn't involve the database first
var handlerCascade = [...]handler{
	// preprocessing
	ignoreBots,
	findLeanWord,
	processIgnores,
	pollRateLimits,
	// gigglesnort
	mentionGigglesnort,
	mentionGigglesnortFallback,
	// basic
	mentionLean,
	mentionLeanFallback,
}

// App handles incoming messages
type App struct {
	Store persist.Store
}

func (a *App) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	e := eventData{s: s, m: m}
	for _, handler := range handlerCascade {
		action, err := handler(a, &e)
		if err != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			shortFuncName := funcName[strings.LastIndexByte(funcName, '.')+1:]
			log.Println(shortFuncName, ":", err)
		}
		if action == stopCascading {
			break
		}
	}
}

// errors will be printed, but the cascade can still continue
type handler func(a *App, e *eventData) (cascadeAction, error)

// eventData is passed to each Handler, including any modifications from earlier Handlers.
type eventData struct {
	// The session
	s *discordgo.Session
	// The MessageCreate event
	m *discordgo.MessageCreate
	// The first word containing "lean" in the message
	matchWord string
	// Whether rate limits apply to the event
	rateLimited bool
	// Whether the matchWord is a gigglesnort word (used by mentionGigglesnortFallback)
	gigglesnort bool
}

type cascadeAction uint8

const (
	keepCascading cascadeAction = iota
	stopCascading
)
