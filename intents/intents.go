package intents

type (
	Response  interface{}
	XResponse struct {
		Version       string                 `json:"version"`
		ResponseField map[string]interface{} `json:"response"`
	}
)
