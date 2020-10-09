package dialog

import (
	"strings"
	"time"
)

type EchoRequest struct {
	Version string              `json:"version"`
	Session *EchoSession        `json:"session"`
	Request *EchoRequestBody    `json:"request"`
	Context *EchoRequestContext `json:"context"`
}

type EchoSession struct {
	New         bool                   `json:"new"`
	SessionID   string                 `json:"sessionId"`
	Application *EchoApplication       `json:"application"`
	Attributes  map[string]interface{} `json:"attributes"`
	User        *EchoUser              `json:"user"`
}

type EchoRequestBody struct {
	Type        string      `json:"type"`
	RequestID   string      `json:"requestId"`
	Timestamp   string      `json:"timestamp"`
	DialogState string      `json:"dialogState"`
	Intent      *EchoIntent `json:"intent"`
	Reason      string      `json:"reason"`
	Locale      string      `json:"locale"`
}

type EchoApplication struct {
	ID string `json:"applicationId"`
}

type EchoUser struct {
	ID          string            `json:"userId"`
	AccessToken string            `json:"accessToken"`
	Permissions map[string]string `json:"permissions"`
}

type EchoPerson struct {
	ID          string `json:"personId"`
	AccessToken string `json:"accessToken"`
}

type EchoDevice struct {
	ID                  string                 `json:"deviceId"`
	SupportedInterfaces map[string]interface{} `json:"supportedInterfaces"`
}

type EchoSystem struct {
	Device         *EchoDevice      `json:"device"`
	Application    *EchoApplication `json:"application"`
	User           *EchoUser        `json:"user"`
	Person         *EchoPerson      `json:"person"`
	APIEndpoint    string           `json:"apiEndpoint"`
	APIAccessToken string           `json:"apiAccessToken"`
}

type EchoRequestContext struct {
	System *EchoSystem `json:"System"`
}

// EchoSlot is the json part for a slot
type EchoSlot struct {
	Name               string `json:"name"`
	Value              string `json:"value"`
	ConfirmationStatus string `json:"confirmationStatus,omitempty"`
	Resolutions        *struct {
		ResolutionsPerAuthority []EchoAuthorityResolution `json:"resolutionsPerAuthority"`
	} `json:"resolutions"`
}

type EchoAuthorityResolution struct {
	Authority string `json:"authority"`
	Status    struct {
		Code string `json:"code"`
	} `json:"code"`
	Values []EchoAuthorityResolutionValue `json:"values"`
}

type EchoAuthorityResolutionValue struct {
	Value *NameID `json:"value"`
}

type NameID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FUNCTIONS //////////////////////////////////////////////////////////////////
func (er *EchoRequest) VerifyTimestamp() bool {
	return time.Since(er.GetTime()) < time.Duration(150)*time.Second
}

func (er *EchoRequest) VerifyAppID(appID string) bool {
	return (er.Context.System.Application.ID == appID)
}

func (er *EchoRequest) GetSessionID() string {
	return er.Session.SessionID
}

func (er *EchoRequest) GetUserID() string {
	return er.Context.System.User.ID
}

func (er *EchoRequest) GetRequestType() string {
	return er.Request.Type
}

func (er *EchoRequest) GetIntentName() string {
	if er.GetRequestType() == "IntentRequest" {
		return er.Request.Intent.Name
	}
	return er.GetRequestType()
}

func (er *EchoRequest) GetTime() time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z", er.Request.Timestamp)
	return t
}

func (res *EchoAuthorityResolution) IsBuiltIn() bool {
	parts := strings.SplitN(res.Authority, ".", 6)
	return len(parts) > 4 && parts[4] == "AMAZON"
}

func (res *EchoAuthorityResolution) IsMatch() bool {
	return res.Status.Code == "ER_SUCCESS_MATCH"
	// return (res.Status.Code == "ER_SUCCESS_MATCH") || (res.IsBuiltIn() && res.Status.Code == "ER_SUCCESS_NO_MATCH")
}
