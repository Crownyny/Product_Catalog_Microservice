package producto

import (
    "errors"
    "time"
)

type ProductoID string

// Entidad raíz del agregado ProductoAgroecologico
type ProductoAgroecologico struct {
    ID               ProductoID
    Nombre           NombreProducto
    Descripcion      DescripcionProducto
    Categoria        Categoria
    TipoProduccion   TipoProduccion
    Temporada        TemporadaLocal
    Estado           EstadoDisponibilidad
    Ubicacion        Ubicacion
    Imagen           Imagen
    ProductorID      string // referencia por identidad al productor
    publicadoEn      time.Time

	eventsPending    []interface{}
}

// Constructor del agregado
func NewProductoAgroecologico(
    id ProductoID,
    nombre NombreProducto,
    desc DescripcionProducto,
    categoria Categoria,
    tipo TipoProduccion,
    temporada TemporadaLocal,
    ubicacion Ubicacion,
    imagen Imagen,
    productorID string,
) (*ProductoAgroecologico, error) {
    if productorID == "" {
        return nil, errors.New("productorID cannot be empty")
    }

    estado := EstadoDisponibilidad{
        Value: Disponible, 
    }

    producto := &ProductoAgroecologico{
        ID:             id,
        Nombre:         nombre,
        Descripcion:    desc,
        Categoria:      categoria,
        TipoProduccion: tipo,
        Temporada:      temporada,
        Estado:         estado,
        Ubicacion:      ubicacion,
        Imagen:         imagen,
        ProductorID:    productorID,
        publicadoEn:    time.Now(),
        eventsPending:  make([]interface{}, 0),
    }
    
    // Generar evento de producto publicado
    producto.addEvent(ProductoPublicado{
        ProductoID: id,
        At:         time.Now(),
    })
    
    return producto, nil
}

func (p *ProductoAgroecologico) MarcarComoExcedente(now time.Time) error {
    if p.Temporada.IsInSeason(now) {
        return errors.New("no se puede marcar como 'Excedente' dentro de la temporada")
    }
    p.Estado = EstadoDisponibilidad{Value: Excedente}
    
    // Generar evento
    p.addEvent(ProductoMarcadoComoExcedente{
        ProductoID: p.ID,
        At:         now,
    })
    
    return nil
}

func (p *ProductoAgroecologico) Agotar() error {
    if p.Estado.Value != Disponible {
        return errors.New("solo un producto 'Disponible' puede marcarse como 'Agotado'")
    }
    p.Estado = EstadoDisponibilidad{Value: Agotado}
    
    // Generar evento
    p.addEvent(ProductoAgotado{
        ProductoID: p.ID,
        At:         time.Now(),
    })
    
    return nil
}

// Recalcula el estado de disponibilidad en base a la temporada actual
func (p *ProductoAgroecologico) RecalcularDisponibilidad(now time.Time) {
    if p.Temporada.IsInSeason(now) {
        p.Estado = EstadoDisponibilidad{Value: Disponible}
    } else if p.Estado.Value != Excedente { 
        p.Estado = EstadoDisponibilidad{Value: Agotado}
    }
}

func (p *ProductoAgroecologico) ActualizarInformacion(nombre NombreProducto, desc DescripcionProducto, imagen Imagen) error {
    // Validar que el producto no esté en un estado que impida actualizaciones
    if p.Estado.Value == Agotado {
        return errors.New("no se puede actualizar información de un producto agotado")
    }
    
    p.Nombre = nombre
    p.Descripcion = desc
    p.Imagen = imagen
    return nil
}

// Métodos para manejar eventos
func (p *ProductoAgroecologico) addEvent(event interface{}) {
    p.eventsPending = append(p.eventsPending, event)
}

func (p *ProductoAgroecologico) GetPendingEvents() []interface{} {
    return p.eventsPending
}

func (p *ProductoAgroecologico) ClearEvents() {
    p.eventsPending = make([]interface{}, 0)
}
