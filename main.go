package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gamefiend/apocalyptica/moves"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	Token string
)

var Apoc = moves.LoadMoves()
var ApocList = MakeList(Apoc)
var Announce = make(map[string]bool)

func MakeList(mv moves.Move) []string {
	l := make([]string, 0)
	for _, c := range mv {
		l = append(l, fmt.Sprintf("**%s** (%s)", c.Full, c.Name))
	}
	return l
}

func getBonus(s string) int {
	reg := regexp.MustCompile(`(?P<bonus>-?\d+)`)
	match := reg.FindStringSubmatch(s)
	result := make(map[string]string)
	if match == nil {
		return 0
		fmt.Printf("no bonus")
	}
	for i, name := range reg.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	fmt.Printf("bonus is %+v\n", result["bonus"])
	bns, e := strconv.Atoi(result["bonus"])
	if e != nil {
		fmt.Print("problem getting bonus: %v\n", e)
		return 0
	}
	return bns
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// new discord session with the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}

	// Register the messageCreate func for a callbackto MessageCreate events
	dg.AddHandler(messageCreate)
	dg.AddHandler(onReady)
	// listen on that websocket connection!
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	//wait until we get Ctl-C/term signal
	fmt.Println("Loading Moves")
	fmt.Println("Apocalyptica is Barfing into your channel now...")
	fmt.Println("Press Ctl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc

	dg.Close()

}

// Addhandler fires this up when Apocalyptica connects to another channel.
func onReady(s *discordgo.Session, m *discordgo.Ready) {
	greetings := "**Apocalyptica**. Apocalypse World 2e bot. !!help for instructions, !moves for moves."

	for _, i := range m.Guilds {
		fmt.Printf("%+v\n", i)
		// Wait for a few seconds to finish connection, otherwise we miss info.
		time.Sleep(2 * time.Second)
		for _, c := range i.Channels {
			fmt.Printf("\t%+v\n", c)
			if c.Type == "text" {
				s.ChannelMessageSend(c.ID, greetings)
				fmt.Printf("Introducing myself in Guild %s Channel %s\n", i.Name, c.Name)
			}
		}
	}
	// s.ChannelMessageSend(m.ChannelID, greetings)
}

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
			fmt.Println(m.ChannelID, d)
		}
	}
}
