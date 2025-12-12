package repositories

import (
	"monitoring-energy-service/internal/domain/entities"
	"monitoring-energy-service/internal/domain/ports/output"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExampleRepository struct {
	db *gorm.DB
}

var _ output.ExampleRepositoryInterface = &ExampleRepository{}

func NewExampleRepository(db *gorm.DB) *ExampleRepository {
	return &ExampleRepository{db: db}
}

func (r *ExampleRepository) FindByID(id uuid.UUID) (*entities.ExampleEntity, error) {
	var entity entities.ExampleEntity
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ExampleRepository) FindAll() ([]*entities.ExampleEntity, error) {
	var entities []*entities.ExampleEntity
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *ExampleRepository) Create(entity *entities.ExampleEntity) (*entities.ExampleEntity, error) {
	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ExampleRepository) Update(entity *entities.ExampleEntity) (*entities.ExampleEntity, error) {
	if err := r.db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *ExampleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.ExampleEntity{}, "id = ?", id).Error
}
