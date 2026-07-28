package configs

import "time"

type CreateConfigRequest struct {
	Name         string `json:"name"`
	EndpointID   uint   `json:"endpoint_id"`
	ConfigTypeID *uint  `json:"config_type_id"`
}

type UpdateConfigRequest struct {
	Name string `json:"name"`
}

type EndpointResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type ConfigTypeResponse struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	ConfigTemplate    string `json:"config_template"`
	UsageInstructions string `json:"usage_instructions"`
	ClientLinks       string `json:"client_links"`
}

type Response struct {
	ID            uint                `json:"id"`
	Name          string              `json:"name"`
	UserID        uint                `json:"user_id"`
	Endpoint      EndpointResponse    `json:"endpoint"`
	ConfigType    *ConfigTypeResponse `json:"config_type,omitempty"`
	ClientIPv4    string              `json:"client_ipv4"`
	ClientIPv6    string              `json:"client_ipv6"`
	PublicKey     string              `json:"public_key"`
	PrivateKey    string              `json:"private_key"`
	PeerPublicKey string              `json:"peer_public_key"`
	ConfigContent string              `json:"config_content"`
	CreatedAt     time.Time           `json:"created_at"`
}

func ToResponse(config Config) Response {
	response := Response{
		ID:            config.ID,
		Name:          config.Name,
		UserID:        config.UserID,
		ClientIPv4:    config.ClientIPv4,
		ClientIPv6:    config.ClientIPv6,
		PublicKey:     config.PublicKey,
		PrivateKey:    config.PrivateKey,
		PeerPublicKey: config.PeerPublicKey,
		ConfigContent: config.ConfigContent,
		Endpoint: EndpointResponse{
			ID:      config.Endpoint.ID,
			Name:    config.Endpoint.Name,
			Address: config.Endpoint.Address,
			Port:    config.Endpoint.Port,
		},
		CreatedAt: config.CreatedAt,
	}

	if config.ConfigType != nil {
		response.ConfigType = &ConfigTypeResponse{
			ID:                config.ConfigType.ID,
			Name:              config.ConfigType.Name,
			ConfigTemplate:    config.ConfigType.ConfigTemplate,
			UsageInstructions: config.ConfigType.UsageInstructions,
			ClientLinks:       config.ConfigType.ClientLinks,
		}
	}

	return response
}

func ToResponseList(configs []Config) []Response {
	response := make([]Response, len(configs))

	for i, config := range configs {
		response[i] = ToResponse(config)
	}

	return response
}
