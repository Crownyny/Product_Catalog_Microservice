package repository

import (
	"Product_Catalog_Microservice/internal/domain/productor"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type ProductorRepository struct {
	mu          sync.RWMutex // To sync the concurrent request
	productores map[productor.ProductorID]*productor.Productor
}

func NewProductorRepository() *ProductorRepository {
	return &ProductorRepository{
		productores: make(map[productor.ProductorID]*productor.Productor),
	}
}

func (pr *ProductorRepository) Save(pro *productor.Productor) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pro.ID = productor.ProductorID(uuid.New().String())

	if _, exist := pr.productores[pro.ID]; exist {
		return fmt.Errorf("El producotr con id %s ya existe", pro.ID)
	}

	pr.productores[pro.ID] = pro
	return nil
}

func (pr *ProductorRepository) GetByID(id productor.ProductorID) (*productor.Productor, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	if prod, ok := pr.productores[id]; ok {
		response := *prod
		return &response, nil
	}
	return nil, fmt.Errorf("No se ha encontrado el productor con id %s", id)
}

func (pr *ProductorRepository) Delete(id productor.ProductorID) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if productorFound, ok := pr.productores[id]; ok {
		productorFound.EstadoActividad = productor.EstadoActividad{
			Value: productor.Inactivo,
		}
		return nil
	}

	return fmt.Errorf("No se ha encontrado el productor con id %s", id)
}
func (pr *ProductorRepository) GetByUbicacion(ubicacion productor.Ubicacion) ([]*productor.Productor, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var result []*productor.Productor
	for _, prod := range pr.productores {
		if prod.Ubicacion == ubicacion {
			result = append(result, prod)
		}
	}
	return result, nil
}

func (pr *ProductorRepository) GetByEstadoVerificacion(estado productor.EstadoVerificacion) ([]*productor.Productor, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var result []*productor.Productor
	for _, prod := range pr.productores {
		if prod.EstadoVerificacion == estado {
			result = append(result, prod)
		}
	}
	return result, nil
}

func (pr *ProductorRepository) GetByReputacionMinima(minReputacion productor.Reputacion) ([]*productor.Productor, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var result []*productor.Productor
	for _, prod := range pr.productores {
		if prod.Reputacion >= minReputacion {
			result = append(result, prod)
		}
	}
	return result, nil
}

func (pr *ProductorRepository) GetVerificados() ([]*productor.Productor, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var result []*productor.Productor
	for _, prod := range pr.productores {
		if prod.EstadoVerificacion.IsVerificado() {
			result = append(result, prod)
		}
	}
	return result, nil
}

func (pr *ProductorRepository) GetPendientesVerificacion() ([]*productor.Productor, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var result []*productor.Productor
	for _, prod := range pr.productores {
		if prod.EstadoVerificacion.IsEnProceso() {
			result = append(result, prod)
		}
	}
	return result, nil
}

func (pr *ProductorRepository) GetAll() ([]*productor.Productor, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	var result []*productor.Productor
	for _, prod := range pr.productores {
		result = append(result, prod)
	}
	return result, nil
}

func (pr *ProductorRepository) UpdateReputacion(id productor.ProductorID, nuevaReputacion productor.Reputacion) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	if prod, ok := pr.productores[id]; ok {
		prod.Reputacion = nuevaReputacion
		return nil
	}
	return fmt.Errorf("No se encontró el productor con id %s", id)
}

func (pr *ProductorRepository) UpdateEstadoVerificacion(id productor.ProductorID, nuevoEstado productor.EstadoVerificacion) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	if prod, ok := pr.productores[id]; ok {
		prod.EstadoVerificacion = nuevoEstado
		return nil
	}
	return fmt.Errorf("No se encontró el productor con id %s", id)
}
