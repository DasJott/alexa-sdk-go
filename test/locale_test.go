package test_test

import (
	"testing"

	"github.com/dasjott/alexa-sdk-go"
	"github.com/dasjott/alexa-sdk-go/dialog"
	"github.com/stretchr/testify/assert"
)

func TestTranslator(t *testing.T) {
	test := assert.New(t)

	l := alexa.Localisation{
		"de-DE": alexa.Translation{
			"one":      "eins",
			"two":      "zwei",
			"three":    "drei",
			"hello":    "Hallo {person}",
			"float":    "Zahl {float}",
			"int":      "Zahl {int}",
			"sentence": "Ich bin {person}, bin {age} Jahre alt und {size}m groß.",
		},
	}

	trans := l.GetTranslator("de-DE")
	test.NotNil(trans)

	test.Equal("eins", trans.GetString("one"))
	test.Equal("zwei", trans.GetString("two"))
	test.Equal("drei", trans.GetString("three"))

	{
		str := trans.GetStringAndReplace("hello", alexa.R{"person": "Batman"})
		test.Equal("Hallo Batman", str)
	}
	{
		str := trans.GetStringAndReplace("int", alexa.R{"int": 1939})
		test.Equal("Zahl 1.939", str)
	}
	{
		str := trans.GetStringAndReplace("float", alexa.R{"float": 19.39})
		test.Equal("Zahl 19,39", str)
	}

	type Vals struct {
		Name string  `alexa:"person"`
		Age  int     `alexa:"age"`
		Size float32 `alexa:"size,1"`
	}
	myData := Vals{"Batman", 39, 17.5}

	str := trans.GetStringWithVariables("sentence", &myData)
	test.Equal("Ich bin Batman, bin 39 Jahre alt und 17,5m groß.", str)
}

func TestVoice(t *testing.T) {
	test := assert.New(t)

	resp := dialog.NewResponse()

	resp.OutputSSML("I am Batman")
	test.Equal("SSML", resp.Response.OutputSpeech.Type)
	test.Equal("<speak>I am Batman</speak>", resp.Response.OutputSpeech.SSML)

	dialog.SetVoice("Alfred")
	resp.OutputSSML("take care, master bruce")
	test.Equal("<speak><voice name=\"Alfred\">take care, master bruce</voice></speak>", resp.Response.OutputSpeech.SSML)
}
