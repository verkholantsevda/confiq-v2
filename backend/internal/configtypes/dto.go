package configtypes

type CreateConfigTypeRequest struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	ConfigTemplate    string `json:"config_template"`
	UsageInstructions string `json:"usage_instructions"`
	ClientLinks       string `json:"client_links"`
	IsActive          bool   `json:"is_active"`
}

type UpdateConfigTypeRequest struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	ConfigTemplate    string `json:"config_template"`
	UsageInstructions string `json:"usage_instructions"`
	ClientLinks       string `json:"client_links"`
	IsActive          bool   `json:"is_active"`
}

type Response struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ConfigTemplate    string `json:"config_template"`
	UsageInstructions string `json:"usage_instructions"`
	ClientLinks       string `json:"client_links"`
	IsActive          bool   `json:"is_active"`
}

// Короткий DTO для выбора типа конфигурации
type ShortResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToResponse(configType ConfigType) Response {
	return Response{
		ID:                configType.ID,
		Name:              configType.Name,
		Description:       configType.Description,
		ConfigTemplate:    configType.ConfigTemplate,
		UsageInstructions: configType.UsageInstructions,
		ClientLinks:       configType.ClientLinks,
		IsActive:          configType.IsActive,
	}
}

func ToResponseList(configTypes []ConfigType) []Response {
	response := make([]Response, len(configTypes))

	for i, configType := range configTypes {
		response[i] = ToResponse(configType)
	}

	return response
}

func ToShortResponseList(configTypes []ConfigType) []ShortResponse {
	response := make([]ShortResponse, len(configTypes))

	for i, configType := range configTypes {
		response[i] = ShortResponse{
			ID:   configType.ID,
			Name: configType.Name,
		}
	}

	return response
}
