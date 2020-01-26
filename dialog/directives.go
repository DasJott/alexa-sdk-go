package dialog

type EchoDirectives []interface{}

func (d *EchoDirectives) Add(directive interface{}) {
	(*d) = append(*d, directive)
}

// SLOTS

// SlotDirective is the directive for slots
type SlotDirective struct {
	Type          string      `json:"type"`
	SlotToElicit  string      `json:"slotToElicit,omitempty"`
	SlotToConfirm string      `json:"slotToConfirm,omitempty"`
	UpdatedIntent *EchoIntent `json:"updatedIntent,omitempty"`
}

// AUDIO PLAYER

type AudioDirective struct {
	Type          string        `json:"type"`
	PlayBehaviour string        `json:"playBehavior",omitempty`  // either play or clear
	ClearBehavior string        `json:"clearBehavior",omitempty` // either play or clear
	Audioitem     EchoAudioItem `json:"audioItem"`
}

type EchoAudioItem struct {
	Stream AudioStream    `json:"stream"`
	Meta   *AudioMetadata `json:"metadata,omitempty"`
}

type AudioStream struct {
	Url                   string `json:"url"`
	Token                 string `json:"token"`
	ExpectedPreviousToken string `json:"expectedPreviousToken,omitempty"`
	OffsetInMilliseconds  int    `json:"offsetInMilliseconds"`
}

type AudioMetadata struct {
	Title           string    `json:"title"`
	Subtitle        string    `json:"subtitle"`
	Art             EchoImage `json:"art"`
	BackgroundImage EchoImage `json:"backgroundImage"`
}
