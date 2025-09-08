package main

import (
	"log"
	"github.com/gin-gonic/gin"

	"Product_Catalog_Microservice/internal/domain/service"
	"Product_Catalog_Microservice/internal/handlers"
	"Product_Catalog_Microservice/internal/repository"
	

)

// Espacio para que el compañero implemente los repositorios reales
// Deben implementar las interfaces:
//   - producto.ProductoRepositoryInterface
//   - productor.ProductorRepositoryInterface

// DummyEventPublisher es una implementación temporal de EventPublisher
type DummyEventPublisher struct{}

func (d *DummyEventPublisher) Publish(event any) error {
	// Aquí podrías loggear el evento o simplemente ignorarlo
	return nil
}


func main() {
	// Repositorios en memoria (simulación por ahora)
	productoRepo := repository.NewProductoRepository()
	productorRepo := repository.NewProductorRepository()

	// Imprimir los IDs de los productores guardados
	if all, err := productorRepo.GetAll(); err == nil {
		log.Println("Productores cargados por defecto:")
		for _, prod := range all {
			log.Printf("ID: %s, Nombre: %s\n", prod.ID, prod.Nombre.Value)
		}
	}

	// Servicio
	eventPublisher := &DummyEventPublisher{}
	catalogoService := service.NewCatalogoService(productorRepo, productoRepo, eventPublisher)

	// Handler
	productoHandler := &handlers.ProductoHandler{Catalogo: catalogoService}

	// Router con Gin
	r := gin.Default()

	// Endpoints
	r.POST("catalogo/producto", productoHandler.PublicarProducto)
	r.POST("catalogo/productos/excedente", productoHandler.MarcarProductoComoExcedente)
	r.PUT("catalogo/productos/disponibilidad", productoHandler.ActualizarDisponibilidadPorTemporada)
  	r.GET("catalogo/completo", productoHandler.GetCatalogoCompleto)
	// Iniciar servidor
	log.Println("Servidor iniciado en :8080")
	r.Run(":8080")
}