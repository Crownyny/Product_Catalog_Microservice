package producto

import "time"

type ProductoRepository interface {
    Save(producto *ProductoAgroecologico) error
    GetByID(id ProductoID) (*ProductoAgroecologico, error)
    Update(nombre NombreProducto, desc DescripcionProducto, imagen Imagen) error
    GetByProductorID(productorID string) ([]*ProductoAgroecologico, error)
    GetByCategoria(categoria Categoria) ([]*ProductoAgroecologico, error)
    GetByEstado(estado EstadoDisponibilidad) ([]*ProductoAgroecologico, error)
    GetByUbicacion(ubicacion Ubicacion) ([]*ProductoAgroecologico, error)
    GetAll() ([]*ProductoAgroecologico, error)
    GetAvailableProducts() ([]*ProductoAgroecologico, error)
    GetProductsInSeason(now time.Time) ([]*ProductoAgroecologico, error)
    UpdateEstadoDisponibilidad(id ProductoID, estado EstadoDisponibilidad) error
}