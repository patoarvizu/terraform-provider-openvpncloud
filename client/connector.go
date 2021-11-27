package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Connector struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name"`
	NetworkItemId   string `json:"networkItemId"`
	NetworkItemType string `json:"networkItemType"`
	VpnRegionId     string `json:"vpnRegionId"`
	IPv4Address     string `json:"ipV4Address"`
	IPv6Address     string `json:"ipV6Address"`
}

const (
	NetworkItemTypeHost    = "HOST"
	NetworkItemTypeNetwork = "NETWORK"
)

func (c *Client) GetConnector(name string) (*Connector, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/connectors", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var connectors []Connector
	err = json.Unmarshal(body, &connectors)
	if err != nil {
		return nil, err
	}
	for _, c := range connectors {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, nil
}

func (c *Client) GetConnectorsForNetwork(networkId string) ([]Connector, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/connectors", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var connectors []Connector
	err = json.Unmarshal(body, &connectors)
	if err != nil {
		return nil, err
	}
	var networkConnectors []Connector
	for _, v := range connectors {
		if v.NetworkItemId == networkId {
			networkConnectors = append(networkConnectors, v)
		}
	}
	return networkConnectors, nil
}

func (c *Client) UpdateNetworkConnectors(network Network) error {
	for _, conn := range network.Connectors {
		if len(conn.Id) == 0 {
			connJson, err := json.Marshal(conn)
			if err != nil {
				return err
			}
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/connectors?networkItemId=%s&networkItemType=NETWORK", c.BaseURL, network.Id), bytes.NewBuffer(connJson))
			if err != nil {
				return err
			}
			_, err = c.DoRequest(req)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) AddNetworkConnector(connector Connector, networkId string) (*Connector, error) {
	connectorJson, err := json.Marshal(connector)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/connectors?networkItemId=%s&networkItemType=NETWORK", c.BaseURL, networkId), bytes.NewBuffer(connectorJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	var conn Connector
	err = json.Unmarshal(body, &conn)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (c *Client) RemoveNetworkConnector(connectorId string, networkId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/connectors/%s?networkItemId=%s&networkItemType=NETWORK", c.BaseURL, connectorId, networkId), nil)
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	return err
}
