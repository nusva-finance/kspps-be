package services

import (
	"nusvakspps/repositories"
)

type SavingTypesService struct {
	repo *repositories.SavingTypesRepository
}

func NewSavingTypesService() *SavingTypesService {
	return &SavingTypesService{
		repo: repositories.NewSavingTypesRepository(),
	}
}

// GetAllTypes retrieves all saving types
func (s *SavingTypesService) GetAllTypes(includeInactive bool) (interface{}, error) {
	var types interface{}
	var err error

	if includeInactive {
		types, err = s.repo.ListAll()
	} else {
		types, err = s.repo.List()
	}

	if err != nil {
		return nil, err
	}

	return types, nil
}

// GetTypeByID retrieves a saving type by ID
func (s *SavingTypesService) GetTypeByID(id uint) (interface{}, error) {
	return s.repo.FindByID(id)
}

// CreateType creates a new saving type
func (s *SavingTypesService) CreateType(req interface{}) (interface{}, error) {
	// Implementation would parse request and call repo.Create
	// This is a placeholder for service logic
	return nil, nil
}

// UpdateType updates an existing saving type
func (s *SavingTypesService) UpdateType(id uint, req interface{}) (interface{}, error) {
	// Implementation would parse request and call repo.Update
	// This is a placeholder for service logic
	return nil, nil
}

// DeleteType deletes a saving type
func (s *SavingTypesService) DeleteType(id uint) error {
	return s.repo.Delete(id)
}

// InitializeDefaultTypes initializes default saving types
func (s *SavingTypesService) InitializeDefaultTypes() error {
	return s.repo.InitializeDefaultTypes()
}
