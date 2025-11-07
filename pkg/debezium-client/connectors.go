package debeziumclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	getConnector          = "/connectors/%s"
	getConnectorStatus    = "/connectors/%s/status"
	getConnectorsStatuses = "/connectors?expand=status"
	postCreateConnectors  = "/connectors"
	deleteConnector       = "/connectors/%s"
)

/*
	TODO:

const (

	updateConnectorConfig = "/connectors/%s/config"
	pauseConnector        = "/connectors/%s/pause"
	resumeConnector       = "/connectors/%s/resume"
	restartConnector      = "/connectors/%s/restart"
	getConnectorTasks     = "/connectors/%s/tasks"
	restartConnectorTask  = "/connectors/%s/tasks/%d/restart"
	listConnectors        = "/connectors"

)
*/
func (c *Client) GetConnector(ctx context.Context, name string) (GetConnectorResponse, error) {
	var getConnectorResponse GetConnectorResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(getConnector, name), nil)
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("GetConnector.NewRequestWithContext: %w", err)
	}
	resp, err := c.cc.Do(req)
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("GetConnector.Do: %w", err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&getConnectorResponse)
	if err != nil {
		return GetConnectorResponse{}, fmt.Errorf("GetConnector.Decode: %w", err)
	}
	return getConnectorResponse, nil

}
func (c *Client) GetConnectorsStatuses(ctx context.Context) (GetConnectorsStatusResponse, error) {
	var connectorResponse GetConnectorsStatusResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+getConnectorsStatuses, nil)
	if err != nil {
		return GetConnectorsStatusResponse{}, fmt.Errorf("GetConnectorsStatuses.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return GetConnectorsStatusResponse{}, fmt.Errorf("GetConnectorsStatuses.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&connectorResponse); err != nil {
		return GetConnectorsStatusResponse{}, fmt.Errorf("GetConnectorsStatuses.DecodeJSON: %w", err)
	}
	return connectorResponse, nil
}
func (c *Client) PostCreateConnectors(ctx context.Context, request CreateConnectorRequest) (bool, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return false, fmt.Errorf("PostCreateConnectors.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+postCreateConnectors, bytes.NewBuffer(data))
	if err != nil {
		return false, fmt.Errorf("PostCreateConnectors.NewRequestWithContext: %w", err)
	}

	resp, err := c.cc.Do(req)
	if err != nil {
		return false, fmt.Errorf("PostCreateConnectors.Client.Do: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errResponse CreateConnectorErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return false, fmt.Errorf("PostCreateConnectors.UnmarshalJSON: %s", errResponse.Message)
		}
	}
	var connectorResponse CreateConnectorRequest
	if err := json.NewDecoder(resp.Body).Decode(&connectorResponse); err != nil {
		return false, fmt.Errorf("PostCreateConnectors.UnmarshalJSON: %w", err)
	}
	defer resp.Body.Close()

	return true, nil
}
func (c *Client) GetConnectorStatusByName(ctx context.Context, connectorName string) (GetConnectorStatusResponse, error) {
	var response GetConnectorStatusResponse
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+fmt.Sprintf(getConnectorStatus, connectorName),
		nil,
	)
	if err != nil {
		return GetConnectorStatusResponse{}, fmt.Errorf("GetConnectorStatusByName.NewRequest: %w", err)
	}
	resp, err := c.cc.Do(req)
	if err != nil {
		return GetConnectorStatusResponse{}, fmt.Errorf("GetConnectorStatusByName.Do: %w", err)
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return GetConnectorStatusResponse{}, fmt.Errorf("GetConnectorStatusByName.DecodeJSON: %w", err)
	}
	defer resp.Body.Close()

	return response, nil
}
func (c *Client) DeleteConnector(ctx context.Context, connectorName string) (bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		c.baseURL+fmt.Sprintf(deleteConnector, connectorName),
		nil,
	)
	if err != nil {
		return false, fmt.Errorf("DeleteConnector.NewRequestWithContext: %w", err)
	}
	resp, err := c.cc.Do(req)
	if err != nil {
		return false, fmt.Errorf("DeleteConnector.Do: %w", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		var errResponse DeleteConnectorError
		err := json.NewDecoder(resp.Body).Decode(&errResponse)
		defer resp.Body.Close()
		if err != nil {
			return false, fmt.Errorf("DeleteConnector.Decode: %s", errResponse.Message)
		}
	}
	return true, nil
}
