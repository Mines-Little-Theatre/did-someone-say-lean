package persist

type Store interface {
	// CheckIgnore checks whether a message sent by the user in the channel should be ignored
	CheckIgnore(userId, channelId string) (bool, error)
	// CheckBypassRateLimit checks whether a message sent by the user in the channel should bypass rate limiting
	CheckBypassRateLimit(userId, channelId string) (bool, error)

	// GetFallbackReaction returns the emoji to react with in lieu of replying, or the empty string
	GetFallbackReaction() (string, error)
	// GetGiggleSnortFallbackReaction returns the emoji to react with in lieu of replying in case of gigglesnort, or the empty string
	GetGigglesnortFallbackReaction() (string, error)
	// GetGigglesnortMessage returns the message associated with the given portmanteau, or the empty string
	GetGigglesnortMessage(word string) (string, error)

	// PollRateLimit checks whether enough time has passed since the last rate-limited event for either ID, and if so, registers a new rate-limited event as a side effect.
	PollRateLimit(userId, channelId string) (bool, error)

	// Close should be called before the application exits
	Close() error
}
