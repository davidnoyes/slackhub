package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SlackMessage struct {
	Text       string `json:"text"`
	Channel    string `json:"channel,omitempty"`
	Username   string `json:"username,omitempty"`
	Icon_url   string `json:"icon_url,omitempty"`
	Icon_emoji string `json:"icon_emoji,omitempty"`
}

type SonarrEpisodes struct {
	Id             int    `json:"Id"`
	EpisodeNumber  int    `json:"EpisodeNumber"`
	SeasonNumber   int    `json:"SeasonNumber"`
	Title          string `json:"Title"`
	AirDate        string `json:"AitDate"`
	AirDateUtc     string `json:"AirDateUtc"`
	Quality        string `json:"Quality"`
	QualityVersion int    `json:"QualityVersion"`
	ReleaseGroup   string `json:"ReleaseGroup"`
	SceneName      string
}

type SonarrSeries struct {
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Path   string `json:"Path"`
	TvdbId int    `json:"TvdbId"`
}

type SonarrData struct {
	EventType string           `json:"EventType"`
	Series    SonarrSeries     `json:"Series"`
	Episodes  []SonarrEpisodes `json:"Episodes"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	DisplayAllProperties()
	router := gin.Default()

	router.POST("/couchpotato", couchpotato)
	router.POST("/sonarr", sonarr)

	router.Run(fmt.Sprintf("%s:%d", ListenAddress(), ListenPort()))
}

func sonarr(c *gin.Context) {
	var sonarrData SonarrData

	body := c.Request.Body
	x, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(x, &sonarrData)

	if err == nil {
		if sonarrData.EventType != "Rename" {
			for _, episode := range sonarrData.Episodes {
				var text = fmt.Sprintf("%s: %s S%d E%d (%s)",
					sonarrData.EventType,
					sonarrData.Series.Title,
					episode.SeasonNumber,
					episode.EpisodeNumber,
					episode.Quality)

				message := SlackMessage{
					Text:       text,
					Username:   "Sonarr",
					Icon_emoji: ":sonarr:",
				}
				send(message)
			}
		} else {
			var text = fmt.Sprintf("%s: %s (%s)",
				sonarrData.EventType,
				sonarrData.Series.Title,
				sonarrData.Series.Path)

			message := SlackMessage{
				Text:       text,
				Username:   "Sonarr",
				Icon_emoji: ":sonarr:",
			}
			send(message)
		}
	}
}

func couchpotato(c *gin.Context) {
	text := c.DefaultPostForm("message", ":worried:")

	message := SlackMessage{
		Text:       text,
		Username:   "CouchPotato",
		Icon_emoji: ":couchpotato:",
	}
	send(message)
}

func send(message SlackMessage) {
	log.Printf("%#v", message)

	messageBuffer := new(bytes.Buffer)
	json.NewEncoder(messageBuffer).Encode(message)
	req, err := http.NewRequest("POST", WebhookUrl(), messageBuffer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Response Status: %s, Body: %s", resp.Status, string(body))
	//log.Println("response Headers:", resp.Header)
}
