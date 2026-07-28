package endpoints

type CreateEndpointRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`

	GroupIDs      []uint `json:"group_ids"`
	ConfigTypeIDs []uint `json:"config_type_ids"`
}

type UpdateEndpointRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`

	GroupIDs      []uint `json:"group_ids"`
	ConfigTypeIDs []uint `json:"config_type_ids"`
}

type Response struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func ToResponse(endpoint Endpoint) Response {
	return Response{
		ID:      endpoint.ID,
		Name:    endpoint.Name,
		Address: endpoint.Address,
		Port:    endpoint.Port,
	}
}

func ToResponseList(endpoints []Endpoint) []Response {
	response := make([]Response, len(endpoints))

	for i, endpoint := range endpoints {
		response[i] = ToResponse(endpoint)
	}

	return response
}

type AddGroupEndpointRequest struct {
	GroupID    uint `json:"group_id"`
	EndpointID uint `json:"endpoint_id"`
}

type GroupEndpointResponse struct {
	GroupID    uint `json:"group_id"`
	EndpointID uint `json:"endpoint_id"`
}

func ToGroupEndpointResponse(item GroupEndpoint) GroupEndpointResponse {
	return GroupEndpointResponse{
		GroupID:    item.GroupID,
		EndpointID: item.EndpointID,
	}
}

func ToGroupEndpointResponseList(items []GroupEndpoint) []GroupEndpointResponse {
	result := make([]GroupEndpointResponse, len(items))

	for i, item := range items {
		result[i] = ToGroupEndpointResponse(item)
	}

	return result
}
