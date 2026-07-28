package groups

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	EndpointIDs []uint `json:"endpoint_ids"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	EndpointIDs []uint `json:"endpoint_ids"`
}

type Response struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EndpointIDs []uint `json:"endpoint_ids"`
	CreatedAt   string `json:"created_at"`
}

func ToResponse(group Group) Response {
	endpointIDs := make([]uint, 0, len(group.Endpoints))
	for _, ep := range group.Endpoints {
		endpointIDs = append(endpointIDs, ep.ID)
	}
	return Response{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		EndpointIDs: endpointIDs,
		CreatedAt:   group.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToResponseList(groups []Group) []Response {
	response := make([]Response, len(groups))

	for i, group := range groups {
		response[i] = ToResponse(group)
	}

	return response
}
