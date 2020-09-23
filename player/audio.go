package player

// NewAudio returns a pointer to a new Audio object.
func NewAudio(url string) *Audio {
	return &Audio{
		url: url,
	}
}

// Audio brings sound to Alexa.
type Audio struct {
	url string
}

// Play plays the sound immediately, within the speech.
func (a *Audio) Play() {
	// play immediately
}

// Stream starts the sound as an audio stream
func (a *Audio) Stream() {
	// audio streaming
}
