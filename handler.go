package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gamefiend/apocalyptica/pkg/moves"
	"log"
	"strings"
)

// Addhandler sends a message to this function any time a message on a channel this bot is listening to is created.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore our own messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	//here is where look for commands.
	// going to move this to another package, but just hardwiring this for an MVP right now.
	switch {
	case m.Content == "!!help":
		help := []string{
			"**Apocalyptica**. Apocalypse World 2e bot. !!help for instructions, !moves for moves.",
			"To use, type !<move> <bonus>",
			"**examples:**",
			"* !goaggro -1",
			"* !goaggro 1",
			"*!moves* shows which moves I currently support.",
		}
		s.ChannelMessageSend(m.ChannelID, strings.Join(help, "\n"))

	case m.Content == "!moves":
		help := "I currently support:"
		s.ChannelMessageSend(m.ChannelID, help)
		s.ChannelMessageSend(m.ChannelID, strings.Join(ApocList, "\n"))
	default:
		moveMsg := moves.FindMove(m.Content, Apoc)
		if len(moveMsg) > 0 {
			bonus := getBonus(m.Content)
			result := moveMsg.Roll(bonus)
			d := moveMsg.Display(result, bonus)
			s.ChannelMessageSend(m.ChannelID, d)
			log.Println(m.ChannelID, m.Author.ID, result, bonus)
		}
	}
}

