package main

import (
	"image"

	"github.com/unix-streamdeck/api/v2"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
)

type NoOpIconHandler struct {
	Running  bool
	Callback func(image image.Image)
}

func (c *NoOpIconHandler) Start(k api.KeyConfigV3, info api.StreamDeckInfoV1, callback func(image image.Image)) {
	if c.Callback == nil {
		c.Callback = callback
	}
	if c.Running {
		img := image.NewRGBA(image.Rect(0, 0, 96, 96)).SubImage(image.Rect(0, 0, 96, 96))
		callback(img)
	}
}

func (c *NoOpIconHandler) IsRunning() bool {
	return c.Running
}

func (c *NoOpIconHandler) SetRunning(running bool) {
	c.Running = running
}

func (c NoOpIconHandler) Stop() {
	c.Running = false
}

func GetModule() streamdeckd.Module {
	return streamdeckd.Module{NewIcon: func() api.IconHandler {
		return &NoOpIconHandler{Running: true}
	}, Name: "NoOp"}
}
