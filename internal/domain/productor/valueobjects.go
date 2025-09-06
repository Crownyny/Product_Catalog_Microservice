package productor

import (
	"errors"
	"regexp"
	"strings"
)

// NombreProducto representa el nombre de un producto como value object.
// Garantiza que el nombre sea válido y cumpla con las reglas de negocio.
type NombreProductor struct {
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
func NewNombreProducto(value string) (NombreProductor, error) {
	if value == "" {
		return NombreProductor{}, errors.New("el nombre del productor no puede estar vacío")
	}
	if len(value) > 80 {
		return NombreProductor{}, errors.New("el nombre del productor no puede superar 80 caracteres")
	}
	return NombreProductor{Value: value}, nil
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

// EstadoVerificacion representa si el productor esta verificado por la plataforma.
// Puede ser "Verificado" o "No Verificado".
type EstadoVerificacion struct {
	Value string
}

// Constantes que definen los estados de verificación válidos
const (
	Verificado     string = "Verificado"     // Productor verificado
	NoVerificado   string = "No Verificado"   // Productor no verificado
	EnProceso	  string = "En Proceso"      // Productor en proceso de verificación
)

// NewEstadoVerificacion crea una nueva instancia de EstadoVerificacion.
// Valida que el estado sea uno de los valores permitidos.
//
// Parámetros:
//   - value: el estado de verificación del productor
//
// Retorna:
//   - EstadoVerificacion: instancia válida del value object
//   - error: error de validación si el estado es inválido
func NewEstadoVerificacion(value string) (EstadoVerificacion, error) {
	switch value {
	case Verificado, NoVerificado, EnProceso:
		return EstadoVerificacion{Value: value}, nil
	default:
		return EstadoVerificacion{}, errors.New("estado de verificación inválido")
	}	
}

func (e EstadoVerificacion) IsVerificado() bool {
	return e.Value == Verificado
}

func (e EstadoVerificacion) IsEnProceso() bool {
	return e.Value == EnProceso
}

// Reputacion representa la reputacion promedio del productor, valor entre 0 y 5 inclusive
type Reputacion float32

// NuevaReputacion crea una nueva instancia de Reputacion.
// Valida que el valor esté entre 0 y 5 inclusive.
//
// Parámetros:
//   - valor: la reputación del productor
//
// Retorna:
//   - Reputacion: instancia válida del value object
//   - error: error de validación si el valor es inválido
func NuevaReputacion(valor float32) (Reputacion, error) {
	if valor < 0 || valor > 5 {
		return 0, errors.New("reputacion debe estar entre 0 y 5")
	}
	return Reputacion(valor), nil
}

// PracticasDeCultivo representa las prácticas utilizadas por el productor en sus cultivos.
// Debe ser un texto validado, acotado y coherente con el lenguaje ubicuo local.
type PracticasDeCultivo struct {
	Descripcion string
}

// NuevaPracticasDeCultivo crea una nueva instancia de PracticasDeCultivo.
// Valida que la descripción no esté vacía, no sea demasiado larga y tenga sentido.
//
// Parámetros:
//   - descripcion: descripción de las prácticas de cultivo del productor
//
// Retorna:
//   - PracticasDeCultivo: instancia válida del value object
//   - error: error de validación si la descripción es inválida
func NuevaPracticasDeCultivo(descripcion string) (PracticasDeCultivo, error) {
	descripcion = strings.TrimSpace(descripcion)
	if descripcion == "" {
		return PracticasDeCultivo{}, errors.New("descripcion de prácticas no puede estar vacía")
	}
	if len(descripcion) > 500 {
		return PracticasDeCultivo{}, errors.New("descripcion de prácticas demasiado larga")
	}

	return PracticasDeCultivo{Descripcion: descripcion}, nil
}

// EstadoActividad representa si el productor está activo en la plataforma.
// Un productor puede estar activo, inactivo o suspendido.
type EstadoActividad struct {
    Value string
}

// Constantes que definen los estados de actividad válidos
const (
    Activo     string = "Activo"     // Productor activo y operativo
    Inactivo   string = "Inactivo"   // Productor temporalmente inactivo
    Suspendido string = "Suspendido" // Productor suspendido por la plataforma
)

// NewEstadoActividad crea una nueva instancia de EstadoActividad.
// Valida que el estado sea uno de los valores permitidos.
//
// Parámetros:
//   - value: el estado de actividad del productor
//
// Retorna:
//   - EstadoActividad: instancia válida del value object
//   - error: error de validación si el estado es inválido
func NewEstadoActividad(value string) (EstadoActividad, error) {
    switch value {
    case Activo, Inactivo, Suspendido:
        return EstadoActividad{Value: value}, nil
    default:
        return EstadoActividad{}, errors.New("estado de actividad inválido")
    }
}

// IsActivo verifica si el productor está activo
func (e EstadoActividad) IsActivo() bool {
    return e.Value == Activo
}