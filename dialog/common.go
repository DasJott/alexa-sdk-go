package dialog

import "github.com/dasjott/alexa-sdk-go/ssml"

// Objects common to request and response

var voice = func(s string) string { return s }

// SetVoice sets a voice to be used for all following output
func SetVoice(name string) {
	voice = ssml.NewVoice(name)
}

// EchoIntent is the json part for an intent
type EchoIntent struct {
	Name               string               `json:"name"`
	Slots              map[string]*EchoSlot `json:"slots"`
	ConfirmationStatus string               `json:"confirmationStatus,omitempty"`
}

type EchoImage struct {
	Sources []*EchoImageSize `json:"sources"`
}

type EchoImageSize struct {
	URL          string `json:"url"`
	Size         string `json:"size"`
	WidthPixels  int    `json:"widthPixels"`
	HeightPixels int    `json:"heightPixels"`
}
