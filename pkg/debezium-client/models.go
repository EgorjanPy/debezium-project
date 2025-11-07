package debeziumclient

type GetConnectorResponse struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
	Tasks  []TaskInfo             `json:"tasks"`
	Type   string                 `json:"type"`
}

type GetConnectorsStatusResponse struct {
	Connectors []struct {
		Status GetConnectorStatusResponse
	} `json:"connector"`
}

type CreateConnectorConfig struct {
	ConnectorClass       string            `json:"connector.class"`
	TasksMax             string            `json:"tasks.max"`
	DatabaseHostname     string            `json:"database.hostname"`
	DatabasePort         string            `json:"database.port"`
	DatabaseUser         string            `json:"database.user"`
	DatabasePassword     string            `json:"database.password"`
	DatabaseDbname       string            `json:"database.dbname"`
	DatabaseServerName   string            `json:"database.server.name"`
	AdditionalParameters map[string]string `json:"-,omitempty"`
}

type TaskInfo struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
}
type GetConnectorStatusResponse struct {
	Name      string `json:"name"`
	Connector struct {
		State    string `json:"state"`
		WorkerID string `json:"worker_id"`
	} `json:"connector"`
	Tasks []TaskInfo `json:"tasks"`
	Type  string     `json:"type"`
}

type CreateConnectorRequest struct {
	Name   string                `json:"name"`
	Config CreateConnectorConfig `json:"config"`
}
type CreateConnectorErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
type CreateConnectorResponse struct {
	Name   string                `json:"name"`
	Config CreateConnectorConfig `json:"config"`
	Tasks  []any                 `json:"tasks"`
	Type   string                `json:"type"`
}
type DeleteConnectorError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
