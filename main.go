package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gamefiend/apocalyptica/pkg/moves"
	"log"
	"os"
	"regexp"
	"strconv"
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
	}
	for i, name := range reg.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	fmt.Printf("bonus is %+v\n", result["bonus"])
	bns, e := strconv.Atoi(result["bonus"])
	if e != nil {
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

func init() {
	flag.Parse()
}

func main(){
	if err := run(os.Args); err !~= nil {
		fmt.Fprint(os.Stderr, "#{err}\n")
		os.Exit(1)
	}
}
