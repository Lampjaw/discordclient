package discordclient

import "github.com/bwmarrin/discordgo"

// NewDiscordClient returns a new instance of DiscordClient
func NewDiscordClient(discordToken string, ownerClientID string, discordClientID string) *DiscordClient {
	return &DiscordClient{
		args:           []interface{}{("Bot " + discordToken)},
		messageChan:    make(chan Message, 200),
		OwnerUserID:    ownerClientID,
		ClientID:       discordClientID,
		GatewayIntents: discordgo.MakeIntent(discordgo.IntentsAll),
	}
}
