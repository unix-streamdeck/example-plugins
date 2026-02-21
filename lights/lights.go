package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/unix-streamdeck/api/v2"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
)

type LightsKeyHandler struct{}

func (LightsKeyHandler) Key(key api.KeyConfigV3, info api.StreamDeckInfoV1) {
	entityId, ok := key.KeyHandlerFields["entity_id"]
	if !ok {
		return
	}
	reqBody, err := json.Marshal(map[string]string{
		"entity_id": entityId.(string),
	})
	if err != nil {
		print(err)
		return
	}
	apiKey, ok := key.KeyHandlerFields["api_key"]
	if !ok {
		return
	}
	domain, ok := key.KeyHandlerFields["domain"]
	if !ok {
		return
	}
	service, ok := key.KeyHandlerFields["service"]
	if !ok {
		return
	}
	baseUrl, ok := key.KeyHandlerFields["base_url"]
	if !ok {
		return
	}
	req, err := http.NewRequest("POST", "http://"+baseUrl.(string)+"/api/services/"+domain.(string)+"/"+service.(string), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey.(string))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}

func GetModule() streamdeckd.Module {

	return streamdeckd.Module{
		Name:      "Lights",
		NewKey:    func() api.KeyHandler { return &LightsKeyHandler{} },
		KeyFields: []api.Field{{Title: "Domain", Name: "domain", Type: "Text"}, {Title: "Service", Name: "service", Type: "Text"}, {Title: "Entity Id", Name: "entity_id", Type: "Text"}, {Title: "Api Key", Name: "api_key", Type: "Text"}, {Title: "Base Url", Name: "base_url", Type: "Text"}},
	}
}
