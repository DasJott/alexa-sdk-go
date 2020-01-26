package dialog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dasjott/alexa-sdk-go/test"
)

var endpoint = "/v1/directives"

// ProgressRequest is a request to make Alexa say things before an intent reponse is ready
type ProgressRequest struct {
	Header struct {
		RequestID string `json:"requestId"`
	} `json:"header"`
	Directive struct {
		Type   string `json:"type"`
		Speech string `json:"speech"`
	} `json:"directive"`

	// for internal use, not for json
	system *EchoSystem
	wait   chan int
}

// NewProgressRequest creates a new ProgressRequest
func NewProgressRequest(speech, requestID string, sys *EchoSystem) *ProgressRequest {
	if sys != nil && sys.APIEndpoint != "" {
		p := ProgressRequest{}
		p.Header.RequestID = requestID
		p.Directive.Speech = "<speak>" + voice(speech) + "</speak>"
		p.Directive.Type = "VoicePlayer.Speak"
		p.system = sys
		return &p
	}
	return nil
}

// Send actually sends the request to where it belongs
func (p *ProgressRequest) Send() {
	data, _ := json.Marshal(p)
	req, err := http.NewRequest(http.MethodPost, p.system.APIEndpoint+endpoint, bytes.NewReader(data))
	req.Header.Add("Authorization", "Bearer "+p.system.APIAccessToken)
	req.Header.Add("Content-Type", "application/json")

	if err == nil {
		p.wait = make(chan int)
		go func() {
			var resp *http.Response
			var err error
			if test.RequestHandler != nil {
				resp = test.RequestHandler(req)
			} else {
				client := http.Client{}
				resp, err = client.Do(req)
			}
			if err != nil {
				fmt.Println("progress error: ", err.Error())
			}
			code := 0
			if resp != nil {
				code = resp.StatusCode
			}
			p.wait <- code
		}()
	} else {
		fmt.Println("progress error: ", err.Error())
	}
}

// Wait returns the response code of the progress after waiting for it
func (p *ProgressRequest) Wait() int {
	if p.wait == nil {
		return 0
	}
	return <-p.wait
}
