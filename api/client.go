package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dasjott/alexa-sdk-go/dialog"
	"github.com/dasjott/alexa-sdk-go/test"
)

const (
	hdrAmznRequestID = "X-Amzn-RequestId"
)

// Client is the client to use the alexa api
// get an instance by using NewClient
type Client struct {
	sys *dialog.EchoSystem
}

// NewClient creates an instance of Client with given setup
func NewClient(esys *dialog.EchoSystem) *Client {
	return &Client{
		sys: esys,
	}
}

// Request to be called with string containing {deviceId}
// Please check the constants from this package
func (c *Client) Request(path string) (string, error) {
	url := c.GetDevicePath(path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", "Bearer "+c.sys.APIAccessToken)

	if err == nil {
		var resp *http.Response

		if test.RequestHandler != nil {
			resp = test.RequestHandler(req)
		} else {
			client := http.Client{}
			resp, err = client.Do(req)
		}

		if err == nil && resp != nil {
			if resp.StatusCode != 200 {
				err = fmt.Errorf("response code %d", resp.StatusCode)
			} else {
				buf := bytes.Buffer{}
				_, err = buf.ReadFrom(resp.Body)
				if err == nil {
					return buf.String(), nil
				}
			}
			err = fmt.Errorf("%s x-amzn-requestid: %s", err.Error(), resp.Header.Get(hdrAmznRequestID))
		}
	}
	return "", err
}

// GetPhoneNumber requests phone number of the current devices user
// this is a shortcut for Request(URLPhoneNumber) responding a decent struct
func (c *Client) GetPhoneNumber() (*PhoneNumber, error) {
	data, err := c.Request(URLPhoneNumber)
	if err == nil {
		var phn PhoneNumber
		err = json.Unmarshal([]byte(data), &phn)
		if err == nil {
			return &phn, nil
		}
	}
	return nil, err
}

// GetRegionAndZip requests region and zip of the current device
// this is a shortcut for Request(URLRegionAndZIP) responding a decent struct
func (c *Client) GetRegionAndZip() (*RegionAndZip, error) {
	data, err := c.Request(URLRegionAndZIP)
	if err == nil {
		var raz RegionAndZip
		err = json.Unmarshal([]byte(data), &raz)
		if err == nil {
			return &raz, nil
		}
	}
	return nil, err
}

// GetAddress requests address of the current device
// this is a shortcut for Request(URLAddress) responding a decent struct
func (c *Client) GetAddress() (*Address, error) {
	data, err := c.Request(URLAddress)
	if err == nil {
		var add Address
		err = json.Unmarshal([]byte(data), &add)
		if err == nil {
			return &add, nil
		}
	}
	return nil, err
}

// GetDevicePath substitutes {deviceId} within path with current device id and prepends the current api url
func (c *Client) GetDevicePath(path string) string {
	return c.sys.APIEndpoint + "/" + strings.ReplaceAll(path, "{deviceId}", c.sys.Device.ID)
}
