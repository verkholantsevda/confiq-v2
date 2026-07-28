package configtypes

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) List() ([]ConfigType, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id uint) (*ConfigType, error) {
	configType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return configType, nil
}

func (s *Service) Create(req CreateConfigTypeRequest) (*ConfigType, error) {

	existing, err := s.repo.GetByName(req.Name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrConfigTypeAlreadyExists
	}

	configType := &ConfigType{
		Name:              req.Name,
		Description:       req.Description,
		ConfigTemplate:    req.ConfigTemplate,
		UsageInstructions: req.UsageInstructions,
		ClientLinks:       req.ClientLinks,
		IsActive:          req.IsActive,
	}

	if err := s.repo.Create(configType); err != nil {
		return nil, err
	}

	return configType, nil
}

func (s *Service) Update(id uint, req UpdateConfigTypeRequest) (*ConfigType, error) {

	configType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrConfigTypeNotFound
	}

	if configType.Name != req.Name {
		existing, err := s.repo.GetByName(req.Name)
		if err != nil {
			return nil, err
		}

		if existing != nil {
			return nil, ErrConfigTypeAlreadyExists
		}
	}

	configType.Name = req.Name
	configType.Description = req.Description
	configType.ConfigTemplate = req.ConfigTemplate
	configType.UsageInstructions = req.UsageInstructions
	configType.ClientLinks = req.ClientLinks
	configType.IsActive = req.IsActive

	if err := s.repo.Update(configType); err != nil {
		return nil, err
	}

	return configType, nil
}

func (s *Service) Delete(id uint) error {

	configType, err := s.repo.GetByID(id)
	if err != nil {
		return ErrConfigTypeNotFound
	}

	if err := s.repo.Delete(configType.ID); err != nil {
		return err
	}

	return nil
}
