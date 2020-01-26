package api

// constants for requests
const (
	// URLAddress is the URL to retrieve the address for current device
	URLAddress = "/v1/devices/{deviceId}/settings/address"
	// URLRegionAndZIP is the URL to retrieve region and zip for current device
	URLRegionAndZIP = "/v1/devices/{deviceId}/settings/address/countryAndPostalCode"

	// URLTimeZone is the URL to retrieve the timezone for current device
	URLTimeZone = "/v2/devices/{deviceId}/settings/System.timeZone"
	// URLDistanceUnits is the URL to retrieve the distance units for current device
	URLDistanceUnits = "/v2/devices/{deviceId}/settings/System.distanceUnits"
	// URLTemperatureUnit is the URL to retrieve temperature unit for current device
	URLTemperatureUnit = "/v2/devices/{deviceId}/settings/System.temperatureUnit"

	// URLFullName is the URL to retrieve full name of current devices user
	URLFullName = "/v2/accounts/~current/settings/Profile.name"
	// URLGivenName is the URL to retrieve given name of current devices user
	URLGivenName = "/v2/accounts/~current/settings/Profile.givenName"
	// URLEmailAddress is the URL to retrieve email address of current devices user
	URLEmailAddress = "/v2/accounts/~current/settings/Profile.email"
	// URLPhoneNumber is the URL to retrieve the phone number of current devices user
	URLPhoneNumber = "/v2/accounts/~current/settings/Profile.mobileNumber"
)

// constants for responses
const (
	// Metric is a possible response to GetDistanceUnits request
	Metric = "METRIC"
	// Imperial is a possible response to GetDistanceUnits request
	Imperial = "IMPERIAL"

	// Celsius is a possible response to GetTemperatureUnit request
	Celsius = "CELSIUS"
	// Fahrenheit is a possible response to GetTemperatureUnit request
	Fahrenheit = "FAHRENHEIT"
)

// PhoneNumber contains contry code and actual phone number for current devices user
type PhoneNumber struct {
	CountryCode string `json:"countryCode"`
	PhoneNumber string `json:"phoneNumber"`
}

// RegionAndZip contains country code and zip information for current device
type RegionAndZip struct {
	CountryCode string `json:"countryCode"`
	PostalCode  string `json:"postalCode"`
}

// Address contains address information of current device
type Address struct {
	StateOrRegion    string `json:"stateOrRegion"`
	City             string `json:"city"`
	CountryCode      string `json:"countryCode"`
	PostalCode       string `json:"postalCode"`
	AddressLine1     string `json:"addressLine1"`
	AddressLine2     string `json:"addressLine2"`
	AddressLine3     string `json:"addressLine3"`
	DistrictOrCounty string `json:"districtOrCounty"`
}
