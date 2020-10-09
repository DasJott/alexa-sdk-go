package alexa

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dasjott/alexa-sdk-go/api"
	"github.com/dasjott/alexa-sdk-go/dialog"
	"github.com/dasjott/alexa-sdk-go/intents"
	"github.com/dasjott/alexa-sdk-go/intents/canfulfill"
)

var (
	// AppID is the ID of the corresponding skill
	AppID string

	// Handlers are intent functions to be called by name
	Handlers IntentHandlers

	CanFulfillIntent func(*dialog.EchoIntent) *canfulfill.Response

	// LocaleStrings are all localized strings
	LocaleStrings Localisation

	// GetTranslation is called with the current locale code.
	// You must provide the according Translation.
	GetTranslation func(locale string) Translation

	// BeforeHandler can be set with a function to implement any checking before every intent.
	// It returns true for going on with the actual intent or false to skip.
	// Remember to implement a appropriate message to the user on skipping!
	BeforeHandler func(*Context)

	random *rand.Rand
)

// Handle is the function you hand over to the lambda.start
var Handle = func(req *dialog.EchoRequest) (intents.Response, error) {
	if req == nil {
		panic("Echo request is nil")
	}

	// if !req.VerifyTimestamp() {
	// 	return "", errors.New("invalid timestamp")
	// }
	if AppID != "" && !req.VerifyAppID(AppID) {
		panic("invalid app id")
	}

	intent := req.GetIntentName()
	fmt.Printf("intent: %s\n", intent)

	if resp := handleXIntents(req); resp != nil {
		return resp, nil
	}

	if Handlers == nil {
		panic("no handlers set")
	}

	var trans *Translator
	if GetTranslation != nil {
		if langmap := GetTranslation(req.Request.Locale); langmap != nil {
			loc := Localisation{req.Request.Locale: langmap}
			trans = loc.GetTranslator(req.Request.Locale)
		}
	} else if LocaleStrings != nil {
		trans = LocaleStrings.GetTranslator(req.Request.Locale)
	}
	if trans == nil {
		panic("language " + req.Request.Locale + " not implemented")
	}

	c := &Context{
		request:    req,
		response:   dialog.NewResponse(),
		translator: trans,
		attributes: req.Session.Attributes,

		System: req.Context.System,
		Intent: req.Request.Intent,
		Time:   req.GetTime(),
	}

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
		if handler, exists := Handlers[intent]; exists {
			handler(c)
		} else if handler, exists := Handlers["Unhandled"]; exists {
			handler(c)
		} else {
			panic("no handler found")
		}
	}

	return handleUserIntents(c)
}

func handleXIntents(req *dialog.EchoRequest) *intents.XResponse {
	name := req.GetRequestType()

	switch {
	case name == "CanFulfillIntentRequest" && CanFulfillIntent != nil:
		content := CanFulfillIntent(req.Request.Intent)
		return canfulfill.NewXResponse(content)
	}
	return nil
}

func handleUserIntents(c *Context) (*dialog.EchoResponse, error) {
	c.progressWait()
	c.response.SessionAttributes = c.attributes
	return c.response, c.err
}

// IntentHandler function for the handler
type IntentHandler func(*Context)

// IntentHandlers for collecting the handler functions
type IntentHandlers map[string]IntentHandler

// Add adds a handler afterwards.
func (h IntentHandlers) Add(name string, handler IntentHandler) {
	h[name] = handler
}

// MultiHandler if you need more than one handler for an intent
func MultiHandler(handlers ...IntentHandler) IntentHandler {
	return func(c *Context) {
		count := len(handlers)
		for i := 0; i < count && !c.abort; i++ {
			(handlers[i])(c)
		}
	}
}

// API sets up a client to call the alexa api
func API(c *Context) *api.Client {
	return api.NewClient(c.request.Context.System)
}
