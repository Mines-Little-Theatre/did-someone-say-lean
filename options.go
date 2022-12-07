package main

import "os"

// optional strings configured via environment variables
var (
	fallbackReaction string
)

func readOptions() {
	fallbackReaction = os.Getenv("LEAN_FALLBACK_REACTION")
}
