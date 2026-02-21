package main

import (
	"log"
	"maps"
	"slices"

	"github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api/v2"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
)

type KeypressOperation string

const (
	PlayPause  KeypressOperation = "PlayPause"
	Play       KeypressOperation = "Play"
	Pause      KeypressOperation = "Pause"
	Previous   KeypressOperation = "Previous"
	Next       KeypressOperation = "Next"
	Shuffle    KeypressOperation = "Shuffle"
	LoopStatus KeypressOperation = "LoopStatus"
)

var operationsMap = map[string]KeypressOperation{
	"PlayPause":  PlayPause,
	"Play":       Play,
	"Pause":      Pause,
	"Previous":   Previous,
	"Next":       Next,
	"Shuffle":    Shuffle,
	"LoopStatus": LoopStatus,
}

type PlayerCtlKeyHandler struct {
	Client *dbus.Conn
}

func (v *PlayerCtlKeyHandler) Key(key api.KeyConfigV3, info api.StreamDeckInfoV1) {
	operation, ok := key.KeyHandlerFields["operation"]
	if !ok {
		log.Println("No MPRIS player operation specified")
	}
	op, ok := operationsMap[operation.(string)]
	if !ok {
		log.Println("Invalid MPRIS player operation specified")
	}
	playerName, ok := key.KeyHandlerFields["player_name"]
	var player *mpris.Player
	players, err := mpris.List(v.Client)
	if err != nil {
		log.Println(err)
		return
	}
	for _, p := range players {
		pl := mpris.New(v.Client, p)
		if ok && playerName != "" {
			if pl.GetShortName() == playerName {
				player = mpris.New(v.Client, p)
				break
			}
		} else {
			status, err := pl.GetPlaybackStatus()
			if err != nil {
				log.Println(err)
				continue
			}
			if status == mpris.PlaybackPlaying {
				player = pl
				break
			}
		}
	}
	if player == nil {
		return
	}
	switch op {
	case PlayPause:
		err = player.PlayPause()
	case Play:
		err = player.Play()
	case Pause:
		err = player.Pause()
	case Previous:
		err = player.Previous()
	case Next:
		err = player.Next()
	case Shuffle:
		shuffle, err := player.GetShuffle()
		if err != nil {
			log.Println(err)
			return
		}
		err = player.SetShuffle(!shuffle)
		break
	case LoopStatus:
		status, err := player.GetLoopStatus()
		if err != nil {
			log.Println(err)
			return
		}
		err = player.SetLoopStatus(getNextLoopStatus(status))
		break
	}
	if err != nil {
		log.Println(err)
		return
	}
}

func getNextLoopStatus(status mpris.LoopStatus) mpris.LoopStatus {
	switch status {
	case "None":
		return "Track"
	case "Track":
		return "Playlist"
	case "Playlist":
		return "None"
	}
	return "None"
}

func GetModule() streamdeckd.Module {

	return streamdeckd.Module{
		Name: "Playerctl",
		NewKey: func() api.KeyHandler {
			client, err := dbus.SessionBus()
			if err != nil {
				panic(err)
			}
			return &PlayerCtlKeyHandler{Client: client}
		},
		KeyFields: []api.Field{{Title: "Player Name", Name: "player_name", Type: "Text"}, {Title: "Operation", Name: "operation", Type: "Select", ListItems: slices.Collect(maps.Keys(operationsMap))}},
	}
}
