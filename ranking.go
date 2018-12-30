package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
)

type Rankings struct {
	Admin             bool     `json:"admin"`
	BannerUrl         string   `json:"banner_url"`
	Guild             Guild    `json:"guild"`
	Page              int      `json:"page"`
	Player            Player   `json:"player"`
	Players           []Player `json:"players"`
	RoleRewards       []string `json:"role_rewards"`
	UserGuildSettings string   `json:"user_guild_settings"`
}

type Guild struct {
	Icon    string `json:"icon"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Premium bool   `json:"premium"`
}

type Player struct {
	Avatar        string `json:"avatar"`
	DetailedXp    []int  `json:"detailed_xp"`
	Discriminator string `json:"discriminator"`
	GuildId       string `json:"guild_id"`
	Id            string `json:"id"`
	Level         int    `json:"level"`
	Username      string `json:"username"`
	Xp            int64  `json:"xp"`
}

func handleRanking(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ranking" {
		response, err := http.Get("https://mee6.xyz/api/plugins/levels/leaderboard/404185377264369665")
		if err != nil {
			fmt.Println("error retrieving dat,", err)
			return
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		var rankings Rankings
		err = json.Unmarshal(contents, &rankings)

		if err != nil {
			fmt.Println("error parsing json,", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"TOP:\n1. %s (%d level with %d messages)\n2. %s (%d level with %d messages)\n3. %s (%d level with %d messages)",
			rankings.Players[0].Username, rankings.Players[0].Level, rankings.Players[0].Xp/20,
			rankings.Players[1].Username, rankings.Players[1].Level, rankings.Players[1].Xp/20,
			rankings.Players[2].Username, rankings.Players[2].Level, rankings.Players[2].Xp/20))
	}
}
