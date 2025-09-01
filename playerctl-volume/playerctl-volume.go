package main

import (
	"context"
	"github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
	"golang.org/x/sync/semaphore"
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type VolumeLcdHandler struct {
	Running    bool
	Quit       chan bool
	Lock       *semaphore.Weighted
	Buff       image.Image
	PlayerName string
	Volume     int
	FirstLoop  bool
	Client     *dbus.Conn
}

func (v *VolumeLcdHandler) Start(knob api.KnobConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if v.Quit == nil {
		v.Quit = make(chan bool)
	}
	if v.Lock == nil {
		v.Lock = semaphore.NewWeighted(1)
	}
	if v.Buff == nil {
		v.Buff = v.GetImage("icon", knob, info)
	}

	playerName, ok := knob.LcdHandlerFields["player_name"]
	if ok {
		v.PlayerName = playerName.(string)
	}
	v.Running = true
	v.Run(knob, info, callback)
}
func (v *VolumeLcdHandler) IsRunning() bool {
	return v.Running
}

func (v *VolumeLcdHandler) SetRunning(running bool) {
	v.Running = running
}

func (v *VolumeLcdHandler) Stop() {
	v.Running = false
	v.Quit <- true
}

func (v *VolumeLcdHandler) GetImage(index string, knob api.KnobConfigV3, info api.StreamDeckInfoV1) image.Image {
	path, ok := knob.LcdHandlerFields[index]
	if !ok {
		log.Println("image missing: " + index)
		return nil
	}
	f, err := os.Open(path.(string))
	defer f.Close()
	if err != nil {
		log.Println(err)
		return nil
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err)
		return nil
	}
	return api.ResizeImageWH(img, info.LcdWidth, info.LcdHeight)
}

func (v *VolumeLcdHandler) Run(knob api.KnobConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	ctx := context.Background()
	err := v.Lock.Acquire(ctx, 1)
	if err != nil {
		return
	}
	for {
		select {
		case <-v.Quit:
			return
		default:
			var player *mpris.Player
			players, err := mpris.List(v.Client)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, p := range players {
				pl := mpris.New(v.Client, p)
				if v.PlayerName != "" {
					if pl.GetShortName() == v.PlayerName {
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
				log.Println("Found no player ")
				continue
			}
			volume, err := player.GetVolume()
			if err != nil {
				log.Println(err)
				continue
			}
			v.Volume = int(math.Round(volume * 100.0))
			text := strconv.Itoa(v.Volume) + "%"
			var img image.Image
			img = v.Buff
			if img == nil {
				log.Println("Creating empty image of", player.GetShortName())
				img = image.NewNRGBA(image.Rect(0, 0, info.LcdWidth, info.LcdHeight))
				img, err = api.DrawText(img, player.GetShortName(), 0, "MIDDLE")
			}
			imgParsed, err := api.DrawText(img, text, 24, "BOTTOM")
			if err != nil {
				log.Println(err)
			} else {
				callback(imgParsed)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

type VolumeKnobOrTouchHandler struct {
	Client *dbus.Conn
}

func (v *VolumeKnobOrTouchHandler) Input(knob api.KnobConfigV3, info api.StreamDeckInfoV1, event api.InputEvent) {
	playerName, ok := knob.KnobOrTouchHandlerFields["player_name"]
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
	volume, err := player.GetVolume()
	if err != nil {
		log.Println(err)
		return
	}
	volume = math.Round(volume * 100.0)

	log.Println("before", volume)

	if event.EventType == api.KNOB_CCW {
		volume -= 1.0
	} else if event.EventType == api.KNOB_CW {
		volume += 1.0
	}
	volume /= 100.0
	err = player.SetVolume(volume)
	if err != nil {
		log.Println(err)
		return
	}
}
func GetModule() streamdeckd.Module {
	return streamdeckd.Module{
		NewLcd: func() api.LcdHandler {
			client, err := dbus.SessionBus()
			if err != nil {
				panic(err)
			}
			return &VolumeLcdHandler{Running: true, Lock: semaphore.NewWeighted(1), FirstLoop: true, Client: client}
		},
		LcdFields: []api.Field{{Title: "Icon", Name: "icon", Type: "File", FileTypes: []string{".png", ".jpg", ".jpeg"}}, {Title: "Player Name", Name: "player_name", Type: "Text"}},
		NewKnobOrTouch: func() api.KnobOrTouchHandler {
			client, err := dbus.SessionBus()
			if err != nil {
				panic(err)
			}
			return &VolumeKnobOrTouchHandler{Client: client}
		},
		KnobOrTouchFields: []api.Field{{Title: "Player Name", Name: "player_name", Type: "Text"}},
		Name:              "PlayerCtlVolume",
	}
}
