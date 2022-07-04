package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Host struct {
	Id             string      `json:"id,omitempty"`
	Name           string      `json:"name"`
	InternetAccess string      `json:"internetAccess"`
	SystemSubnets  []string    `json:"systemSubnets"`
	Connectors     []Connector `json:"connectors"`
}

func (c *Client) GetHosts() ([]Host, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/hosts", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var hosts []Host
	err = json.Unmarshal(body, &hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func (c *Client) GetHostByName(name string) (*Host, error) {
	hosts, err := c.GetHosts()
	if err != nil {
		return nil, err
	}
	for _, h := range hosts {
		if h.Name == name {
			return &h, nil
		}
	}
	return nil, nil
}
