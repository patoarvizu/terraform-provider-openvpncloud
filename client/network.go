package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Network struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Egress         bool        `json:"egress"`
	InternetAccess string      `json:"internetAcess"`
	SystemSubnets  []string    `json:"systemSubnets"`
	Routes         []Route     `json:"routes"`
	Connectors     []Connector `json:"connectors"`
}

const (
	InternetAccessBlocked        = "BLOCKED"
	InternetAccessGlobalInternet = "GLOBAL_INTERNET"
	InternetAccessLocal          = "LOCAL"
)

func (c *Client) GetNetworkByName(name string) (*Network, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/networks", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var networks []Network
	err = json.Unmarshal(body, &networks)
	if err != nil {
		return nil, err
	}
	for _, n := range networks {
		if n.Name == name {
			return &n, nil
		}
	}
	return nil, nil
}

func (c *Client) CreateNetwork(network Network) (*Network, error) {
	networkJson, err := json.Marshal(network)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/networks", c.BaseURL), bytes.NewBuffer(networkJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var n Network
	err = json.Unmarshal(body, &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
