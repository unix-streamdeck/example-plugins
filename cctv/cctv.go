package main

import (
	"errors"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
	"golang.org/x/sync/semaphore"
	"image"
	"log"
	"net/http"
	"time"
)

type CCTVIconHandler struct {
	Status   bool
	Running  bool
	Lock     *semaphore.Weighted
	Callback func(image image.Image)
	Quit     chan bool
	Url      string
}

func (c *CCTVIconHandler) Start(k api.KeyConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	c.Running = true
	if c.Lock == nil {
		c.Lock = semaphore.NewWeighted(1)
	}
	if c.Quit == nil {
		c.Quit = make(chan bool)
	}
	if c.Callback == nil {
		c.Callback = callback
	}
	url, ok := k.IconHandlerFields["url"]
	if !ok {
		log.Println("Missing fields: url")
		c.Quit <- true
		c.Running = false
		return
	}
	c.Url = url.(string)
	go c.loop()
}

func (c *CCTVIconHandler) loop() {
	for {
		select {
		case <-c.Quit:
			return
		default:
			c.updateIcon()
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func (c *CCTVIconHandler) updateIcon() {
	img, err := getImage(c.Url)
	if err != nil {
		log.Println(err)
		return
	}
	c.Callback(img)
}

func (c *CCTVIconHandler) IsRunning() bool {
	return c.Running
}

func (c *CCTVIconHandler) SetRunning(running bool) {
	c.Running = running
}

func (c *CCTVIconHandler) Stop() {
	c.Running = false
	c.Quit <- true
}

type CCTVKeyHandler struct{}

func (CCTVKeyHandler) Key(key api.KeyConfigV3, info api.StreamDeckInfoV1) {
	if key.IconHandler != "CCTV" {
		return
	}
	handler := key.IconHandlerStruct.(*CCTVIconHandler)
	handler.updateIcon()
}

func GetModule() streamdeckd.Module {

	return streamdeckd.Module{
		Name:       "CCTV",
		NewIcon:    func() api.IconHandler { return &CCTVIconHandler{Running: true, Lock: semaphore.NewWeighted(1)} },
		IconFields: []api.Field{{Title: "URL", Name: "url", Type: "Text"}},
	}
}

func getImage(url string) (image.Image, error) {
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
