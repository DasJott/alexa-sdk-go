# Alexa SDK for GO
Alexa SDK for native go on AWS Lambda<br>

Installation:
``` bash
go get https://github.com/dasjott/alexa-sdk-go
```
<br>

# Usage / Quick Start
## __main__
Your main function is fairly clear
```go
func main() {
	alexa.AppID = "your-skill-id-1234"
	alexa.Handlers = handlers
	alexa.LocaleStrings = locales
	lambda.Start(alexa.Handle)
}
```

## __handlers__
The 'handlers' variable mentioned in main:
```go
var handlers = alexa.IntentHandlers{
	"LaunchRequest": func(c *alexa.Context) {
		sayHello(c, c.T("HELLO"))
	},
	"AMAZON.HelpIntent": func(c *alexa.Context) {
		c.Ask(c.T("HELP"))
	},
	"AMAZON.CancelIntent": bye,
	"AMAZON.StopIntent": bye,
	"AMAZON.NoIntent": bye,
	"AMAZON.YesIntent": func(c *alexa.Context) {
		sayHello(c, c.T("HELLO"))
	},
}

func bye(c *alexa.Context) {
	c.Tell(c.T("BYE"))
}
func sayHello(c *alexa.Context, speech string) {
	c.Tell(speech)
}
```

## __locales__
The 'locales' variable mentioned in main:
```go
var locales = alexa.Localisation{
	"de-DE": alexa.Translation{
		"HELLO":     []string{"Hallo.", "Guten Tag."},
		"BYE":       []string{"Tschüss.", "Bis dann."},
		"HELP":      "Ich kann Hallo sagen. Möchtest du jetzt Hallo sagen?",
	},
	"en-US": alexa.Translation{
		"HELLO":     "Hello.",
		"BYE":       []string{"Bye.", "See you."},
		"HELP":      "I can say hello. Do you want me to say hello?",
	},
}
```
As you may have noticed, some fields are string arrays but are accessed like single strings.<br>
If you do so, the translation kit then chooses one string randomly from the given array, enabling you an easy way of implementing a more natural response behaviour of Alexa.

# __Alexa__
The alexa package is the main package of this SDK. You already saw it in the previous section.<br>
It provides you with the following functionality:<br>
- `alexa.AppID = "your-skill-id-1234`<br>
	Fill in the skills ID here. It is checked on every request and a panic is released on a mismatch.<br>
	Let it empty to skip that check (and panic)<br>

- `alexa.Handlers = alexa.IntentHandlers{}`<br>
	The IntentHandlers type is a mapping of intent names to alexa.IntentHandler functions.<br>
	[See here, how to use it](#handlers)

- `alexa.IntentHandler`
	The type for a intent function: `func(c *alexa.Context)`. Put those mapped to the intent name into the alexa.IntenHandlers{}<br>
	[See here, how to use it](#handlers)

- `alexa.LocaleStrings = alexa.Localisation{}`<br>
	The IntentHandlers type is a mapping of language code strings to alexa.Translation functions.<br>
	[See here, how to use it](#locales)

- `alexa.Translation{}`<br>
	A translation of all your strings. It is basicly a map, where the key is string and the value either string or []string.<br>
	[See here, how to use it](#locales) and [see here, how to access those strings in your intent code.](#localization)

- `alexa.BeforeHandler`<br>
	You can assign a function `func(*Context)` to this property. This function is then called before any intent is called.<br>
	If c.Abort() was called within this function, no following intent is called.

- `alexa.MultiHandler(func(c *alexa.Context), func(c *alexa.Context), func(c *alexa.Context))`<br>
	Provide this function instead of providing an intent function directly, if you want to use more than one function on that intent.<br>
	You can add as many IntentHandler functions as you want.

- `alexa.Handle()`
	This is the function the Lambda.Start() function wants to have. Just provide it [as shown here](#main)

# __Context__
Your intent functions are provided with an alexa.Context pointer. That contains all the information you need.<br>

## __Conversational__
This SDK keeps conversation with your customer simple. Basicly you want to ask or tell something. And sometimes you want to provide a card for the Alexa App.<br>
On most functions (if they work with cards) you can immediately add card information.

Conversational functions are:
- `c.Tell("I am Batman")`<br>
	Tells something and ends the session.<br>
	Usable with __card__.<br>

- `c.Ask("Why so serious?")` or `c.Ask("Why so serious?", "didn't you hear me?")`<br>
	Asks the customer and keeps the session open, thus waiting for an answer. That can invoke an other intent then.<br>
	The first argument is mandatory. It is the actual question. The second is reprompt speech. If you provide more than one (after the first), the reprompt speech is randomly chosen out of those.<br>
	Usable with __card__.<br>

- `c.ElicitSlot("mySlot", "What do you want?", "Are you sure?", optionalIntent)`<br>
	Asks the user about a specific slot value. First argument is the slots name, second is the actual speech (question), third is reprompt, if needed. The last one can be a modified \*dialog.EchoIntent, if you need to do so. Otherwise let it be nil.<br>
	Usable with __card__.<br>

- `c.ConfirmSlot("TheSlot", "please confirm", optionalIntent)`<br>
	Let Alexa confirm the slot named in the first parameter by the customer. The second parameter is the speech rendered by Alexa and the third one can be a modified \*dialog.EchoIntent, if you need to do so. Otherwise let it be nil.<br>
	Usable with __card__.<br>

- `c.ConfirmIntent("Confirm all the stuff please", optionalIntent)`<br>
	Let Alexa confirm all information given by the customer. The second parameter can be a modified \*dialog.EchoIntent, if you need to do so. Otherwise let it be nil.<br>
	Usable with __card__.<br>

- `c.Delegate(optionalIntent)`<br>
	Delegate slot fulfillment to Alexa. You can check the dialog state `c.DialogState()` for e.g. 'STARTED' or 'COMPLETED' to know, whether all the slots are filled or you have to again Delegate.<br>
	You can provide a \*dialog.EchoIntent to send a modified intent, if you need to do so. Otherwise let it be nil.<br>

### __With card__
Just add one of these method calls to the conversation chain:<br>
- `SimpleCard("My Title", "An explaining text")`<br>
	Adds a simple card with a title and some text.<br>

- `StandardCard("My Title", "An explaining text", "https://url/to/small/image.png", "https://url/to/large/image.png")`<br>
	Adds a card with a title, some text, a small and a large picture.<br>

- `LinkAccountCard()`<br>
	Adds a link account card, remembering the customer to do some account linking first.<br>

Examples:<br>
`c.Tell("I am Batman").SimpleCard("Dark Knight said", "He is Batman")`<br>
`c.Tell("Please link your Gotham account first").LinkAccountCard()`<br>

## __Progressive Request__
Before you actually respond to the intent request, you can send progressive requests to keep your customer entertained while your response may take longer.<br>
The SDK makes sure the request was responded before another one may be sent or before your intent response is sent.<br>

- `c.Progress("Sorry, this may take a while")`<br>
	Sends a progress to Alexa to be rendered. You can use ssml.<br>
	This method returns immediately as the request runs parallel to the subsequent code.<br>

## __Localization__
As you could already notice, localization is rather easy with this SDK. It always automaticly chooses the current language and returns the translation.<br>
There are three different functions provided:
- `c.T("FOO")`<br>
	Returns the translated string for the key "FOO"<br>
	If FOO is an array of strings, a random string is chosen.
	If more than one key is provided, the values are concatenated.

- `c.TA("BAR")`<br>
	Returns the translated array for the key "BAR"<br>
	If BAR is just a string, it returns an array with this one entry.

- `c.TR("BAZ", &MyStruct)` or `c.TR("BAZ", alexa.R{"foo":"bar"})`<br>
	Returns the translated string for the key "BAZ" and substitutes variables with the given struct. You can also provide an `alexa.R` for spontanious values. For preparing the struct read on.

### __Runtime values for localized strings__
You can place values into your translated strings on runtime. You use the `c.TR()` function for this and as a second argument you provide a `alexa.R` or any struct.<br>
Place variables to be substituted within your string as follows:<br>
`"BAZ": "You have {NUM} cookies"`<br>
The key would be "NUM" then. The translation kit either uses the struct field name or if provided an alexa tag. Example:
```go
type MyStruct struct {
	Num int `json:"num" alexa:"NUM"`
}
```
The tag always wins over the field name.<br>
Please note that cascaded structs go concatenated with dots as a value. A struct 'foo' contains another struct 'bar' which contains a value "baz". The replacement variable would then be 'bar.baz'.

**Please note that float numbers are rendered according to language.**

## __Attributes__
Supported types: interface{}, string, bool, int, float32.<br>
To access the session attributes, you simply use the Attr method of the context.
```go
// to write attributes
c.Attr("my_number", 42)
c.Attr("my_str", "what a value")
c.Attr("my_flag", true)

// to read attributes
num := c.Attr("my_number").Int()
str := c.Attr("my_str").String()
flag := c.Attr("my_flag").Bool()

// delete an attribute
c.Attr("old_val", nil)
```
__Please note that the SDK tries to cast values except for bool!__.<br>

<br>
The Attr method returns an Attr object which not only provides you with type cast methods, but also an Exists Method. Therefore you can do the following:
```go
if a := c.Attr("page"); a.Exists() {
	c.Tell(c.TR("YOUR_PAGE", a.R()))
}
c.Tell(c.T("NO_PAGE"))
```
As you can see here, the Attr object also provides a method `R()` wich can directly be used as an input to `TR()` of localisation, as it returns a suitable map.<br>
You can call the `R()` method with or without parameters. If you call it without parameters, the name of the key is the attributes name. If you provide parameters, these will be used as key names and thus will be replaced with the same value (can be useful).

## __Slots__
For slots the alexa.Context provides the Method `Slot("slotname")`. It returns an object providing you with three values of the slot. If the slot does not exist in the request, those three values are empty, but never does the Slot method return nil. The three values of the slot here are:
- __ID__<br>
	The ID you can assign to a slot.
- __Value__<br>
	The value you can assign to a slot.
- __Spoken__<br>
	The words actually spoken by the customer.
- __ConfirmationStatus__<br>
	The status of the confirmation of this slot, if in a dialog
- __Match__<br>
	Is true if the actual spoken words are matching a value of this slot or its synonyms

## __Other methods__
- `c.NewSession()`<br>
	Returns true if the session is just started and false otherwise.

- `c.SessionID()`<br>
	Returns the current sessions id.

- `c.Locale()`<br>
	Returns the current used language code. Possible values are de-DE, en-GB, en-US and more. Please see Amazon docs.

- `c.DialogState()`<br>
	Returns the current dialog state.

## __Properties__
- `c.Intent`<br>
	The name of the currently called intent.

- `c.System`<br>
	Returns the system property of the request. That contains information about the calling system and user.

