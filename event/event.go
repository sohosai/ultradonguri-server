package event

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.com/donguri-back/client"
	"example.com/donguri-back/telop"
	"example.com/donguri-back/util"
)

func GetEvents() ([]Event, error) {
	file, err := os.Open("events.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []Event
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&events); err != nil {
		return nil, err
	}

	// 余剰トークン検出
	if dec.More() {
		return nil, fmt.Errorf("trailing data after JSON array")
	}
	return events, nil
}

func FindMusicById(id string) (*Music, error) {
	events, err := GetEvents()

	if err != nil {
		return nil, err
	}

	for _, event := range events {
		for _, music := range event.Musics {
			if id == music.ID {
				return &music, nil
			}
		}
	}

	return nil, fmt.Errorf("The music with id %s doesn't exist.", id)
}

func FindEventById(id string) (*Event, error) {
	events, err := GetEvents()

	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if id == event.ID {
			return &event, nil
		}
	}

	return nil, fmt.Errorf("The music with id %s doesn't exist.", id)
}

type HM int

func (t *HM) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("time must be a string: %w", err)
	}
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid time %q (want HH:MM)", s)
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return fmt.Errorf("invalid time %q (0<=HH<=23, 0<=MM<=59)", s)
	}
	*t = HM(h*60 + m)
	return nil
}

func (t HM) MarshalJSON() ([]byte, error) {
	min := int(t)
	return json.Marshal(fmt.Sprintf("%02d:%02d", min/60, min%60))
}

type Music struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	ShouldBeMuted bool   `json:"should_be_muted"`
	Intro         string `json:"intro"`
}

type Event struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Performer   string  `json:"performer"`
	Description string  `json:"description"`
	StartsAt    HM      `json:"starts_at"`
	EndsAt      HM      `json:"ends_at"`
	Musics      []Music `json:"musics"`
}

type MusicTelopClient struct {
	telopClient *telop.TelopClient
	music       util.Option[Music]
}

const MUSIC_TELOP_NAME = "Music Telop"

func NewMusicTelopClient(sharedClient *client.SharedClient, sceneName string) (*MusicTelopClient, error) {
	telopClient, err := telop.NewTelopClient(sharedClient, sceneName, MUSIC_TELOP_NAME)

	if err != nil {
		return nil, err
	}

	return &MusicTelopClient{
		telopClient: telopClient,
		music:       util.None[Music](),
	}, nil
}

// 保持しているmusicを変更する。
func (self *MusicTelopClient) SetMusic(music *Music) {
	self.music = util.Some(*music)
}

// musicの変更をOBSのテロップへ反映する
func (self *MusicTelopClient) ApplyMusicChange() {
	if self.music.IsSome() {
		self.telopClient.SetText(self.music.Unwrap().Title)
		self.telopClient.SetVisible(true)
	}
}

// 保持しているmusicの参照を返す
func (self *MusicTelopClient) GetMusic() (*Music, error) {
	if self.music.IsNone() {
		return nil, fmt.Errorf("Music telop is not configured yet.")
	}

	music := self.music.Unwrap()
	return &music, nil

}

type EventTelopClient struct {
	telopClient *telop.TelopClient
	event       util.Option[Event]
}

const EVENT_TELOP_NAME = "Event Telop"

func NewEventTelopClient(sharedClient *client.SharedClient, sceneName string) (*EventTelopClient, error) {
	telopClient, err := telop.NewTelopClient(sharedClient, sceneName, EVENT_TELOP_NAME)

	if err != nil {
		return nil, err
	}

	return &EventTelopClient{
		telopClient: telopClient,
		event:       util.None[Event](),
	}, nil
}

// 保持しているeventを変更する。
func (self *EventTelopClient) SetEvent(event *Event) {
	self.event = util.Some(*event)
}

// musicの変更をOBSのテロップへ反映する
func (self *EventTelopClient) ApplyEventChange() {
	if self.event.IsSome() {
		self.telopClient.SetText(self.event.Unwrap().Title)
		self.telopClient.SetVisible(true)
	}
}

// 保持しているmusicの参照を返す
func (self *EventTelopClient) GetEvent() (*Event, error) {
	if self.event.IsNone() {
		return nil, fmt.Errorf("Event telop is not configured yet.")
	}

	music := self.event.Unwrap()
	return &music, nil

}
