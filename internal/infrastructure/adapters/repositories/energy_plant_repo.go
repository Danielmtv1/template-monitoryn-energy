package repositories

import (
	"errors"

	"monitoring-energy-service/internal/domain/entities"
	domainerrors "monitoring-energy-service/internal/domain/errors"
	"monitoring-energy-service/internal/domain/ports/output"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnergyPlantRepository implementa la capa de persistencia para plantas de energía
//
// PROPÓSITO:
// Proporciona acceso a la base de datos PostgreSQL para validar que las plantas existen.
// Se usa principalmente para validar plant_source_id antes de guardar eventos.
//
// CAMBIO REALIZADO: Archivo creado desde cero
// RAZÓN: Necesitábamos validar que los eventos solo se guarden si la planta existe
type EnergyPlantRepository struct {
	db *gorm.DB
}

var _ output.EnergyPlantRepositoryInterface = &EnergyPlantRepository{}

// NewEnergyPlantRepository crea una nueva instancia del repositorio de plantas
// PARÁMETROS: db - Conexión GORM a PostgreSQL
func NewEnergyPlantRepository(db *gorm.DB) *EnergyPlantRepository {
	return &EnergyPlantRepository{db: db}
}

// FindByID busca una planta por su UUID
// CAMBIO: Método nuevo
// RAZÓN: Permite validar que una planta existe antes de guardar eventos
// CAMBIO: Ahora traduce gorm.ErrRecordNotFound a domainerrors.ErrNotFound
// RAZÓN: Permite a los handlers distinguir entre 404 (not found) y 500 (internal error)
func (r *EnergyPlantRepository) FindByID(id uuid.UUID) (*entities.EnergyPlants, error) {
	var plant entities.EnergyPlants
	if err := r.db.First(&plant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}
	return &plant, nil
}

// Exists verifica rápidamente si una planta existe por UUID
// CAMBIO: Método nuevo
// RAZÓN: Método optimizado para solo verificar existencia sin cargar toda la entidad
func (r *EnergyPlantRepository) Exists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entities.EnergyPlants{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
