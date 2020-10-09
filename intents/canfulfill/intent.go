package canfulfill

import "github.com/dasjott/alexa-sdk-go/intents"

/*
{
    "version":"1.0",
    "response":{
        "canFulfillIntent": {
            "canFulfill": "YES",
            "slots":{
                "slotName1": {
                    "canUnderstand": "YES",
                    "canFulfill": "YES"
                },
               "slotName2": {
                    "canUnderstand": "YES",
                    "canFulfill": "YES"
                }
            }
        }
    }
}
*/

const (
	NO    = "NO"
	Yes   = "YES"
	MAYBE = "MAYBE"
)

type (
	Answer string

	Slot struct {
		CanUnderstand Answer `json:"canUnderstand"`
		CanFulfill    Answer `json:"canFulfill"`
	}

	Slots map[string]Slot

	Response struct {
		CanFullfill Answer `json:"canFulfill"`
		Slots       Slots  `json:"slots"`
	}
)

func NewXResponse(content *Response) *intents.XResponse {
	return &intents.XResponse{
		Version: "1.0",
		ResponseField: map[string]interface{}{
			"canFulfillIntent": content,
		},
	}
}
