package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const SEARCH_URL = "https://api.boardgameatlas.com/api/search"

type BoardGameAtlas struct {
	// "members"
	clientId string
}

// Game
type Game struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Price         string `json:"price"`
	YearPublished uint   `json:"year_published"`
	Description   string `json:"description"`
	Url           string `json:"official_url"`
	ImageUrl      string `json:"image_url"`
	RulesUrl      string `json:"rules_url"`
}

type SearchResult struct {
	Games []Game `json:"games"`
	Count uint   `json:"count"`
}

func New(clientId string) BoardGameAtlas {
	return BoardGameAtlas{clientId}
}

func (b BoardGameAtlas) Search(ctx context.Context, query string, limit uint, skip uint) (*SearchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, SEARCH_URL, nil)
	if nil != err {
		// return an error object
		return nil, fmt.Errorf("cannot create HTTP client: %v", err)
	}

	// Get query string object
	queryObj := req.URL.Query()

	// Populate the URL with query params
	queryObj.Add("name", query)
	queryObj.Add("limit", fmt.Sprintf("%d", limit))
	queryObj.Add("skip", strconv.Itoa(int(skip)))
	queryObj.Add("client_id", b.clientId)
	// Encode the query params and add back to the base URI
	req.URL.RawQuery = queryObj.Encode()
	// fmt.Printf("URL = %s\n", req.URL.String())

	// Make the call
	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return nil, fmt.Errorf("cannot create HTTP client for invocation: %v", err)
	}

	// HTTP status code >= 400 is error
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error HTTP status: %s", resp.Status)
	}

	// Deserialise the JSON payload to struct
	var result SearchResult
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); nil != err {
		return nil, fmt.Errorf("cannot deserialize JSON payload: %v", err)
	}

	return &result, nil
}
