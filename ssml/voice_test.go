package ssml_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/dasjott/alexa-sdk-go/ssml"
)

func TestVoice(t *testing.T) {
	test := assert.New(t)

	hans := ssml.NewVoice("Hans")
	test.Equal("<voice name=\"Hans\">pretty cool</voice>", hans("pretty cool"))
	test.Equal("<voice name=\"Nicole\">also nice</voice>", ssml.Voice("Nicole", "also nice"))
}
