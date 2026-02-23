package main

import (
	"context"
	"errors"
	"image"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Endg4meZer0/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/unix-streamdeck/api/v2"
	"golang.org/x/sync/semaphore"
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
					_, err := pl.GetVolume()
					if err != nil {
						continue
					}
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
				img, err = v.FindImage(player)
				if err != nil {
					log.Println(err)
				}
				if img == nil {
					log.Println("Creating empty image of", player.GetShortName())
					img = image.NewNRGBA(image.Rect(0, 0, info.LcdWidth, info.LcdHeight))
					img, err = api.DrawText(img, player.GetShortName(), 0, "MIDDLE")
				} else {
					img = api.ResizeImageWH(img, info.LcdWidth, info.LcdHeight)
				}
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

func (v *VolumeLcdHandler) FindImage(player *mpris.Player) (image.Image, error) {
	metadata, err := player.GetMetadata()
	if err == nil && metadata != nil {
		artUrl, err := metadata.ArtURL()
		if err == nil && artUrl != "" {
			img, err := ExtractImage(artUrl)
			if err != nil {
				log.Println(err)
				err = nil
			}
			if img != nil {
				return img, nil
			}
		}
		if err != nil {
			log.Println(err)
		}
		err = nil
	}
	if err != nil {
		log.Println(err)
	}
	err = nil
	app, err := player.GetIdentity()
	if err != nil {
		log.Println(err)
		return nil, errors.New("couldn't find image")
	}

	if app == "" {
		return nil, errors.New("couldn't find image")
	}

	//apps := LinuxApps.GetApps()
	//
	//var icon string
	//for _, desktopEntry := range apps {
	//	if desktopEntry.Name == app {
	//		icon = desktopEntry.IconName
	//		break
	//	}
	//}
	//if icon == "" {
	//	return nil, errors.New("couldn't find image")
	//}
	//
	//img, err := ExtractImage(icon)
	//if err != nil {
	//	log.Println(err)
	//	return nil, errors.New("couldn't find image")
	//}
	//if img != nil {
	//	return img, nil
	//}

	return nil, errors.New("couldn't find image")
}

func ExtractImage(icon string) (image.Image, error) {
	match, err := regexp.MatchString(`(https?://)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)`, icon)
	if match {
		return getHttpImage(icon)
	}
	if err != nil {
		log.Println(err)
		err = nil
	}
	match, err = regexp.MatchString(`(/)+[a-zA-Z0-9\\\-_/ .]*\.+[a-z0-9A-Z]+`, icon)
	if match {
		return loadImage(icon)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return nil, errors.New("couldn't find image")
}

func getHttpImage(url string) (image.Image, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("Couldn't get Image from URL")
	}
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
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
			_, err := pl.GetVolume()
			if err != nil {
				log.Println("Caught the thing")
				continue
			}
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
func GetModule() api.Module {
	return api.Module{
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
