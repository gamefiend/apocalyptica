package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func run(args []string) error {
	var wg sync.WaitGroup
	//Grab Token from the Environment
	Token := os.Getenv("DISCORD_TOKEN")
	if len(Token) == 0 {
		fmt.Println("environment variable DISCORD_TOKEN not set")
		os.Exit(1)
	}
	//launch a small web server

	go func() {
		s, err := os.Getwd()
		if err != nil {
			log.Println("Can't access working directory", err.Error())
			return err
		}
		httpStart(s)
	}()
	// new discord session with the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return err
	}
	// Register the messageCreate func for a callback to MessageCreate events
	dg.AddHandler(messageCreate)
	dg.AddHandler(onReady)
	// listen on that websocket connection!
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return err
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
	fmt.Println("received ", signals)
	fmt.Println("Stopping....")
	wg.Done()
}

func httpStart(s string) {
	httpd := fmt.Sprintf("%s/static", s)
	fmt.Println(httpd)
	fs := http.FileServer(http.Dir(httpd))
	http.Handle("/", fs)
	log.Printf("listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
