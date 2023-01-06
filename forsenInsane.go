package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)
type WineTime struct {
	Code string `json:"code"`
	Url string `json:"url"`
}
type monkaE struct {
	Username string `json:"username"`
	Message string `json:"message"`
	Timestamp string `json:"timestamp"`
	Emotes []WineTime `json:"emotes"`
}
type dankHug struct {
	Vod_id string `json:"vod_id"`
	From string `json:"from"`
	To string `json:"to"`
}
type Sadge struct {
	Vods []dankHug `json:"vods"`
}
type PoroSad struct {
	Frogs []monkaE `json:"frogs"`
}
type FeelsGoodMan struct {
	Username string `json:"username"`
	Message string `json:"message"`
	Emotes []WineTime `json:"emotes"`
}
var pwd, _ = os.Getwd()
var masterkey = "12345"
func happE(c *fiber.Ctx) error {
	err := os.MkdirAll(fmt.Sprintf("%s/forsenBussin/", pwd), 0750)
	if err != nil {
		return c.Status(400).SendString("PepeLaugh [Unable to create folder]")
	}
	_, err = os.Stat(fmt.Sprintf("%s/forsenBussin/vods.json", pwd))
	if err != nil {
		f, err := os.OpenFile(fmt.Sprintf("%s/forsenBussin/vods.json", pwd), os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			return c.Status(400).SendString("PepeLaugh [Unable to create json file]")
		}
		vod := Sadge{
			Vods: []dankHug{},
		}
		pretty, _ := json.MarshalIndent(vod, "", "\t")
		if err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while converting struct to json]")
		}
		if _, err := f.WriteString(string(pretty)); err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while writing string to file]")
		}
		if err := f.Close(); err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while closing file]")
		}
	}
	
    if c.Method() == "GET" {
		f, err := os.ReadFile(fmt.Sprintf("%s/forsenBussin/vods.json", pwd))
		if err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while reading VODs file]")
		}
		var vod Sadge
		err = json.Unmarshal(f, &vod)
		if err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while converting file to json]")
		}
		return c.JSON(vod)
	}
   if c.Method() == "POST" {
		key := c.Get("Authorization")
		if key != masterkey {
		return c.Status(400).SendString("PepeLaugh [Not Authorized]")
		}
		var small dankHug
		if err := c.BodyParser(&small); err != nil {
			return c.Status(400).SendString("PepeLaugh [Invalid request body]")
		}
		if small.From == "" || small.To == "" || small.Vod_id == "" {
			return c.Status(400).SendString("PepeLaugh [Empty]")
		}
		f, err := os.ReadFile(fmt.Sprintf("%s/forsenBussin/vods.json", pwd))
		if err != nil {
			return c.Status(400).SendString("PepeLaugh [Unable to read VODs file]")
		}
		var big Sadge
		err = json.Unmarshal(f, &big)
		if err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while converting VODs file to json]")
		}
		big.Vods = append(big.Vods, small)
		pretty, _ := json.MarshalIndent(big, "", "\t")
		file, err := os.OpenFile(fmt.Sprintf("%s/forsenBussin/vods.json", pwd), os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong opening VODs file]")
		}
		if _, err := file.WriteString(string(pretty)); err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while appending to VODs file]")
		}
		if err := file.Close(); err != nil {
			return c.Status(400).SendString("PepeLaugh [Something went wrong while closing VODs file]")
		}
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.Status(400).SendString("PepeLaugh [Method Not Allowed]")
	}

}
func FeelsStrongMan(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" && to == "" {
		return c.Status(400).SendString("I don't care that much...") 
	} else if from == "" {
		return c.Status(400).SendString("I don't care that much...") 
	} else if to == "" {
		return c.Status(400).SendString("I don't care that much...") 
	}
	from = strings.Replace(from, " ", "+", 1)
	to = strings.Replace(to, " ", "+", 1)
	fromTrans, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...") 
	}
	toTrans, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...") 
	}
	fromGay := time.Date(fromTrans.Year(), fromTrans.Month(), fromTrans.Day(), 0, 0, 0, 0, fromTrans.Location())
	toGay := time.Date(toTrans.Year(), toTrans.Month(), toTrans.Day(), 0, 0, 0, 0, toTrans.Location())
	diff := toGay.Sub(fromGay)
	var urls []string
	for i := 0; i <= int(diff.Hours()); i += 24 {
		bufferTime := fromTrans.Add(time.Duration(i) * time.Hour)
		url := fmt.Sprintf("%s-%s-%s.json", bufferTime.Format("02"), bufferTime.Format("01"),bufferTime.Format("2006"))
		urls = append(urls, url)
	}

	logs := []monkaE{}
	for _, i := range urls {
		path := fmt.Sprintf(pwd+"/forsenCD/%s", i)
		if _, err := os.Stat(path); err != nil {
			if _, err := os.Stat(fmt.Sprintf("%s.gz", path)); err != nil {
				return c.Status(400).SendString("I don't care that much...") 
			}
			path = fmt.Sprintf(pwd+"/forsenCD/%s.gz", i)
		}
		var f []byte
		if strings.HasSuffix(path, ".gz") {
			var all bytes.Buffer
			file, err := os.Open(path)
			if err != nil {
				return c.Status(400).SendString("I don't care that much...") 
			}
			gz, err := gzip.NewReader(file)
			if err != nil {
				return c.Status(400).SendString("I don't care that much...")
			}
			defer file.Close()
			defer gz.Close()
			pre := bufio.NewScanner(gz)
			for pre.Scan() {
				line := pre.Text()
				all.WriteString(line)
			}
			f = all.Bytes()
		} else {
			f, err = os.ReadFile(path)
			if err != nil {
				return c.Status(400).SendString("I don't care that much...") 
			}
		}
		var xqcL PoroSad
		err = json.Unmarshal(f, &xqcL)
		if err != nil {
			return c.Status(400).SendString("I don't care that much...")
		}
		logs = append(logs, xqcL.Frogs...)
	}
	forsenPossessed := []monkaE{}
	for _, raw := range logs {
		t, _ := time.Parse(time.RFC3339, raw.Timestamp)
		if t.Unix() < fromTrans.Unix() || t.Unix() > toTrans.Unix() {
			continue
		}
		forsenPossessed = append(forsenPossessed, raw)
	}
	data := map[string][]FeelsGoodMan{}
	for _, v := range forsenPossessed {
		x, err := time.Parse(time.RFC3339, v.Timestamp)
		if err != nil {
			return c.Status(400).SendString("I don't care that much...")
		}
		almost := x.Unix()
		nnn := int(almost)
		unix := fmt.Sprintf("%d", nnn)
		if data[unix] == nil {
			data[unix] = []FeelsGoodMan{{Username: v.Username,  Message:  v.Message, Emotes: v.Emotes}}
		} else {
			data[unix] = append(data[unix], FeelsGoodMan{Username: v.Username,  Message:  v.Message, Emotes: v.Emotes})
		}
    }
	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	c.Context().SetBody(payload)
	c.Context().SetContentType(fiber.MIMEApplicationJSON)
	return nil
}
func scarE(c *fiber.Ctx) error {
	type NOIDONTTHINKSO struct {
		Code string `json:"code"`
		Url string `json:"url"`
	}
	type cringeE struct {
		Seventv []NOIDONTTHINKSO `json:"7tv"`
		Bttv []NOIDONTTHINKSO `json:"bttv"`
		Twitch []NOIDONTTHINKSO `json:"twitch"`
		Ffz []NOIDONTTHINKSO `json:"ffz"`
	}
	var big cringeE
	type BTTVEmote struct {
		ID        string    `json:"id"`
		Code      string    `json:"code"`
		ImageType string `json:"imageType"`
		UserID    string    `json:"userId"`
	}
	type BTTVChannelEmotes struct {
		ID            string        `json:"id"`
		Bots          []interface{} `json:"bots"`
		Avatar        string        `json:"avatar"`
		ChannelEmotes []BTTVEmote   `json:"channelEmotes"`
		SharedEmotes  []BTTVEmote   `json:"sharedEmotes"`
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.betterttv.net/3/cached/emotes/global", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var bttvGlobalEmotes []BTTVEmote
	err = json.Unmarshal(body, &bttvGlobalEmotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	for _, emote := range bttvGlobalEmotes {
		big.Bttv = append(big.Bttv, NOIDONTTHINKSO{Code:emote.Code, Url: fmt.Sprintf("https://cdn.betterttv.net/emote/%s/1x", emote.ID),})
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", "https://api.betterttv.net/3/cached/users/twitch/50465063", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	resp, err = client.Do(req)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var bttvChannelEmotes BTTVChannelEmotes
	err = json.Unmarshal(body, &bttvChannelEmotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	for _, emote := range bttvChannelEmotes.ChannelEmotes {
		big.Bttv = append(big.Bttv, NOIDONTTHINKSO{Code:emote.Code, Url: fmt.Sprintf("https://cdn.betterttv.net/emote/%s/1x", emote.ID),})
	}	
	for _, emote := range bttvChannelEmotes.SharedEmotes {
		big.Bttv = append(big.Bttv, NOIDONTTHINKSO{Code:emote.Code, Url: fmt.Sprintf("https://cdn.betterttv.net/emote/%s/1x", emote.ID),})
	}
	type FFZImages struct {
		The1X string  `json:"1x"`
		The2X *string `json:"2x"`
		The4X *string `json:"4x"`
	}
	type FFZUser struct {
		ID          int64       `json:"id"`
		Name        string        `json:"name"`
		DisplayName string `json:"displayName"`
	}
	type FFZEmote struct {
		ID        int64     `json:"id"`
		User      FFZUser   `json:"user"`
		Code      string    `json:"code"`
		Images    FFZImages `json:"images"`
		ImageType string `json:"imageType"`
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", "https://api.betterttv.net/3/cached/frankerfacez/emotes/global", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	resp, err = client.Do(req)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var ffzGlobalEmotes []FFZEmote
	err = json.Unmarshal(body, &ffzGlobalEmotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	for _, emote := range ffzGlobalEmotes {
		big.Ffz = append(big.Ffz, NOIDONTTHINKSO{Code:emote.Code, Url: emote.Images.The1X,})
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", "https://api.betterttv.net/3/cached/frankerfacez/users/twitch/50465063", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	resp, err = client.Do(req)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var ffzChannelEmotes []FFZEmote
	err = json.Unmarshal(body, &ffzChannelEmotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	for _, emote := range ffzChannelEmotes {
		big.Ffz = append(big.Ffz, NOIDONTTHINKSO{Code:emote.Code, Url: emote.Images.The1X,})
	}
	type SevenTVGlobalEmotes struct {
		ID         string        `json:"id"`
		Name       string        `json:"name"`
		Tags       []interface{} `json:"tags"`
		Immutable  bool          `json:"immutable"`
		Privileged bool          `json:"privileged"`
		Emotes     []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Flags     int    `json:"flags"`
			Timestamp int64  `json:"timestamp"`
			ActorID   string `json:"actor_id"`
			Data      struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Flags     int    `json:"flags"`
				Lifecycle int    `json:"lifecycle"`
				Listed    bool   `json:"listed"`
				Animated  bool   `json:"animated"`
				Owner     struct {
					ID          string `json:"id"`
					Username    string `json:"username"`
					DisplayName string `json:"display_name"`
					AvatarURL   string `json:"avatar_url"`
					Style       struct {
					} `json:"style"`
					Roles []string `json:"roles"`
				} `json:"owner"`
				Host struct {
					URL   string `json:"url"`
					Files []struct {
						Name       string `json:"name"`
						StaticName string `json:"static_name"`
						Width      int    `json:"width"`
						Height     int    `json:"height"`
						FrameCount int    `json:"frame_count"`
						Size       int    `json:"size"`
						Format     string `json:"format"`
					} `json:"files"`
				} `json:"host"`
			} `json:"data"`
			OriginID string `json:"origin_id,omitempty"`
		} `json:"emotes"`
		EmoteCount int `json:"emote_count"`
		Capacity   int `json:"capacity"`
		Origins    []struct {
			ID     string        `json:"id"`
			Weight int           `json:"weight"`
			Slices []interface{} `json:"slices"`
		} `json:"origins"`
		Owner struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			AvatarURL   string `json:"avatar_url"`
			Style       struct {
				Color int `json:"color"`
			} `json:"style"`
			Roles []string `json:"roles"`
		} `json:"owner"`
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", "https://7tv.io/v3/emote-sets/global", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	resp, err = client.Do(req)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var global_7tv_emotes SevenTVGlobalEmotes
	err = json.Unmarshal(body, &global_7tv_emotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	for _, emote := range global_7tv_emotes.Emotes {
		big.Seventv = append(big.Seventv, NOIDONTTHINKSO{Code: emote.Name, Url: fmt.Sprintf("https:%s/1x.webp", emote.Data.Host.URL),})
	}
	type SevenTVChannelEmotes struct {
		ID            string      `json:"id"`
		Platform      string      `json:"platform"`
		Username      string      `json:"username"`
		DisplayName   string      `json:"display_name"`
		LinkedAt      int64       `json:"linked_at"`
		EmoteCapacity int         `json:"emote_capacity"`
		EmoteSetID    interface{} `json:"emote_set_id"`
		EmoteSet      struct {
			ID         string        `json:"id"`
			Name       string        `json:"name"`
			Tags       []interface{} `json:"tags"`
			Immutable  bool          `json:"immutable"`
			Privileged bool          `json:"privileged"`
			Emotes     []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Flags     int    `json:"flags"`
				Timestamp int64  `json:"timestamp"`
				ActorID   string `json:"actor_id"`
				Data      struct {
					ID        string   `json:"id"`
					Name      string   `json:"name"`
					Flags     int      `json:"flags"`
					Tags      []string `json:"tags"`
					Lifecycle int      `json:"lifecycle"`
					Listed    bool     `json:"listed"`
					Animated  bool     `json:"animated"`
					Owner     struct {
						ID          string `json:"id"`
						Username    string `json:"username"`
						DisplayName string `json:"display_name"`
						AvatarURL   string `json:"avatar_url"`
						Style       struct {
							Color int `json:"color"`
						} `json:"style"`
						Roles []string `json:"roles"`
					} `json:"owner"`
					Host struct {
						URL   string `json:"url"`
						Files []struct {
							Name       string `json:"name"`
							StaticName string `json:"static_name"`
							Width      int    `json:"width"`
							Height     int    `json:"height"`
							FrameCount int    `json:"frame_count"`
							Size       int    `json:"size"`
							Format     string `json:"format"`
						} `json:"files"`
					} `json:"host"`
				} `json:"data"`
			} `json:"emotes"`
			EmoteCount int `json:"emote_count"`
			Capacity   int `json:"capacity"`
			Owner      struct {
				ID          string `json:"id"`
				Username    string `json:"username"`
				DisplayName string `json:"display_name"`
				AvatarURL   string `json:"avatar_url"`
				Style       struct {
					Color int `json:"color"`
				} `json:"style"`
				Roles []string `json:"roles"`
			} `json:"owner"`
		} `json:"emote_set"`
		User struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			CreatedAt   int64  `json:"created_at"`
			AvatarURL   string `json:"avatar_url"`
			Biography   string `json:"biography"`
			Style       struct {
				Color int `json:"color"`
			} `json:"style"`
			Editors []struct {
				ID          string `json:"id"`
				Permissions int    `json:"permissions"`
				Visible     bool   `json:"visible"`
				AddedAt     int64  `json:"added_at"`
			} `json:"editors"`
			Roles       []string `json:"roles"`
			Connections []struct {
				ID            string      `json:"id"`
				Platform      string      `json:"platform"`
				Username      string      `json:"username"`
				DisplayName   string      `json:"display_name"`
				LinkedAt      int64       `json:"linked_at"`
				EmoteCapacity int         `json:"emote_capacity"`
				EmoteSetID    interface{} `json:"emote_set_id"`
				EmoteSet      struct {
					ID         string        `json:"id"`
					Name       string        `json:"name"`
					Tags       []interface{} `json:"tags"`
					Immutable  bool          `json:"immutable"`
					Privileged bool          `json:"privileged"`
					Capacity   int           `json:"capacity"`
					Owner      interface{}   `json:"owner"`
				} `json:"emote_set"`
			} `json:"connections"`
		} `json:"user"`
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", "https://7tv.io/v3/users/twitch/50465063", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	resp, err = client.Do(req)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var channel_7tv_emotes SevenTVChannelEmotes
	err = json.Unmarshal(body, &channel_7tv_emotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	for _, emote := range channel_7tv_emotes.EmoteSet.Emotes {
		big.Seventv = append(big.Seventv, NOIDONTTHINKSO{Code: emote.Name, Url: fmt.Sprintf("https:%s/1x.webp", emote.Data.Host.URL),})

	}
	type Images struct {
		URL1X string `json:"url_1x"`
		URL2X string `json:"url_2x"`
		URL4X string `json:"url_4x"`
	}
	type TwitchGlobalEmote struct {
		ID        string      `json:"id"`
		Name      string      `json:"name"`
		Images    Images      `json:"images"`
		Format    []string    `json:"format"`
		Scale     []string    `json:"scale"`
		ThemeMode []string `json:"theme_mode"`
	}
	type TwitchGlobalEmotes struct {
		Data     []TwitchGlobalEmote `json:"data"`
		Template string              `json:"template"`
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", "https://api.twitch.tv/helix/chat/emotes/global", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	req.Header.Add("Client-ID", "g5zg0400k4vhrx2g6xi4hgveruamlv")
	req.Header.Add("Authorization", "Bearer "+"12345")
	resp, err = client.Do(req)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(400).SendString("I don't care that much...")
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var twitchGlobalEmotes TwitchGlobalEmotes
	err = json.Unmarshal(body, &twitchGlobalEmotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	for _, emote := range twitchGlobalEmotes.Data {

		big.Twitch = append(big.Twitch, NOIDONTTHINKSO{Code: emote.Name, Url: strings.Replace(emote.Images.URL1X, "/static/light", "/default/dark", -1),})
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", "https://api.twitch.tv/helix/chat/emotes?broadcaster_id=50465063", nil)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	req.Header.Add("Client-ID", "g5zg0400k4vhrx2g6xi4hgveruamlv")
	req.Header.Add("Authorization", "Bearer "+"12345")

	resp, err = client.Do(req)
	if err != nil {
		return c.Status(400).SendString("1I don't care that much...")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(400).SendString("I don't care that much...")
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	var twitchChannelEmotes TwitchGlobalEmotes
	err = json.Unmarshal(body, &twitchChannelEmotes)
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}

	for _, emote := range twitchChannelEmotes.Data {
		big.Twitch = append(big.Twitch, NOIDONTTHINKSO{Code: emote.Name, Url: strings.Replace(emote.Images.URL1X, "/static/light", "/default/dark", -1),})
	}
	payload, err := json.MarshalIndent(big, "", "\t")
	if err != nil {
		return c.Status(400).SendString("I don't care that much...")
	}
	c.Context().SetBody(payload)
	c.Context().SetContentType(fiber.MIMEApplicationJSON)
	return nil
}
func FeelsOkayMan(w *sync.WaitGroup, c string) {
	defer w.Done()
	fmt.Printf("[!] Started logging #%s!\n", c)
	client := twitch.NewAnonymousClient()
	client.OnPrivateMessage(
		func(message twitch.PrivateMessage) {
			var emotes []WineTime

			for _, i := range message.Emotes {
				emoteUrl := fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v2/%s/default/dark/1.0", i.ID)
				emotes = append(emotes, WineTime{Code: i.Name, Url: emoteUrl,})
			}
			err := os.MkdirAll(fmt.Sprintf("%s/forsenCD/", pwd), 0750)
			if err != nil {
				panic(err)
			}
			time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", message.Time.Year(), message.Time.Month(), message.Time.Day(), message.Time.Hour(), message.Time.Minute(), message.Time.Second())
			filename := fmt.Sprintf("%s/forsenCD/%02d-%02d-%d.json", pwd, message.Time.Day(), message.Time.Month(), message.Time.Year())
			_, err = os.Stat(filename)
			 if err != nil {
				chat := PoroSad{
					Frogs: []monkaE{},
				}
				chat.Frogs = append(chat.Frogs, monkaE{Username: message.User.Name, Message:  message.Message, Timestamp: time, Emotes: emotes})
				things, _ := json.MarshalIndent(chat, "", "\t")
				f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					panic(err)
				}
				if _, err := f.WriteString(string(things)); err != nil {
					panic(err)
				}
				if err := f.Close(); err != nil {
					panic(err)
				}
			} else {
				f, err := os.ReadFile(filename)
				if err != nil {
					panic(err)
				}
				var chat PoroSad
				err = json.Unmarshal(f, &chat)
				if err != nil {
					panic(err)
				}
				chat.Frogs = append(chat.Frogs, monkaE{Username: message.User.Name, Message:  message.Message, Timestamp: time, Emotes: emotes})
				pretty, _ := json.MarshalIndent(chat, "", "\t")
				f1, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					panic(err)
				}
				if _, err := f1.WriteString(string(pretty)); err != nil {
					panic(err)
				}
				if err := f1.Close(); err != nil {
					panic(err)
				}
			}
		})
	var latest twitch.ClearChatMessage
	client.OnClearChatMessage(
		func(message twitch.ClearChatMessage) {
			err := os.MkdirAll(fmt.Sprintf("%s/forsenCD/", pwd), 0750)
			if err != nil {
				panic(err)
			}
			time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", message.Time.Year(), message.Time.Month(), message.Time.Day(), message.Time.Hour(), message.Time.Minute(), message.Time.Second())
			filename := fmt.Sprintf("%s/forsenCD/%02d-%02d-%d.json", pwd, message.Time.Day(), message.Time.Month(), message.Time.Year())
			_, err = os.Stat(filename)
			 if err != nil {
				if latest.BanDuration == message.BanDuration && latest.RoomID == message.RoomID && message.Time.Year() == latest.Time.Year() && message.Time.Month() == latest.Time.Month() && latest.Time.Day() == message.Time.Day() && latest.Time.Hour() == message.Time.Hour() && message.Time.Minute() == latest.Time.Minute() && message.Time.Second() == latest.Time.Second() && latest.TargetUserID == message.TargetUserID && message.TargetUsername == latest.TargetUsername {
					return
				}
				chat := PoroSad{
					Frogs: []monkaE{},
				}
				payload := fmt.Sprintf("%s has been banned.", message.TargetUsername)
				if message.BanDuration != 0 {
					payload = fmt.Sprintf("%s has been timed out for %d seconds.", message.TargetUsername, message.BanDuration)
				}
				chat.Frogs = append(chat.Frogs, monkaE{Username: "bot", Message:  payload, Timestamp: time,})
				things, _ := json.MarshalIndent(chat, "", "\t")
				f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					panic(err)
				}
				if _, err := f.WriteString(string(things)); err != nil {
					panic(err)
				}
				if err := f.Close(); err != nil {
					panic(err)
				}
				latest = message
			} else {
				if latest.BanDuration == message.BanDuration && latest.RoomID == message.RoomID && message.Time.Year() == latest.Time.Year() && message.Time.Month() == latest.Time.Month() && latest.Time.Day() == message.Time.Day() && latest.Time.Hour() == message.Time.Hour() && message.Time.Minute() == latest.Time.Minute() && message.Time.Second() == latest.Time.Second() && latest.TargetUserID == message.TargetUserID && message.TargetUsername == latest.TargetUsername {
					return
				}
				f, err := os.ReadFile(filename)
				if err != nil {
					panic(err)
				}
				var chat PoroSad
				err = json.Unmarshal(f, &chat)
				if err != nil {
					panic(err)
				}
				payload := fmt.Sprintf("%s has been banned.", message.TargetUsername)
				if message.BanDuration != 0 {
					payload = fmt.Sprintf("%s has been timed out for %d seconds.", message.TargetUsername, message.BanDuration)
				}
				chat.Frogs = append(chat.Frogs, monkaE{Username: "bot", Message:  payload, Timestamp: time,})
				pretty, _ := json.MarshalIndent(chat, "", "\t")
				f1, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					panic(err)
				}
				if _, err := f1.WriteString(string(pretty)); err != nil {
					panic(err)
				}
				if err := f1.Close(); err != nil {
					panic(err)
				}
				latest = message
			}
	})
	client.OnUserNoticeMessage(func (message twitch.UserNoticeMessage) {
			err := os.MkdirAll(fmt.Sprintf("%s/forsenCD/", pwd), 0750)
			if err != nil {
				panic(err)
			}
			time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", message.Time.Year(), message.Time.Month(), message.Time.Day(), message.Time.Hour(), message.Time.Minute(), message.Time.Second())
			filename := fmt.Sprintf("%s/forsenCD/%02d-%02d-%d.json", pwd, message.Time.Day(), message.Time.Month(), message.Time.Year())
			_, err = os.Stat(filename)
			 if err != nil {
				chat := PoroSad{
					Frogs: []monkaE{},
				}
				chat.Frogs = append(chat.Frogs, monkaE{Username: "bot", Message:  message.SystemMsg, Timestamp: time,})
				things, _ := json.MarshalIndent(chat, "", "\t")
				f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					panic(err)
				}
				if _, err := f.WriteString(string(things)); err != nil {
					panic(err)
				}
				if err := f.Close(); err != nil {
					panic(err)
				}
			} else {
				f, err := os.ReadFile(filename)
				if err != nil {
					panic(err)
				}
				var chat PoroSad
				err = json.Unmarshal(f, &chat)
				if err != nil {
					panic(err)
				}
				chat.Frogs = append(chat.Frogs, monkaE{Username: "bot", Message:  message.SystemMsg, Timestamp: time,})
				pretty, _ := json.MarshalIndent(chat, "", "\t")
				f1, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					panic(err)
				}
				if _, err := f1.WriteString(string(pretty)); err != nil {
					panic(err)
				}
				if err := f1.Close(); err != nil {
					panic(err)
				}
			}
	})
	client.Join(c)
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
func gzipFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	gzipFile, err := os.Create(filePath + ".gz")
	if err != nil {
		panic(err)
	}
	defer gzipFile.Close()

	w := gzip.NewWriter(gzipFile)
	_, err = w.Write(content)
	if err != nil {
		panic(err)
	}
	w.Close()

	err = os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}
func forsenScoots(w *sync.WaitGroup) {
	defer w.Done()
	fmt.Printf("[!] Scanning for logs...\n")
	for {
		files, err := os.ReadDir(fmt.Sprintf("%s/forsenCD/", pwd))
		if err != nil {
			continue
		}
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".gz") {
				raw := strings.Split(file.Name(), ".")[0]
				day, _ := strconv.Atoi(strings.Split(raw, "-")[0])
				month, _ := strconv.Atoi(strings.Split(raw, "-")[1])
				year, _ := strconv.Atoi(strings.Split(raw, "-")[2])
				big_bang := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
				diff := time.Since(big_bang)
				if int(diff.Hours()/24) > 7 {
					gzipFile(fmt.Sprintf("%s/forsenCD/%s", pwd, file.Name()))
					fmt.Printf("[!] Archiving -> %s\n", fmt.Sprintf("%s/forsenCD/%s", pwd, file.Name()))
				}
			}
			 
		}
		time.Sleep(60*time.Second)
	}
}
func flushE(w *sync.WaitGroup) {
	defer w.Done()
	surskity := fiber.New(fiber.Config{DisableStartupMessage: true})
	surskity.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST",
	}))
	surskity.Static("/", "./assets")	
	surskity.All("/vods", happE)
	surskity.Get("/logs", FeelsStrongMan)
	surskity.Get("/emotes", scarE)
	fmt.Printf("[!] Listening on port 8000!\n")
	surskity.Listen(":8000")
}
func main() {
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	wg.Add(3)
	go FeelsOkayMan(wg, "surskity")
	go forsenScoots(wg)
	go flushE(wg)
	wg.Wait()
}
