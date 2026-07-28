package endpoints

import (
	"confiq/internal/configtypes"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(endpoint *Endpoint) error {
	return r.db.Create(endpoint).Error
}

func (r *Repository) GetByID(id uint) (*Endpoint, error) {
	var endpoint Endpoint

	if err := r.db.
		Preload("ConfigTypes").
		First(&endpoint, id).Error; err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (r *Repository) List() ([]Endpoint, error) {
	var endpoints []Endpoint

	if err := r.db.Order("id").Find(&endpoints).Error; err != nil {
		return nil, err
	}

	return endpoints, nil
}

func (r *Repository) Update(endpoint *Endpoint) error {
	return r.db.Save(endpoint).Error
}

func (r *Repository) ReplaceGroups(endpointID uint, groupIDs []uint) error {
	var endpoint Endpoint
	if err := r.db.First(&endpoint, endpointID).Error; err != nil {
		return err
	}

	if len(groupIDs) == 0 {
		return r.db.
			Where("endpoint_id = ?", endpointID).
			Delete(&GroupEndpoint{}).Error
	}

	// Update the join table directly to avoid importing the groups package.
	if err := r.db.Where("endpoint_id = ?", endpointID).Delete(&GroupEndpoint{}).Error; err != nil {
		return err
	}

	for _, id := range groupIDs {
		if err := r.db.Create(&GroupEndpoint{
			GroupID:    id,
			EndpointID: endpointID,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ReplaceConfigTypes(endpointID uint, configTypeIDs []uint) error {
	var endpoint Endpoint
	if err := r.db.First(&endpoint, endpointID).Error; err != nil {
		return err
	}

	if len(configTypeIDs) == 0 {
		return r.db.Model(&endpoint).Association("ConfigTypes").Clear()
	}

	types := make([]configtypes.ConfigType, 0, len(configTypeIDs))
	for _, id := range configTypeIDs {
		types = append(types, configtypes.ConfigType{ID: id})
	}

	return r.db.Model(&endpoint).Association("ConfigTypes").Replace(types)
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Endpoint{}, id).Error
}

func (r *Repository) AddGroupEndpoint(groupID, endpointID uint) error {
	groupEndpoint := GroupEndpoint{
		GroupID:    groupID,
		EndpointID: endpointID,
	}

	return r.db.Create(&groupEndpoint).Error
}

func (r *Repository) RemoveGroupEndpoint(groupID, endpointID uint) error {
	return r.db.Delete(&GroupEndpoint{}, "group_id = ? AND endpoint_id = ?", groupID, endpointID).Error
}

func (r *Repository) ListGroupsEndpoints() ([]GroupEndpoint, error) {
	var groupEndpoints []GroupEndpoint

	if err := r.db.Order("group_id, endpoint_id").Find(&groupEndpoints).Error; err != nil {
		return nil, err
	}

	return groupEndpoints, nil
}

func (r *Repository) ListForUser(userID uint) ([]Endpoint, error) {
	var endpoints []Endpoint

	err := r.db.
		Table("endpoints").
		Joins("JOIN group_endpoints ge ON ge.endpoint_id = endpoints.id").
		Joins("JOIN users u ON u.group_id = ge.group_id").
		Where("u.id = ?", userID).
		Order("endpoints.id").
		Find(&endpoints).Error

	return endpoints, err
}

func (r *Repository) ListConfigTypes(endpointID uint) ([]configtypes.ConfigType, error) {
	var endpoint Endpoint

	err := r.db.
		Preload("ConfigTypes").
		First(&endpoint, endpointID).Error
	if err != nil {
		return nil, err
	}

	return endpoint.ConfigTypes, nil
}

func (r *Repository) ValidateUserEndpoint(userID, endpointID uint) error {
	var count int64

	err := r.db.
		Table("group_endpoints").
		Joins("JOIN users u ON u.group_id = group_endpoints.group_id").
		Where("u.id = ?", userID).
		Where("group_endpoints.endpoint_id = ?", endpointID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
