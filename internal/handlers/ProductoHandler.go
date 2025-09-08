// ...existing code...
package handlers


import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"Product_Catalog_Microservice/internal/domain/producto"
	"Product_Catalog_Microservice/internal/domain/productor"
	"Product_Catalog_Microservice/internal/domain/service"
)

type ProductoHandler struct {
    Catalogo *service.CatalogoService
}

// POST /productos/publicar
func (h *ProductoHandler) PublicarProducto(c *gin.Context) {
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
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
        return
    }

    // Generación de IDs y value objects
    productorID := req.ProductorID
    productoID := producto.ProductoID(uuid.New().String()) // forzado en backend

    nombre, err := producto.NewNombreProducto(req.Nombre)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    desc, err := producto.NewDescripcionProducto(req.Descripcion)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    categoria, err := producto.NewCategoria(req.Categoria)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    tipo := producto.TipoProduccion(req.TipoProduccion)

    temporadaInicio, err := time.Parse("2006-01-02", req.TemporadaInicio)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha de inicio inválido"})
        return
    }
    temporadaFin, err := time.Parse("2006-01-02", req.TemporadaFin)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha de fin inválido"})
        return
    }
    temporada, err := producto.NewTemporadaLocal(temporadaInicio, temporadaFin)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ubicacion, err := producto.NewUbicacion(req.ZonaVeredal, req.Finca)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    imagen, err := producto.NewImagen(req.ImagenURL, req.ImagenDesc)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    minReputacion, err := productor.NuevaReputacion(req.MinReputacion)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, prod)
}

// POST /productos/excedente
func (h *ProductoHandler) MarcarProductoComoExcedente(c *gin.Context) {
    type requestBody struct {
        ProductoID string `json:"producto_id"`
        Fecha      string `json:"fecha"` // formato: "2006-01-02"
    }

    var req requestBody
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
        return
    }

    productoID := producto.ProductoID(req.ProductoID)
    fecha, err := time.Parse("2006-01-02", req.Fecha)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
        return
    }

    if err := h.Catalogo.MarcarProductoComoExcedente(productoID, fecha); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

// PUT /productos/disponibilidad
func (h *ProductoHandler) ActualizarDisponibilidadPorTemporada(c *gin.Context) {
    now := time.Now()

    if err := h.Catalogo.ActualizarDisponibilidadPorTemporada(now); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
// ...existing code...

func (h *ProductoHandler) GetCatalogoCompleto(c *gin.Context) {
    catalogo, err := h.Catalogo.GetCatalogoCompleto()
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, catalogo)
}