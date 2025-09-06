package service

import (
    "errors"
    "time"

    "Product_Catalog_Microservice/internal/domain/producto"
    "Product_Catalog_Microservice/internal/domain/productor"
)

// EventPublisher define la interfaz para publicar eventos de dominio
type EventPublisher interface {
    Publish(event any) error
}

type CatalogoService struct {
    productorRepo  productor.ProductorRepositoryInterface
    productoRepo   producto.ProductoRepositoryInterface
    eventPublisher EventPublisher
}

func NewCatalogoService(
    productorRepo productor.ProductorRepositoryInterface,
    productoRepo producto.ProductoRepositoryInterface,
    eventPublisher EventPublisher,
) *CatalogoService {
    return &CatalogoService{
        productorRepo:  productorRepo,
        productoRepo:   productoRepo,
        eventPublisher: eventPublisher,
    }
}

// PublicarProducto valida que el productor pueda publicar y crea el producto
func (s *CatalogoService) PublicarProducto(
    productorID productor.ProductorID,
    productoID producto.ProductoID,
    nombre producto.NombreProducto,
    desc producto.DescripcionProducto,
    categoria producto.Categoria,
    tipo producto.TipoProduccion,
    temporada producto.TemporadaLocal,
    ubicacion producto.Ubicacion,
    imagen producto.Imagen,
    minReputacion productor.Reputacion,
) (*producto.ProductoAgroecologico, error) {
    
    // Verificar que el productor existe y puede publicar
    prod, err := s.productorRepo.GetByID(productorID)
    if err != nil {
        return nil, errors.New("productor no encontrado")
    }
    
    if !prod.PuedePublicar(minReputacion) {
        return nil, errors.New("el productor no está autorizado para publicar productos")
    }
    
    // Crear el producto (esto genera el evento ProductoPublicado)
    nuevoProducto, err := producto.NewProductoAgroecologico(
        productoID,
        nombre,
        desc,
        categoria,
        tipo,
        temporada,
        ubicacion,
        imagen,
        string(productorID),
    )
    if err != nil {
        return nil, err
    }
    
    // Guardar el producto
    if err := s.productoRepo.Save(nuevoProducto); err != nil {
        return nil, err
    }
    
    // Publicar eventos generados por el agregado
    s.publishPendingEvents(nuevoProducto)
    
    return nuevoProducto, nil
}

// IniciarVerificacionProductor inicia el proceso de verificación de un productor
func (s *CatalogoService) IniciarVerificacionProductor(productorID productor.ProductorID) error {
    prod, err := s.productorRepo.GetByID(productorID)
    if err != nil {
        return errors.New("productor no encontrado")
    }
    
    // Esto genera el evento ProductorEnVerificacion
    if err := prod.IniciarProcesosVerificacion(); err != nil {
        return err
    }
    
    // Actualizar el estado en el repositorio
    if err := s.productorRepo.UpdateEstadoVerificacion(productorID, prod.EstadoVerificacion); err != nil {
        return err
    }
    
    // Publicar eventos generados por el agregado
    s.publishPendingEvents(prod)
    
    return nil
}

// CompletarVerificacionProductor completa la verificación de un productor
func (s *CatalogoService) CompletarVerificacionProductor(productorID productor.ProductorID) error {
    prod, err := s.productorRepo.GetByID(productorID)
    if err != nil {
        return errors.New("productor no encontrado")
    }
    
    // Esto genera el evento ProductorVerificado
    if err := prod.VerificarProductor(); err != nil {
        return err
    }
    
    // Actualizar el estado en el repositorio
    if err := s.productorRepo.UpdateEstadoVerificacion(productorID, prod.EstadoVerificacion); err != nil {
        return err
    }
    
    // Publicar eventos generados por el agregado
    s.publishPendingEvents(prod)
    
    return nil
}

// ActualizarReputacionProductor actualiza la reputación de un productor
func (s *CatalogoService) ActualizarReputacionProductor(
    productorID productor.ProductorID, 
    nuevaReputacion productor.Reputacion,
) error {
    prod, err := s.productorRepo.GetByID(productorID)
    if err != nil {
        return errors.New("productor no encontrado")
    }
    
    // Esto genera el evento ReputacionActualizada si la reputación cambia
    if err := prod.ActualizarReputacion(nuevaReputacion); err != nil {
        return err
    }
    
    // Actualizar la reputación en el repositorio
    if err := s.productorRepo.UpdateReputacion(productorID, nuevaReputacion); err != nil {
        return err
    }
    
    // Publicar eventos generados por el agregado
    s.publishPendingEvents(prod)
    
    return nil
}

// MarcarProductoComoExcedente marca un producto como excedente
func (s *CatalogoService) MarcarProductoComoExcedente(
    productoID producto.ProductoID, 
    now time.Time,
) error {
    prod, err := s.productoRepo.GetByID(productoID)
    if err != nil {
        return errors.New("producto no encontrado")
    }
    
    // Esto genera el evento ProductoMarcadoComoExcedente
    if err := prod.MarcarComoExcedente(now); err != nil {
        return err
    }
    
    // Actualizar el estado en el repositorio
    if err := s.productoRepo.UpdateEstadoDisponibilidad(productoID, prod.Estado); err != nil {
        return err
    }
    
    // Publicar eventos generados por el agregado
    s.publishPendingEvents(prod)
    
    return nil
}

// AgotarProducto marca un producto como agotado
func (s *CatalogoService) AgotarProducto(productoID producto.ProductoID) error {
    prod, err := s.productoRepo.GetByID(productoID)
    if err != nil {
        return errors.New("producto no encontrado")
    }
    
    // Esto genera el evento ProductoAgotado
    if err := prod.Agotar(); err != nil {
        return err
    }
    
    // Actualizar el estado en el repositorio
    if err := s.productoRepo.UpdateEstadoDisponibilidad(productoID, prod.Estado); err != nil {
        return err
    }
    
    // Publicar eventos generados por el agregado
    s.publishPendingEvents(prod)
    
    return nil
}

// ActualizarInformacionProducto actualiza la información básica de un producto
func (s *CatalogoService) ActualizarInformacionProducto(
    productoID producto.ProductoID,
    nombre producto.NombreProducto,
    desc producto.DescripcionProducto,
    imagen producto.Imagen,
) error {
    prod, err := s.productoRepo.GetByID(productoID)
    if err != nil {
        return errors.New("producto no encontrado")
    }
    
    if err := prod.ActualizarInformacion(nombre, desc, imagen); err != nil {
        return err
    }
    
     if err := s.productoRepo.Update(prod); err != nil {
        return err
     }
    

    return nil
}

// GetProductosByProductor obtiene todos los productos de un productor
func (s *CatalogoService) GetProductosByProductor(productorID productor.ProductorID) ([]*producto.ProductoAgroecologico, error) {
    // Verificar que el productor existe
    _, err := s.productorRepo.GetByID(productorID)
    if err != nil {
        return nil, errors.New("productor no encontrado")
    }
    
    return s.productoRepo.GetByProductorID(string(productorID))
}

// GetProductosDisponiblesEnZona obtiene productos disponibles de productores verificados en una zona
func (s *CatalogoService) GetProductosDisponiblesEnZona(ubicacion productor.Ubicacion) ([]*producto.ProductoAgroecologico, error) {
    // Obtener productores verificados en la zona
    productoresZona, err := s.productorRepo.GetByUbicacion(ubicacion)
    if err != nil {
        return nil, err
    }
    
    var todosProductos []*producto.ProductoAgroecologico
    
    for _, prod := range productoresZona {
        if prod.EstadoVerificacion.IsVerificado() && prod.EstadoActividad.IsActivo() {
            productos, err := s.productoRepo.GetByProductorID(string(prod.ID))
            if err != nil {
                continue // Continúar con el siguiente productor
            }
            
            // Filtrar solo productos disponibles
            for _, producto := range productos {
                if producto.Estado.Value == "Disponible" {
                    todosProductos = append(todosProductos, producto)
                }
            }
        }
    }
    
    return todosProductos, nil
}

// ActualizarDisponibilidadPorTemporada actualiza la disponibilidad de productos según la temporada
func (s *CatalogoService) ActualizarDisponibilidadPorTemporada(now time.Time) error {
    productos, err := s.productoRepo.GetAll()
    if err != nil {
        return err
    }
    
    for _, prod := range productos {
        estadoAnterior := prod.Estado.Value
        prod.RecalcularDisponibilidad(now)
        
        // Solo actualizar si el estado cambió
        if prod.Estado.Value != estadoAnterior {
            if err := s.productoRepo.UpdateEstadoDisponibilidad(prod.ID, prod.Estado); err != nil {
                // Log el error pero continúa con los demás productos
                continue
            }
            
            // Publicar eventos si los hay (RecalcularDisponibilidad podría generar eventos)
            s.publishPendingEvents(prod)
        }
    }
    
    return nil
}

// GetCatalogoCompleto obtiene el catálogo completo con información de productores
func (s *CatalogoService) GetCatalogoCompleto() (*CatalogoCompleto, error) {
    productos, err := s.productoRepo.GetAvailableProducts()
    if err != nil {
        return nil, err
    }
    
    productores, err := s.productorRepo.GetVerificados()
    if err != nil {
        return nil, err
    }
    
    return &CatalogoCompleto{
        Productos:   productos,
        Productores: productores,
        GeneradoEn:  time.Now(),
    }, nil
}

// GetProductoresAptosParaPublicar obtiene productores que pueden publicar productos
func (s *CatalogoService) GetProductoresAptosParaPublicar(minReputacion productor.Reputacion) ([]*productor.Productor, error) {
    productores, err := s.productorRepo.GetByReputacionMinima(minReputacion)
    if err != nil {
        return nil, err
    }
    
    var productoresAptos []*productor.Productor
    for _, prod := range productores {
        if prod.PuedePublicar(minReputacion) {
            productoresAptos = append(productoresAptos, prod)
        }
    }
    
    return productoresAptos, nil
}

// Método auxiliar para publicar eventos pendientes de cualquier agregado
func (s *CatalogoService) publishPendingEvents(aggregate any) {
    var events []interface{}
    
    // Type assertion para obtener eventos según el tipo de agregado
    switch agg := aggregate.(type) {
    case *producto.ProductoAgroecologico:
        events = agg.GetPendingEvents()
        agg.ClearEvents()
    case *productor.Productor:
        events = agg.GetPendingEvents()
        agg.ClearEvents()
    }
    
    // Publicar cada evento
    for _, event := range events {
        if err := s.eventPublisher.Publish(event); err != nil {
			//TODO: IDK what the hell put here, but is a recommended validation
        }
    }
}

// CatalogoCompleto representa una vista completa del catálogo
type CatalogoCompleto struct {
    Productos   []*producto.ProductoAgroecologico
    Productores []*productor.Productor
    GeneradoEn  time.Time
}