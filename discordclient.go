package discordclient

// NewDiscordClient returns a new instance of DiscordClient
func NewDiscordClient(discordToken string, ownerClientID string, discordClientID string) *DiscordClient {
	return &DiscordClient{
		args:        []interface{}{("Bot " + discordToken)},
		OwnerUserID: ownerClientID,
		ClientID:    discordClientID,
	}
}
