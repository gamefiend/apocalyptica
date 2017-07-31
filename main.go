package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gamefiend/apocalyptica/moves"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	Intro = flag.Bool("intro", false, "Channel Introductions")
)

var Apoc = moves.LoadMoves("basic.json")
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

// Addhandler fires this up when Apocalyptica connects to another channel.
func onReady(s *discordgo.Session, m *discordgo.Ready) {
	greetings := "**Apocalyptica**. Apocalypse World 2e bot. !!help for instructions, !moves for moves."
	log.Printf("Invited to %s servers\n", len(m.Guilds))
	for _, i := range m.Guilds {
		log.Printf("CONNECT [%s]\n", i.Name)
		// Wait for a few seconds to finish connection, otherwise we miss info.
		time.Sleep(2 * time.Second)
		for _, c := range i.Channels {
			log.Printf("%s.%s [%s]\n", i.Name, c.Name, c.ID)
			if c.Type == "text" && *Intro == true {
				s.ChannelMessageSend(c.ID, greetings)
				fmt.Printf("%s.%s Introduction\n", i.Name, c.Name)
			}
		}
	}
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
			log.Println(m.ChannelID, m.Author.ID, moveMsg, result, bonus)
		}
	}
}

func init() {
	flag.Parse()
}

func main() {
	var wg sync.WaitGroup
	//Grab Token from the Environment
	Token := os.Getenv("DISCORD_TOKEN")
	if len(Token) == 0 {
		fmt.Println("environment variable DISCORD_TOKEN not set")
		os.Exit(1)
	}
	//launch a small web server

	go func() {
		s, e := os.Getwd()
		if e != nil {
			log.Println("Can't access working directory", e)
		}
		httpwd := fmt.Sprintf("%s/static", s)
		fmt.Println(httpwd)
		fs := http.FileServer(http.Dir(httpwd))
		http.Handle("/", fs)
		log.Printf("listening on port 8080")
		http.ListenAndServe(":8080", nil)
	}()
	// new discord session with the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
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
	fmt.Println(
		`Loading Moves...
Apocalyptica is Barfing into your channel now...
Press Ctl-C to exit.`)
	wg.Add(1)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	signals := <-sc
	dg.Close()
	fmt.Println("recieved ", signals)
	fmt.Println("Stopping....")
	wg.Done()
}
