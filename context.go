package alexa

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dasjott/alexa-sdk-go/dialog"
)

var random *rand.Rand

// Context is the object sent to every intent, collecting infos for response
type Context struct {
	attributes
	request    *dialog.EchoRequest
	handlers   IntentHandlers
	response   *dialog.EchoResponse
	translator *Translator
	err        error
	abort      bool
	progress   *dialog.ProgressRequest
	// System contains informations about the calling Device and User
	System *dialog.EchoSystem
	// Intent is the intents name
	Intent *dialog.EchoIntent
	// Time is the requests timestamp as go time.Time
	Time time.Time
}

func (c *Context) start(req *dialog.EchoRequest) {
	if c.translator == nil {
		panic("no translator set")
	}
	if c.attributes == nil {
		c.attributes = make(attributes)
	}

	random = rand.New(rand.NewSource(time.Now().Unix()))

	if BeforeHandler != nil {
		BeforeHandler(c)
	}
	if !c.abort {
		c.onIntent(req.GetIntentName())
	}
}

func (c *Context) onIntent(name string) {
	fmt.Printf("intent: %s\n", name)
	if handler, exists := c.handlers[name]; exists {
		handler(c)
	} else if handler, exists := c.handlers["Unhandled"]; exists {
		handler(c)
	} else {
		panic("no handler found")
	}
}

func (c *Context) getResult() (*dialog.EchoResponse, error) {
	c.progressWait()
	c.response.SessionAttributes = c.attributes
	return c.response, c.err
}

func (c *Context) progressWait() {
	if c.progress != nil {
		c.progress.Wait()
		c.progress = nil
	}
}

// Slot gets a slot by name. The pointer is never nil.
func (c *Context) Slot(name string) *Slot {
	if c.request.Request.Intent.Slots != nil {
		if slot, ok := c.request.Request.Intent.Slots[name]; ok {
			return slotFromEchoSlot(&slot)
		}
	}
	return &Slot{}
}

// NewSession determines whether this is a new session that was opened with this call
func (c *Context) NewSession() bool {
	return c.request.Session.New
}

// SessionID is the unique ID of this session
func (c *Context) SessionID() string {
	return c.request.Session.SessionID
}

// Locale gets the locale string like one of:
// de-DE, en-AU, en-CA, en-GB, en-IN, en-US, ja-JP, fr-FR
func (c *Context) Locale() string {
	return c.request.Request.Locale
}

// DialogState gets the current state of the dialog
func (c *Context) DialogState() string {
	return c.request.Request.DialogState
}

// T gets a translated string according to the given key. If the value is an array, a random value is chosen.
func (c *Context) T(key ...string) string {
	for i, k := range key {
		key[i] = c.translator.GetString(k)
	}
	return strings.Join(key, " ")
}

// TA gets a translated string array according to the given key.
func (c *Context) TA(key string) []string {
	return c.translator.GetArray(key)
}

// TR gets a translated string according to the given key. If the value is an array, a random value is chosen.
// Variables in {brackets} will be replaced. Use either the alexa.R or a struct for providing variables (tag name would be 'alexa')!
func (c *Context) TR(key string, replace interface{}) string {
	if repR, ok := replace.(R); ok {
		return c.translator.GetStringAndReplace(key, repR)
	}
	return c.translator.GetStringWithVariables(key, replace)
}

// Tell something to the user
func (c *Context) Tell(speech string) *Cardable {
	c.response.EndSession().OutputSSML(speech)
	return &Cardable{c}
}

// Ask the user something
func (c *Context) Ask(speechOutput string, repromptSpeech ...string) *Cardable {
	c.response.OutputSSML(speechOutput)
	if count := len(repromptSpeech); count > 0 {
		reprompt := repromptSpeech[random.Intn(count)]
		c.response.RepromptSSML(reprompt)
	}
	return &Cardable{c}
}

// AudioPlay is for AudioPlayer.Play
func (c *Context) AudioPlay(audio, id, speech string) *Cardable {
	// c.response.OutputSSML(speech).AudioDirective("AudioPlayer.Play", "REPLACE_ALL")
	return &Cardable{c}
}

// AudioStop is for AudioPlayer.Stop, AudioPlayer.ClearQueue
func (c *Context) AudioStop(speech string) *Cardable {
	// c.response.OutputSSML(speech).AudioDirective("AudioPlayer.Stop", "")
	return &Cardable{c}
}

// ElicitSlot action to fullfill a slot of a certain intent
func (c *Context) ElicitSlot(slotToElicit, speechOutput, repromptSpeech string, updatedIntent *dialog.EchoIntent) *Cardable {
	c.response.OutputSSML(speechOutput).RepromptSSML(repromptSpeech).SlotDirective("Dialog.ElicitSlot", slotToElicit, "", updatedIntent)
	return &Cardable{c}
}

// ConfirmSlot confirm a slot value by Alexa
func (c *Context) ConfirmSlot(slotToConfirm, speechOutput, repromptSpeech string, updatedIntent *dialog.EchoIntent) *Cardable {
	c.response.OutputSSML(speechOutput).RepromptSSML(repromptSpeech).SlotDirective("Dialog.ConfirmSlot", "", slotToConfirm, updatedIntent)
	return &Cardable{c}
}

// ConfirmIntent confirm all the slots given to the intent by Alexa
func (c *Context) ConfirmIntent(speechOutput, repromptSpeech string, updatedIntent *dialog.EchoIntent) *Cardable {
	c.response.OutputSSML(speechOutput).RepromptSSML(repromptSpeech).SlotDirective("Dialog.ConfirmIntent", "", "", updatedIntent)
	return &Cardable{c}
}

// Delegate a slot fullfillment to Alexa
func (c *Context) Delegate(updatedIntent *dialog.EchoIntent) {
	c.response.SlotDirective("Dialog.Delegate", "", "", updatedIntent)
}

// Progress sends a progress for the user to be entertained while waiting
func (c *Context) Progress(speech string) {
	c.progressWait()
	c.progress = dialog.NewProgressRequest(speech, c.request.Request.RequestID, c.System)
	if c.progress != nil {
		c.progress.Send()
	}
}

// Now returns the time of the request on users side
func (c *Context) Now() time.Time {
	return time.Now()
}

// Abort prevents the execution of a following handler within an alexa.MultiHandler chain.
func (c *Context) Abort() {
	c.abort = true
}

// Cardable is returned by functions, to make them cardable
type Cardable struct {
	cc *Context
}

// SimpleCard adds a simple card to the response
func (c *Cardable) SimpleCard(title, content string) {
	c.cc.response.SimpleCard(title, content)
}

// StandardCard adds a standard card to the response
func (c *Cardable) StandardCard(title, content, smallImageURL, largeImageURL string) {
	c.cc.response.StandardCard(title, content, smallImageURL, largeImageURL)
}

// LinkAccountCard adds a link account card to the response
func (c *Cardable) LinkAccountCard() {
	c.cc.response.LinkAccountCard()
}

// AskPermissionCard adds a "ask for permission" card to the response
// you can use the constants from this package, prefixed with Permission
func (c *Cardable) AskPermissionCard(permissions []string) {
	c.cc.response.AskPermissionCard(permissions)
}
