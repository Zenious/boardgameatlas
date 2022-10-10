package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/zenious/boardgameatlas/api"
)

// "go run ."" to run program
func main() {
	// Define the command line argument
	query := flag.String("query", "", "Boardgame name to search")
	clientId := flag.String("clientId", "", "Client id to be passed in for BoardGameAtlas")
	skip := flag.Uint("skip", 0, "Skips the number of results provided. It's generally used for paging results.")
	limit := flag.Uint("limit", 10, "Limits the number of results returned.")
	timeout := flag.Uint("timeout", 10, "Timeout")
	// Parse the command line
	flag.Parse()

	if isNull(*query) {
		log.Fatalln("Please use --query to set the boardgame name to search")
	}
	if isNull(*clientId) {
		log.Fatalln("Please use --clientId to set the clientId")
	}
	fmt.Printf("query=%s, clientId=%s, limit=%d, skip=%d\n", *query, *clientId, *limit, *skip)

	// Create a new instance of the Boardgame Atlas client
	bga := api.New(*clientId)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout*uint(time.Second)))
	defer cancel()

	// Make the invocation
	result, err := bga.Search(ctx, *query, *limit, *skip)
	if nil != err {
		log.Fatalf("cannot search for boardgame: %v", err)
	}

	// Colors
	boldGreen := color.New(color.Bold).Add(color.FgHiGreen).SprintFunc()

	for _, g := range result.Games {
		fmt.Printf("%s : %s\n", boldGreen("Name"), g.Name)
		fmt.Printf("%s : %s\n", boldGreen("Description"), g.Description)
		fmt.Printf("%s : %s\n\n", boldGreen("Url"), g.Url)
	}
}

func isNull(s string) bool {
	return len(strings.TrimSpace(s)) <= 0
}
