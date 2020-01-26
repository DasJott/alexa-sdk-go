package dialog

import "encoding/json"

type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes"`
	Response          EchoResponseBody       `json:"response"`
}

type EchoResponseBody struct {
	OutputSpeech     *EchoOutput    `json:"outputSpeech,omitempty"`
	Card             *EchoCard      `json:"card,omitempty"`
	Reprompt         *EchoReprompt  `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool           `json:"shouldEndSession"`
	Directives       EchoDirectives `json:"directives,omitempty"`
}

type EchoReprompt struct {
	OutputSpeech EchoOutput `json:"outputSpeech,omitempty"`
}

type EchoOutput struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

type EchoCard struct {
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	Text        string         `json:"text"`
	Image       *EchoCardImage `json:"image,omitempty"`
	Permissions []string       `json:"permissions,omitempty"`
}

type EchoCardImage struct {
	SmallImageURL string `json:"smallImageUrl"`
	LargeImageURL string `json:"largeImageUrl"`
}

// FUNCTIONS //////////////////////////////////////////////////////////////////

func NewResponse() *EchoResponse {
	er := &EchoResponse{
		Version: "1.0",
		Response: EchoResponseBody{
			ShouldEndSession: false,
		},
		SessionAttributes: make(map[string]interface{}),
	}
	return er
}

func (er *EchoResponse) OutputText(text string) *EchoResponse {
	er.Response.OutputSpeech = &EchoOutput{
		Type: "PlainText",
		Text: text,
	}
	return er
}

func (er *EchoResponse) OutputSSML(text string) *EchoResponse {
	er.Response.OutputSpeech = &EchoOutput{
		Type: "SSML",
		SSML: "<speak>" + voice(text) + "</speak>",
	}
	return er
}

func (er *EchoResponse) SlotDirective(directiveType, slotToElicit, slotToConfirm string, updatedIntent *EchoIntent) *EchoResponse {
	er.Response.Directives.Add(SlotDirective{
		Type:          directiveType,
		SlotToElicit:  slotToElicit,
		SlotToConfirm: slotToConfirm,
		UpdatedIntent: updatedIntent,
	})
	return er
}

// AudioDirective is for AudioPlayer.Play, AudioPlayer.Stop, AudioPlayer.ClearQueue
func (er *EchoResponse) AudioDirective(directiveType, behaviour, url, id string) *EchoResponse {
	// dir := AudioDirective{
	// 	Type: directiveType,
	// }
	// if AudioPlayer.ClearQueue {
	// 	dir.ClearBehavior = behavior
	// } else {
	// 	PlayBehaviour = behaviour
	// }
	// if url != "" {
	// 	dir.Audioitem.Stream.Url = url
	// 	dir.Audioitem.Stream.Token = id
	// 	dir.Audioitem.Stream.
	// }

	// er.Response.Directives.Add(&dir)
	return er
}

func (er *EchoResponse) Card(title string, content string) *EchoResponse {
	return er.SimpleCard(title, content)
}

func (er *EchoResponse) SimpleCard(title string, content string) *EchoResponse {
	er.Response.Card = &EchoCard{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}
	return er
}

func (er *EchoResponse) StandardCard(title string, content string, smallImg string, largeImg string) *EchoResponse {
	er.Response.Card = &EchoCard{
		Type:  "Standard",
		Title: title,
		Text:  content,
	}

	er.Response.Card.Image = &EchoCardImage{
		SmallImageURL: smallImg,
		LargeImageURL: largeImg,
	}

	return er
}

func (er *EchoResponse) LinkAccountCard() *EchoResponse {
	er.Response.Card = &EchoCard{
		Type: "LinkAccount",
	}
	return er
}

func (er *EchoResponse) AskPermissionCard(permissions []string) *EchoResponse {
	er.Response.Card = &EchoCard{
		Type:        "AskForPermissionsConsent",
		Permissions: permissions,
	}
	return er
}

func (er *EchoResponse) RepromptText(text string) *EchoResponse {
	if text != "" {
		er.Response.Reprompt = &EchoReprompt{
			OutputSpeech: EchoOutput{
				Type: "PlainText",
				Text: text,
			},
		}
	}
	return er
}

func (er *EchoResponse) RepromptSSML(text string) *EchoResponse {
	if text != "" {
		er.Response.Reprompt = &EchoReprompt{
			OutputSpeech: EchoOutput{
				Type: "SSML",
				SSML: "<speak>" + voice(text) + "</speak>",
			},
		}
	}
	return er
}
func (er *EchoResponse) EndSession() *EchoResponse {
	er.Response.ShouldEndSession = true
	return er
}

func (er *EchoResponse) Data() []byte {
	jsonStr, err := json.Marshal(&er)
	if err != nil {
		return nil
	}
	return jsonStr
}

func (er *EchoResponse) String() string {
	return string(er.Data())
}
