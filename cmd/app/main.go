package main

import (
	"log"
	"net/http"

	"Product_Catalog_Microservice/internal/domain/service"
	"Product_Catalog_Microservice/internal/handlers"
	"Product_Catalog_Microservice/internal/repository"

	// Importa los paquetes de dominio para los tipos
	"Product_Catalog_Microservice/internal/domain/producto"
	"Product_Catalog_Microservice/internal/domain/productor"
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
	// TODO: Instanciar los repositorios reales cuando estén implementados
	var productoRepo producto.ProductoRepositoryInterface
	var productorRepo productor.ProductorRepositoryInterface

	productoRepo = repository.NewProductoRepository()
	productorRepo = repository.NewProductorRepository()

	eventPublisher := &DummyEventPublisher{}

	catalogoService := service.NewCatalogoService(productorRepo, productoRepo, eventPublisher)
	productoHandler := &handlers.ProductoHandler{Catalogo: catalogoService}

	http.HandleFunc("/productos/publicar", productoHandler.PublicarProducto)
	http.HandleFunc("/productos/excedente", productoHandler.MarcarProductoComoExcedente)
	http.HandleFunc("/productos/actualizar-disponibilidad", productoHandler.ActualizarDisponibilidadPorTemporada)

	log.Println("Servidor iniciado en :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
