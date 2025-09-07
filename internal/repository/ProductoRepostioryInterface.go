package repository

import (
	"Product_Catalog_Microservice/internal/domain/producto"
	"fmt"
	"sync"
	"time"
)

type ProductoRepository struct {
	mu        sync.RWMutex                                            //To sync the concurrent request
	productos map[producto.ProductoID]*producto.ProductoAgroecologico //map to save the Productos Agroecologicos by ID
}

func NewProductoRepository() *ProductoRepository {
	return &ProductoRepository{
		productos: make(map[producto.ProductoID]*producto.ProductoAgroecologico),
	}
}

func (pr *ProductoRepository) Save(producto *producto.ProductoAgroecologico) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if _, exist := pr.productos[producto.ID]; exist {
		return fmt.Errorf("El producto con id %s ya existe", producto.ID)
	}

	pr.productos[producto.ID] = producto
	return nil
}

func (pr *ProductoRepository) GetByID(id producto.ProductoID) (*producto.ProductoAgroecologico, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	if prod, ok := pr.productos[id]; ok {
		response := &prod
		return *response, nil
	}

	return nil, fmt.Errorf("No se ha encontrado del producto con id %s", id)
}

func (pr *ProductoRepository) Update(producto *producto.ProductoAgroecologico) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if _, ok := pr.productos[producto.ID]; ok {
		pr.productos[producto.ID] = producto
		return nil
	}

	return fmt.Errorf("Producto con id %s no encontrado", producto.ID)
}

func (pr *ProductoRepository) GetByProductorID(productorID string) ([]*producto.ProductoAgroecologico, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	var result []*producto.ProductoAgroecologico

	for _, prod := range pr.productos {
		if prod.ProductorID == productorID {
			result = append(result, prod)
		}
	}

	return result, nil
}

func (pr *ProductoRepository) GetByCategoria(categoria producto.Categoria) ([]*producto.ProductoAgroecologico, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	var result []*producto.ProductoAgroecologico

	for _, prod := range pr.productos {
		if prod.Categoria == categoria {
			result = append(result, prod)
		}
	}

	return result, nil
}

func (pr *ProductoRepository) GetByEstado(estado producto.EstadoDisponibilidad) ([]*producto.ProductoAgroecologico, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	var result []*producto.ProductoAgroecologico

	for _, prod := range pr.productos {
		if prod.Estado == estado {
			result = append(result, prod)
		}
	}

	return result, nil
}

func (pr *ProductoRepository) GetByUbicacion(ubicacion producto.Ubicacion) ([]*producto.ProductoAgroecologico, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	var result []*producto.ProductoAgroecologico

	for _, prod := range pr.productos {
		if prod.Ubicacion == ubicacion {
			result = append(result, prod)
		}
	}

	return result, nil
}

func (pr *ProductoRepository) GetAll() ([]*producto.ProductoAgroecologico, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	result := make([]*producto.ProductoAgroecologico, 0, len(pr.productos)) // Reserv memory to no reallocate
	for _, prod := range pr.productos {
		result = append(result, prod)
	}
	return result, nil

}

func (pr *ProductoRepository) GetAvailableProducts() ([]*producto.ProductoAgroecologico, error) {
	return pr.GetByEstado(producto.EstadoDisponibilidad{Value: producto.Disponible})
}

func (pr *ProductoRepository) GetProductsInSeason(now time.Time) ([]*producto.ProductoAgroecologico, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	var result []*producto.ProductoAgroecologico

	for _, prod := range pr.productos {
		if prod.Temporada.IsInSeason(now) {
			result = append(result, prod)
		}
	}

	return result, nil
}

func (pr *ProductoRepository) UpdateEstadoDisponibilidad(id producto.ProductoID, estado producto.EstadoDisponibilidad) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if prod, ok := pr.productos[id]; ok {
		prod.Estado = estado
		return nil
	}

	return fmt.Errorf("No se encontro el producto con id %s", id)
}
