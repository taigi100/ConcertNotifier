package providers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"../utils"
	"github.com/tidwall/gjson"
)

//SongKickConnector type represents a simple connector to SongKick
type songKickConnector struct {
	apikey      string
	hostURL     string
	artistsURL  string
	concertsURL string
}

//SongKickConfig configuration object required to provide events for a specific user
type SongKickConfig struct {
	Username string `json:"username"`
}

type songKickArtist struct {
	id string
}

//NewSongKickConnector - constructor for a songkick provider
func NewSongKickConnector() (songKickConnector, error) {
	connector := songKickConnector{}
	connector.apikey = "3v8qBPX0ePTaOmbm"
	connector.hostURL = "http://api.songkick.com"
	connector.artistsURL = "/api/3.0/users/{username}/artists/tracked.json?apikey=" + connector.apikey
	connector.concertsURL = "/api/3.0/artists/{artist_id}/calendar.json?apikey=" + connector.apikey
	return connector, nil
}

func (provider songKickConnector) GetAllEvents(cfg interface{}) ([]Event, error) {
	config, ok := cfg.(SongKickConfig)
	if ok == false {
		return nil, fmt.Errorf("unexpected type %T config object for SongKick Provider. Please use a SongKickConfig object", config)
	}

	artists, err := provider.getAllFollowedArtists(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to get artists user is following")
	}

	var retVal []Event
	for _, artist := range artists {
		artistEvents, err := provider.getAllArtistsConcerts(artist)
		if err != nil {
			continue
		}
		retVal = append(retVal, artistEvents...)
	}
	return retVal, err
}

func (provider songKickConnector) getAllFollowedArtists(config SongKickConfig) ([]songKickArtist, error) {
	url := provider.hostURL + strings.Replace(provider.artistsURL, "{username}", config.Username, -1)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Couldn't fetch artists user's following: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read response body")
	}
	bodyString := string(body)
	status := gjson.Get(bodyString, "resultsPage.status").String()

	if status == "error" {
		return nil, fmt.Errorf("Couldn't fetch artists user's following: %v", gjson.Get(bodyString, "resultsPage.error.message"))
	}

	artistsString := gjson.Get(bodyString, "resultsPage.results.artist.#.id").Array()
	var artists []songKickArtist

	for _, val := range artistsString {
		artists = append(artists, songKickArtist{id: val.String()})
	}

	return artists, nil
}

func (provider songKickConnector) getAllArtistsConcerts(artist songKickArtist) ([]Event, error) {
	url := provider.hostURL + strings.Replace(provider.concertsURL, "{artist_id}", artist.id, -1)
	//TODO: See if you can extract the following lines in a function
	//duplicate with the getAllArtistsFollowed function
	res, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't fetch artists concerts: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read response body")
	}
	bodyString := string(body)
	status := gjson.Get(bodyString, "resultsPage.status").String()
	fmt.Println(status)
	if status == "error" {
		return nil, fmt.Errorf("Couldn't fetch artists concerts: %v", gjson.Get(bodyString, "resultsPage.error.message"))
	}

	concertsJSON := gjson.Get(bodyString, "resultsPage.results.event")
	var events []Event
	concertsJSON.ForEach(func(key, value gjson.Result) bool {
		event := Event{}
		event.Title = gjson.Get(value.String(), "displayName").String()
		event.Link = gjson.Get(value.String(), "uri").String()
		if gjson.Get(value.String(), "Location").Exists() {
			event.Location = utils.Location{}
			event.Location.Lat = gjson.Get(value.String(), "Location.lat").Float()
			event.Location.Long = gjson.Get(value.String(), "Location.lng").Float()
		}
		if gjson.Get(value.String(), "start.date").Exists() {
			event.StartDate, err = time.Parse("2006-01-02", gjson.Get(value.String(), "start.date").String())
			if err != nil { //TODO: make outer function return some sort of error.
				return false
			}
			end := gjson.Get(value.String(), "end.date")
			if end.Exists() {
				event.EndDate, _ = time.Parse("2006-01-02", end.String())
			} else {
				event.EndDate = event.StartDate
			}
		}
		events = append(events, event)
		return true
	})
	return events, nil
}
