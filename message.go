package discordclient

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Message defines discord message helpers
type Message interface {
	Channel() string
	UserName() string
	UserID() string
	UserAvatar() string
	Message() string
	RawMessage() string
	MessageID() string
	Type() MessageType
	Timestamp() (time.Time, error)
	ResolveGuildID() (string, error)
	ResolveMessageChannel() (*discordgo.Channel, error)
	IsMentionTrigger(string) (bool, string)
	IsBot() bool
}

// DiscordMessage holds received message information
type DiscordMessage struct {
	Discord          *DiscordClient
	DiscordgoMessage *discordgo.Message
	MessageType      MessageType
	Nick             *string
	Content          *string
	messageChannel   *discordgo.Channel
}

// Channel returns the message channel id
func (m *DiscordMessage) Channel() string {
	return m.DiscordgoMessage.ChannelID
}

// UserName returns the message username
func (m *DiscordMessage) UserName() string {
	me := m.DiscordgoMessage
	if me.Author == nil {
		return ""
	}

	if m.Nick == nil {
		n := m.Discord.NicknameForID(me.Author.ID, me.Author.Username, me.ChannelID)
		m.Nick = &n
	}
	return *m.Nick
}

// UserID returns the message userID
func (m *DiscordMessage) UserID() string {
	if m.DiscordgoMessage.Author == nil {
		return ""
	}

	return m.DiscordgoMessage.Author.ID
}

// UserAvatar returns the url to the message senders avatar
func (m *DiscordMessage) UserAvatar() string {
	if m.DiscordgoMessage.Author == nil {
		return ""
	}

	return discordgo.EndpointUserAvatar(m.DiscordgoMessage.Author.ID, m.DiscordgoMessage.Author.Avatar)
}

// Message returns the message in human readable text
func (m *DiscordMessage) Message() string {
	if m.Content == nil {
		c := m.DiscordgoMessage.ContentWithMentionsReplaced()
		c = m.Discord.replaceRoleNames(m.DiscordgoMessage, c)
		c = m.Discord.replaceChannelNames(m.DiscordgoMessage, c)

		m.Content = &c
	}
	return *m.Content
}

// RawMessage gets the raw message
func (m *DiscordMessage) RawMessage() string {
	return m.DiscordgoMessage.Content
}

// MessageID gets the ID of the message
func (m *DiscordMessage) MessageID() string {
	return m.DiscordgoMessage.ID
}

// Type gets the type of the message
func (m *DiscordMessage) Type() MessageType {
	return m.MessageType
}

// Timestamp gets the timestamp of the message
func (m *DiscordMessage) Timestamp() (time.Time, error) {
	return m.DiscordgoMessage.Timestamp.Parse()
}

// ResolveGuildID resolves the GuildID from the message channel
func (m *DiscordMessage) ResolveGuildID() (string, error) {
	channel, err := m.ResolveMessageChannel()

	if err != nil {
		return "", err
	}

	return channel.GuildID, nil
}

// ResolveMessageChannel resolves the Channel from the message
func (m *DiscordMessage) ResolveMessageChannel() (*discordgo.Channel, error) {
	if m.messageChannel == nil {
		channel, err := m.Discord.Channel(m.Channel())

		if err != nil {
			return nil, err
		}

		m.messageChannel = channel
	}

	return m.messageChannel, nil
}

// IsMentionPrefix returns true when the client user is the first part of the message
func (m *DiscordMessage) IsMentionPrefix() (string, bool) {
	parts := strings.Fields(m.RawMessage())

	if parts[0] == fmt.Sprintf("<@%s>", m.Discord.UserID()) || parts[0] == fmt.Sprintf("<@!%s>", m.Discord.UserID()) {
		return parts[0], true
	}

	return "", false
}

// IsMentionTrigger returns true when the client user is the first part of the message with a trigger word
func (m *DiscordMessage) IsMentionTrigger(trigger string) (bool, string) {
	parts := strings.Fields(m.RawMessage())

	if len(parts) < 2 {
		return false, ""
	}

	if _, isMentioned := m.IsMentionPrefix(); !isMentioned {
		return false, ""
	}

	mentionTrigger := fmt.Sprintf("%s %s", parts[0], parts[1])

	return parts[1] == trigger, mentionTrigger
}

// IsBot returns true if the message author is a bot. This will always be false if DiscordClient.AllowBots is false.
func (m *DiscordMessage) IsBot() bool {
	return m.DiscordgoMessage.Author != nil && m.DiscordgoMessage.Author.Bot
}
