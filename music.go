package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	// Youtube video meta source url
	URL_META = "http://www.youtube.com/get_video_info?&video_id="
)

type Video struct {
	Id, Title, Author, Keywords, Thumbnail_url string
	Avg_rating                                 float32
	View_count, Length_seconds                 int
	Formats                                    []Format
	Filename                                   string
}

type Format struct {
	Itag                     int
	Video_type, Quality, Url string
}

func handleAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	pattern := regexp.MustCompile(`^!add (https?:\/\/)?(www\.)?(youtube\.com\/watch\?v=|youtu\.be\/)?(\S{3,})$`)

	matched := pattern.FindStringSubmatch(m.Content)
	if matched != nil {
		videoId := matched[len(matched)-1]
		fmt.Printf("processing %s\n", videoId)
		video, err := getVideo(videoId)

		if err != nil {
			fmt.Println("error during processing video,", err)
			return
		}

		if video.Length_seconds == 0 {
			return
		}

		if video.Length_seconds > 60*8 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s %d:%02d - предупреждение!!!",
				video.Title, video.Length_seconds/60, video.Length_seconds%60))

			// notify Tosori
			s.ChannelMessageSend("529035794145607681", fmt.Sprintf(
				"Продолжительность трека %s %d:%02d by @%s",
				video.Title, video.Length_seconds/60, video.Length_seconds%60, m.Author.Username))

		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"%s %d:%02d - продолжай в том же духе...",
				video.Title, video.Length_seconds/60, video.Length_seconds%60))
		}
	}
}

func getVideo(video_id string) (Video, error) {
	// fetch video meta from youtube
	query_string, err := fetchMeta(video_id)
	if err != nil {
		return Video{}, err
	}

	meta, err := parseMeta(video_id, query_string)

	if err != nil {
		return Video{}, err
	}

	return *meta, nil
}

func fetchMeta(video_id string) (string, error) {
	resp, err := http.Get(URL_META + video_id)

	// fetch the meta information from http
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	query_string, _ := ioutil.ReadAll(resp.Body)

	return string(query_string), nil
}

func parseMeta(video_id, query_string string) (*Video, error) {
	// parse the query string
	u, _ := url.Parse("?" + query_string)

	// parse url params
	query := u.Query()

	// no such video
	//if query.Get("errorcode") != "" || query.Get("status") == "fail" {
	//	fmt.Println(query)
	//	return nil, errors.New(query.Get("reason"))
	//}

	// collate the necessary params
	video := &Video{
		Id:            video_id,
		Title:         query.Get("title"),
		Author:        query.Get("author"),
		Keywords:      query.Get("keywords"),
		Thumbnail_url: query.Get("thumbnail_url"),
	}

	v, _ := strconv.Atoi(query.Get("view_count"))
	video.View_count = v

	r, _ := strconv.ParseFloat(query.Get("avg_rating"), 32)
	video.Avg_rating = float32(r)

	l, _ := strconv.Atoi(query.Get("length_seconds"))
	video.Length_seconds = l

	// further decode the format data
	format_params := strings.Split(query.Get("url_encoded_fmt_stream_map"), ",")

	// every video has multiple format choices. collate the list.
	for _, f := range format_params {
		furl, _ := url.Parse("?" + f)
		fquery := furl.Query()

		itag, _ := strconv.Atoi(fquery.Get("itag"))

		video.Formats = append(video.Formats, Format{
			Itag:       itag,
			Video_type: fquery.Get("type"),
			Quality:    fquery.Get("quality"),
			Url:        fquery.Get("url") + "&signature=" + fquery.Get("sig"),
		})
	}

	return video, nil
}
