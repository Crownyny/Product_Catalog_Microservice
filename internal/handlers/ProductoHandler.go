// ...existing code...
package handlers

import (
	"Product_Catalog_Microservice/internal/domain/producto"
	"Product_Catalog_Microservice/internal/domain/productor"
	"Product_Catalog_Microservice/internal/domain/service"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ProductoHandler struct {
	Catalogo *service.CatalogoService
}

// ...existing code...

func (h *ProductoHandler) PublicarProducto(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		ProductorID     string  `json:"productor_id"`
		ProductoID      string  `json:"producto_id"`
		Nombre          string  `json:"nombre"`
		Descripcion     string  `json:"descripcion"`
		Categoria       string  `json:"categoria"`
		TipoProduccion  string  `json:"tipo_produccion"`
		TemporadaInicio string  `json:"temporada_inicio"` // formato: "2006-01-02"
		TemporadaFin    string  `json:"temporada_fin"`    // formato: "2006-01-02"
		ZonaVeredal     string  `json:"zona_veredal"`
		Finca           string  `json:"finca"`
		ImagenURL       string  `json:"imagen_url"`
		ImagenDesc      string  `json:"imagen_desc"`
		MinReputacion   float32 `json:"min_reputacion"`
	}

	var req requestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Value objects y validaciones
	productorID := req.ProductorID
	//productoID := req.ProductoID //<= el id deberia generarse en el front
	//dado que no tenemos se genera en el controlador

	productoID := producto.ProductoID(uuid.New().String())

	nombre, err := producto.NewNombreProducto(req.Nombre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	desc, err := producto.NewDescripcionProducto(req.Descripcion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	categoria, err := producto.NewCategoria(req.Categoria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tipo := producto.TipoProduccion(req.TipoProduccion)
	temporadaInicio, err := time.Parse("2006-01-02", req.TemporadaInicio)
	if err != nil {
		http.Error(w, "Formato de fecha de inicio inválido", http.StatusBadRequest)
		return
	}
	temporadaFin, err := time.Parse("2006-01-02", req.TemporadaFin)
	if err != nil {
		http.Error(w, "Formato de fecha de fin inválido", http.StatusBadRequest)
		return
	}
	temporada, err := producto.NewTemporadaLocal(temporadaInicio, temporadaFin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ubicacion, err := producto.NewUbicacion(req.ZonaVeredal, req.Finca)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	imagen, err := producto.NewImagen(req.ImagenURL, req.ImagenDesc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	minReputacion, err := productor.NuevaReputacion(req.MinReputacion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prod, err := h.Catalogo.PublicarProducto(
		productor.ProductorID(productorID),
		producto.ProductoID(productoID),
		nombre,
		desc,
		categoria,
		tipo,
		temporada,
		ubicacion,
		imagen,
		minReputacion,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prod)
}

func (h *ProductoHandler) MarcarProductoComoExcedente(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		ProductoID string `json:"producto_id"`
		Fecha      string `json:"fecha"` // formato: "2006-01-02"
	}

	var req requestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	productoID := producto.ProductoID(req.ProductoID)
	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		http.Error(w, "Formato de fecha inválido", http.StatusBadRequest)
		return
	}

	err = h.Catalogo.MarcarProductoComoExcedente(productoID, fecha)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductoHandler) ActualizarDisponibilidadPorTemporada(w http.ResponseWriter, r *http.Request) {
	// Oaquí se usa la fecha actual del servidor
	now := time.Now()

	err := h.Catalogo.ActualizarDisponibilidadPorTemporada(now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ...existing code...
