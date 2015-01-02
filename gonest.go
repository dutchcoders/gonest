/*
The MIT License (MIT)

Copyright (c) 2014 DutchCoders [https://github.com/dutchcoders/]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package gonest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const USER_AGENT = "gonest library v0.1"

type Nest struct {
	clientid string
	token    string
}

type Alarm struct {
}

type TemperatureScale string

type Thermostat struct {
	DeviceId               string           `json:"device_id"`
	Locale                 string           `json:"locale"`
	SoftwareVersion        string           `json:"software_version"`
	Name                   string           `json:"name"`
	NameLong               string           `json:"name_long"`
	LastConnection         string           `json:"last_connection"`
	IsOnline               bool             `json:"is_online"`
	CanCool                bool             `json:"can_cool"`
	CanHeat                bool             `json:"can_heat"`
	IsUsingEmergencyHeat   bool             `json:"is_using_emergency_heat"`
	HasFan                 bool             `json:"has_fan"`
	FanTimerActive         bool             `json:"fan_timer_active"`
	FanTimerTimeout        bool             `json:"fan_timer_timeout"`
	HasLeaf                bool             `json:"has_leaf"`
	TemperatureScale       TemperatureScale `json:"temperature_scale"`
	TargetTemperatureF     float64          `json:"target_temperature_f"`
	TargetTemperatureC     float64          `json:"target_temperature_c"`
	TargetTemperatureHighF float64          `json:"target_temperature_high_f"`
	TargetTemperatureHighC float64          `json:"target_temperature_high_c"`
	TargetTemperatureLowF  float64          `json:"target_temperature_low_f"`
	TargetTemperatureLowC  float64          `json:"target_temperature_low_c"`
	AwayTemperatureHighF   float64          `json:"away_temperature_high_f"`
	AwayTemperatureHighC   float64          `json:"away_temperature_high_c"`
	AwayTemperatureLowF    float64          `json:"away_temperature_low_f"`
	AwayTemperatureLowC    float64          `json:"away_temperature_low_c"`
	HvacMode               string           `json:"hvac_mode"`
	AmbientTemperatureF    float64          `json:"ambient_temperature_f"`
	AmbientTemperatureC    float64          `json:"ambient_temperature_c"`
	Humidity               float64          `json:"humidity"`
}

type Devices struct {
	Thermostats map[string]Thermostat `json:"thermostats"`
	Alarms      map[string]Alarm      `json:"smoke_co_alarms"`
}

type Response struct {
	Devices Devices `json:"devices"`
}

type OAuth2Response struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

type OAuth2ResponseError struct {
	error
	Name        string `json:"error"`
	Description string `json:"error_description"`
}

func (e *OAuth2ResponseError) Error() string {
	return e.Name
}

// authorizes the pincode from the authorization request and returns
// access token
func (nest *Nest) Authorize(secret string, code string) error {
	url := fmt.Sprintf("http://api.home.nest.com/oauth2/access_token?code=%s&client_id=%s&client_secret=%s&grant_type=authorization_code", code, nest.clientid, secret)

	client := &http.Client{}

	var err error
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept", "application/json")

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != 200 {
		var r OAuth2ResponseError
		if err := dec.Decode(&r); err == io.EOF {
		} else if err != nil {
			return err
		}
		return errors.New(r.Description)
	} else {
		var r OAuth2Response
		if err := dec.Decode(&r); err == io.EOF {
		} else if err != nil {
			return err
		}

		nest.token = r.Token
		return nil
	}
}

func (nest *Nest) All(o interface{}) error {
	return nest.get("", o)
}

func (nest *Nest) Devices(o interface{}) error {
	return nest.get("devices", o)
}

// gets response from nest api
func (nest *Nest) get(path string, nr interface{}) error {
	url := fmt.Sprintf("https://developer-api.nest.com/%s?auth=%s", path, nest.token)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept", "application/json")

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(&nr); err == io.EOF {
	} else if err != nil {
		return err
	}

	return nil
}

// connects to nest api and returns nest object
func Connect(clientid string, token string) (*Nest, error) {
	nest := &Nest{clientid: clientid, token: token}

	if token == "" {
		return nest, errors.New(fmt.Sprintf("No authorization token, register at: https://home.nest.com/login/oauth2?client_id=%s&state=STATE", clientid))
	}

	return nest, nil
}
