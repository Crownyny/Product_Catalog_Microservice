// Package producto contiene los value objects del dominio de productos
// para el microservicio de catálogo de productos agrícolas.
package producto

import (
	"errors"
	"regexp"
	"time"
)

// NombreProducto representa el nombre de un producto como value object.
// Garantiza que el nombre sea válido y cumpla con las reglas de negocio.
type NombreProducto struct {
	Value string
}

// NewNombreProducto crea una nueva instancia de NombreProducto.
// Valida que el nombre no esté vacío y no supere los 100 caracteres.
//
// Parámetros:
//   - value: el nombre del producto
//
// Retorna:
//   - NombreProducto: instancia válida del value object
//   - error: error de validación si el nombre es inválido
func NewNombreProducto(value string) (NombreProducto, error) {
	if value == "" {
		return NombreProducto{}, errors.New("el nombre del producto no puede estar vacío")
	}
	if len(value) > 100 {
		return NombreProducto{}, errors.New("el nombre del producto no puede superar 100 caracteres")
	}
	return NombreProducto{Value: value}, nil
}

// DescripcionProducto representa la descripción de un producto como value object.
// Garantiza que la descripción tenga una longitud adecuada para ser informativa.
type DescripcionProducto struct {
	Value string
}

// NewDescripcionProducto crea una nueva instancia de DescripcionProducto.
// Valida que la descripción tenga entre 10 y 500 caracteres.
//
// Parámetros:
//   - value: la descripción del producto
//
// Retorna:
//   - DescripcionProducto: instancia válida del value object
//   - error: error de validación si la descripción es inválida
func NewDescripcionProducto(value string) (DescripcionProducto, error) {
	if len(value) < 10 {
		return DescripcionProducto{}, errors.New("la descripción debe tener al menos 10 caracteres")
	}
	if len(value) > 500 {
		return DescripcionProducto{}, errors.New("la descripción no puede superar 500 caracteres")
	}
	return DescripcionProducto{Value: value}, nil
}

// Categoria representa las categorías válidas de productos agrícolas.
// Define un conjunto limitado de categorías para clasificar los productos.
type Categoria string

// Constantes que definen las categorías válidas de productos
const (
	CategoriaFruta     Categoria = "Fruta"           // Productos frutales
	CategoriaHortaliza Categoria = "Hortaliza"       // Vegetales y hortalizas
	CategoriaTuberculo Categoria = "Tubérculo"       // Tubérculos como papa, yuca
	CategoriaMedicinal Categoria = "PlantaMedicinal" // Plantas con propiedades medicinales
	CategoriaLacteo    Categoria = "Lácteo"          // Productos lácteos
)

// NewCategoria crea una nueva instancia de Categoria.
// Valida que la categoría sea una de las categorías predefinidas válidas.
//
// Parámetros:
//   - value: el valor de la categoría como string
//
// Retorna:
//   - Categoria: instancia válida del value object
//   - error: error de validación si la categoría no es válida
func NewCategoria(value string) (Categoria, error) {
	switch Categoria(value) {
	case CategoriaFruta, CategoriaHortaliza, CategoriaTuberculo, CategoriaMedicinal, CategoriaLacteo:
		return Categoria(value), nil
	default:
		return "", errors.New("categoría inválida")
	}
}

// TipoProduccion representa los diferentes métodos de producción agrícola.
// Define los tipos de producción según las prácticas utilizadas.
type TipoProduccion string

// Constantes que definen los tipos de producción válidos
const (
	ProduccionAgroecologica TipoProduccion = "Agroecologico" // Producción agroecológica
	ProduccionOrganica      TipoProduccion = "Organico"      // Producción orgánica 
	ProduccionTradicional   TipoProduccion = "Tradicional"   // Producción tradicional
)

// TemporadaLocal representa el período de temporada local de un producto.
// Define cuándo está disponible naturalmente en la región.
type TemporadaLocal struct {
	Inicio time.Time // Fecha de inicio de la temporada
	Fin    time.Time // Fecha de fin de la temporada
}

// NewTemporadaLocal crea una nueva instancia de TemporadaLocal.
// Valida que la fecha de fin no sea anterior a la fecha de inicio.
//
// Parámetros:
//   - inicio: fecha de inicio de la temporada
//   - fin: fecha de fin de la temporada
//
// Retorna:
//   - TemporadaLocal: instancia válida del value object
//   - error: error de validación si las fechas son inválidas
func NewTemporadaLocal(inicio, fin time.Time) (TemporadaLocal, error) {
	if fin.Before(inicio) {
		return TemporadaLocal{}, errors.New("la fecha de fin no puede ser antes del inicio")
	}

	if fin.Before(time.Now()) {
		return TemporadaLocal{}, errors.New("la fecha de fin no puede estar en el pasado")
	}

	if fin.Sub(inicio).Hours() > 24*365 {
		return TemporadaLocal{}, errors.New("la temporada no puede durar más de un año")
	}

	return TemporadaLocal{Inicio: inicio, Fin: fin}, nil
}

// Funcion auxiliar para saber si actualmente está en temporada
func (t TemporadaLocal) IsInSeason(now time.Time) bool {
    return (now.Equal(t.Inicio) || now.After(t.Inicio)) &&
           (now.Equal(t.Fin) || now.Before(t.Fin))
}


// EstadoDisponibilidad representa el estado actual de disponibilidad de un producto.
// Indica si el producto está disponible, agotado o en excedente.
type EstadoDisponibilidad struct {
	Value string
}

// Constantes que definen los estados de disponibilidad válidos
const (
	Disponible string = "Disponible" // Producto disponible para venta
	Agotado    string = "Agotado"    // Producto temporalmente agotado
	Excedente  string = "Excedente"  // Producto en excedente/abundancia
)

// NewEstadoDisponibilidad crea una nueva instancia de EstadoDisponibilidad.
// Valida que el estado de disponibilidad sea uno de los estados predefinidos válidos.
//
// Parámetros:
//   - value: el valor del estado de disponibilidad como string
//
// Retorna:
//   - EstadoDisponibilidad: instancia válida del value object
//   - error: error de validación si el estado no es válido
// Retorna:
//   - EstadoDisponibilidad: instancia válida del value object
//   - error: error de validación si el estado no es válido
func NewEstadoDisponibilidad(value string) (EstadoDisponibilidad, error) {
    switch value {
    case Disponible, Agotado, Excedente:
        return EstadoDisponibilidad{Value: value}, nil
    default:
        return EstadoDisponibilidad{}, errors.New("estado de disponibilidad inválido")
    }
}

// Ubicacion representa la ubicación geográfica donde se produce el producto.
// Incluye información sobre la zona veredal y la finca específica.
type Ubicacion struct {
	ZonaVeredal string // Zona veredal donde se encuentra la finca
	Finca       string // Nombre de la finca productora
}

// NewUbicacion crea una nueva instancia de Ubicacion.
// Valida que tanto la zona veredal como la finca estén especificadas,
// que no excedan la longitud máxima y que no contengan caracteres prohibidos.
//
// Parámetros:
//   - zona: nombre de la zona veredal (máximo 50 caracteres)
//   - finca: nombre de la finca (máximo 80 caracteres)
//
// Retorna:
//   - Ubicacion: instancia válida del value object
//   - error: error de validación si algún campo es inválido
func NewUbicacion(zona, finca string) (Ubicacion, error) {
    // Validar campos vacíos
    if zona == "" || finca == "" {
        return Ubicacion{}, errors.New("zona veredal y finca no pueden estar vacíos")
    }

    // Validar longitud máxima
    if len(zona) > 40 {
        return Ubicacion{}, errors.New("la zona veredal no puede superar 40 caracteres")
    }
    if len(finca) > 50 {
        return Ubicacion{}, errors.New("el nombre de la finca no puede superar 50 caracteres")
    }

    // Validar caracteres prohibidos
    if err := validarCaracteresProhibidos(zona, "zona veredal"); err != nil {
        return Ubicacion{}, err
    }
    if err := validarCaracteresProhibidos(finca, "finca"); err != nil {
        return Ubicacion{}, err
    }

    return Ubicacion{ZonaVeredal: zona, Finca: finca}, nil
}

// validarCaracteresProhibidos valida que el texto solo contenga caracteres permitidos
// para nombres de ubicaciones (letras, números, espacios, guiones, apostrofes, puntos).
func validarCaracteresProhibidos(texto, campo string) error {
    // Permite letras (incluye acentos), números, espacios, guiones, apostrofes y puntos
    patron := regexp.MustCompile(`^[a-zA-ZáéíóúñüÁÉÍÓÚÑÜ0-9\s\-'\.]+$`)
    if !patron.MatchString(texto) {
        return errors.New("el campo " + campo + " contiene caracteres no permitidos")
    }
    return nil
}

// Imagen representa una imagen asociada a un producto.
// Contiene la URL de la imagen y una descripción corta para accesibilidad.
type Imagen struct {
	URL              string // URL de la imagen del producto
	DescripcionCorta string // Descripción corta de la imagen para accesibilidad
}

// NewImagen crea una nueva instancia de Imagen.
// Valida que la URL tenga un formato válido (HTTP o HTTPS).
//
// Parámetros:
//   - url: URL de la imagen (debe comenzar con http:// o https://)
//   - desc: descripción corta de la imagen
//
// Retorna:
//   - Imagen: instancia válida del value object
//   - error: error de validación si la URL no es válida
func NewImagen(url, desc string) (Imagen, error) {
	regex := regexp.MustCompile(`^https?://`)
	if !regex.MatchString(url) {
		return Imagen{}, errors.New("la URL de la imagen no es válida")
	}
	return Imagen{URL: url, DescripcionCorta: desc}, nil
}
