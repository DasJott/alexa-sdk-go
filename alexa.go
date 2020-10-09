package alexa

import (
	"github.com/dasjott/alexa-sdk-go/api"
	"github.com/dasjott/alexa-sdk-go/dialog"
)

// AppID is the ID of the corresponding skill
var AppID string

// Handlers are intent functions to be called by name
var Handlers IntentHandlers

var CanFulfillIntent func(*dialog.EchoIntent)

// LocaleStrings are all localized strings
var LocaleStrings Localisation

// GetTranslation is called with the current locale code.
// You must provide the according Translation.
var GetTranslation func(locale string) Translation

// BeforeHandler can be set with a function to implement any checking before every intent.
// It returns true for going on with the actual intent or false to skip.
// Remember to implement a appropriate message to the user on skipping!
var BeforeHandler func(*Context)

// Handle is the function you hand over to the lambda.start
var Handle = func(req *dialog.EchoRequest) (*dialog.EchoResponse, error) {
	if req == nil {
		panic("Echo request is nil")
	}

	// if !req.VerifyTimestamp() {
	// 	return "", errors.New("invalid timestamp")
	// }
	if AppID != "" && !req.VerifyAppID(AppID) {
		panic("invalid app id")
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
		handlers:   Handlers,
		response:   dialog.NewResponse(),
		translator: trans,
		attributes: req.Session.Attributes,

		System: req.Context.System,
		Intent: req.Request.Intent,
		Time:   req.GetTime(),
	}

	c.start(req)
	return c.getResult()
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
