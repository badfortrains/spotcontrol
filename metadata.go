package spotcontrol

import (
	"encoding/json"
	"fmt"
	Spotify "github.com/badfortrains/spotcontrol/proto"
	"github.com/golang/protobuf/proto"
	"net/url"
)

type Artist struct {
	Image string
	Name  string
	Uri   string
}

type Album struct {
	Artists []Artist
	Image   string
	Name    string
	Uri     string
}

type Track struct {
	Album      Album
	Artists    []Artist
	Image      string
	Name       string
	Uri        string
	Duration   int
	Popularity int
}

type TopHit struct {
	Image string
	Name  string
	Uri   string
	Log   struct {
		Origin string
		TopHit string `json:"top_hit"`
	}
	Artists []Artist
	Album   Album
}

type SearchResult struct {
	Artists struct {
		Hits  []Artist
		Total int
	}
	Albums struct {
		Hits  []Album
		Total int
	}
	Tracks struct {
		Hits  []Track
		Total int
	}
}

func (c *SpircController) Search(search string) (*SearchResult, error) {
	url := "hm://searchview/km/v2/search/" + url.QueryEscape(search) + "?limit=12&tracks-limit=100&catalogue=&country=US&locale=en&platform=zelda&username="
	done := make(chan interface{})

	go c.session.mercurySendRequest(mercuryRequest{
		method:  "GET",
		uri:     url,
		payload: [][]byte{},
	}, func(res mercuryResponse) {
		result := &SearchResult{}
		err := json.Unmarshal(res.combinePayload(), result)
		if err != nil {
			done <- err
		} else {
			done <- result
		}
	})

	result := <-done
	v, ok := result.(*SearchResult)
	if ok {
		return v, nil
	} else {
		return nil, result.(error)
	}
}

type SuggestResult struct {
	Sections []struct {
		RawItems json.RawMessage `json:"items"`
		Typ      string          `json:"type"`
	} `json:"sections"`
	Albums  []Artist
	Artists []Album
	Tracks  []Track
	TopHits []TopHit
	Error   error
}

func parseSuggest(body []byte) (*SuggestResult, error) {
	result := &SuggestResult{}
	err := json.Unmarshal(body, result)
	if err != nil {
		fmt.Println("err", err)
	}

	for _, s := range result.Sections {
		switch s.Typ {
		case "top-results":
			err = json.Unmarshal(s.RawItems, &result.TopHits)
		case "album-results":
			err = json.Unmarshal(s.RawItems, &result.Albums)
		case "artist-results":
			err = json.Unmarshal(s.RawItems, &result.Artists)
		case "track-results":
			err = json.Unmarshal(s.RawItems, &result.Tracks)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (res *mercuryResponse) combinePayload() []byte {
	body := make([]byte, 0)
	for _, p := range res.payload {
		body = append(body, p...)
	}
	return body
}

func (c *SpircController) Suggest(search string) (*SuggestResult, error) {
	url := "hm://searchview/km/v3/suggest/" + url.QueryEscape(search) + "?limit=3&intent=2516516747764520149&sequence=0&catalogue=&country=&locale=&platform=zelda&username="
	done := make(chan interface{})

	go c.session.mercurySendRequest(mercuryRequest{
		method:  "GET",
		uri:     url,
		payload: [][]byte{},
	}, func(res mercuryResponse) {
		result, err := parseSuggest(res.combinePayload())
		if err != nil {
			done <- err
		} else {
			done <- result
		}
	})

	result := <-done
	v, ok := result.(*SuggestResult)
	if ok {
		return v, nil
	} else {
		return nil, result.(error)
	}
}

func (c *SpircController) GetTrack(id string) {
	url := "hm://metadata/3/track/" + id
	c.session.mercurySendRequest(mercuryRequest{
		method:  "GET",
		uri:     url,
		payload: [][]byte{},
	}, func(res mercuryResponse) {
		track := &Spotify.Track{}
		err := proto.Unmarshal(res.payload[0], track)

		if err != nil {
			fmt.Println("error unmarshaling track")
		}

		fmt.Println("track", *track.Name)
	})

}
