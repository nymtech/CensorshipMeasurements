// Package example contains a nym experiment.
//
package nym

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ooni/probe-cli/v3/internal/model"
)

const (
	testName    = "nym"
	testVersion = "0.1.0"
)

// Gateway description is the struct that is sent by the validator api
type GatewayDescription struct {
	BlockHeight int `json:"block_height"`
	Gateway     struct {
		ClientsPort int    `json:"clients_port"`
		Host        string `json:"host"`
		IdentityKey string `json:"identity_key"`
		Location    string `json:"location"`
		MixPort     int    `json:"mix_port"`
		SphinxKey   string `json:"sphinx_key"`
		Version     string `json:"version"`
	} `json:"gateway"`
	Owner        string `json:"owner"`
	PledgeAmount struct {
		Amount string `json:"amount"`
		Denom  string `json:"denom"`
	} `json:"pledge_amount"`
	Proxy struct{} `json:"proxy"`
}

// Config contains the experiment config.
//
// This contains all the settings that user can set to modify the behaviour
// of this experiment. By tagging these variables with `ooni:"..."`, we allow
// miniooni's -O flag to find them and set them.
type Config struct {
	NymValidatorURL string `json:"nym_validator_url"`
}

// TestKeys contains the experiment's result.
//
// This is what will end up into the Measurement.TestKeys field
// when you run this experiment.
//
// In other words, the variables in this struct will be
// the specific results of this experiment.
type TestKeys struct {
	ValidatorAPIReachable       bool  `json:"validator_api_reachable"`
	ValidatorAPIGettingGateways bool  `json:"validator_api_gettingGateways"`
	GatewaysTotal               int64 `json:"gateways_total"`
	GatewaysAccessible          int64 `json:"gateways_accessible"`
}

// Measurer performs the measurement.
type Measurer struct {
	config Config
}

// ExperimentName implements model.ExperimentMeasurer.ExperimentName.
func (m Measurer) ExperimentName() string {
	return testName
}

// ExperimentVersion implements model.ExperimentMeasurer.ExperimentVersion.
func (m Measurer) ExperimentVersion() string {
	return testVersion
}

// Run implements model.ExperimentMeasurer.Run.
func (m Measurer) Run(ctx context.Context, args *model.ExperimentArgs) error {
	//callbacks := args.Callbacks
	measurement := args.Measurement
	sess := args.Session
	testkeys := &TestKeys{ValidatorAPIReachable: false, ValidatorAPIGettingGateways: false, GatewaysTotal: 0, GatewaysAccessible: 0}
	measurement.TestKeys = testkeys
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// We start by parsing the input URL. If we cannot parse it, of
	// course this is a hard error and we cannot continue.
	apiURL := m.config.NymValidatorURL + "/api/v1/gateways"
	sess.Logger().Infof("Using API %s", apiURL)
	parsedURL, err := url.Parse(apiURL)
	_ = parsedURL
	if err != nil {
		return err
	}

	clnt := &http.Client{}
	_ = clnt

	//resp, err := m.HTTPClientGET(ctx, clnt, parsedURL)
	resp, err := http.Get(apiURL)
	if err != nil {
		testkeys.ValidatorAPIReachable = false
		return nil
	}
	testkeys.ValidatorAPIReachable = true

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		testkeys.ValidatorAPIGettingGateways = false
		return nil
	}

	var gateways []GatewayDescription
	err = json.Unmarshal(body, &gateways)
	if err != nil {
		testkeys.ValidatorAPIGettingGateways = false
		return nil
	}
	testkeys.ValidatorAPIGettingGateways = true

	for _, g := range gateways {
		testkeys.GatewaysTotal++
		u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", g.Gateway.Host, g.Gateway.ClientsPort), Path: "/"}
		sess.Logger().Infof("Websocket connecting to %s", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			sess.Logger().Warnf("Websocket handshake error: %s", err)
			continue
		}
		defer c.Close()

		// Cleanly close the connection by sending a close message
		err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			sess.Logger().Warnf("Websocket closing errior: %s", err)
			continue
		}
		testkeys.GatewaysAccessible++
	}

	return nil
}

// NewExperimentMeasurer creates a new ExperimentMeasurer.
func NewExperimentMeasurer(config Config) model.ExperimentMeasurer {
	return Measurer{config: config}
}

// SummaryKeys contains summary keys for this experiment.
//
// Note that this structure is part of the ABI contract with ooniprobe
// therefore we should be careful when changing it.
type SummaryKeys struct {
	IsAnomaly bool `json:"-"`
}

// GetSummaryKeys implements model.ExperimentMeasurer.GetSummaryKeys.
func (m Measurer) GetSummaryKeys(measurement *model.Measurement) (interface{}, error) {
	sk := SummaryKeys{IsAnomaly: false}
	tk, ok := measurement.TestKeys.(*TestKeys)
	if !ok {
		return sk, errors.New("invalid test keys type")
	}
	sk.IsAnomaly = !tk.ValidatorAPIReachable || !tk.ValidatorAPIGettingGateways || tk.GatewaysAccessible == 0
	return sk, nil
}
