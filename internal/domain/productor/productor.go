package productor

import (
	"errors"
	"time"
)

type ProductorID string

type Productor struct {
	ID               ProductorID
	Nombre           NombreProductor
	Ubicacion        Ubicacion
	EstadoVerificacion EstadoVerificacion
	EstadoActividad  EstadoActividad
	Reputacion       Reputacion
	PracticasCultivo PracticasDeCultivo
	    // Agregar eventos pendientes
    eventsPending      []interface{}
}

// NewProductor crea un nuevo Productor con validaciones para mantener invariantes
func NewProductor(
	id ProductorID,
	nombre NombreProductor,
	ubicacion Ubicacion,
	estadoVerificacion EstadoVerificacion,
	estadoActividad EstadoActividad,
	reputacion Reputacion,
	practicasCultivo PracticasDeCultivo,
) (*Productor, error) {

	if id == "" {
		return nil, errors.New("el ID del productor no puede estar vacío")
	}

	return &Productor{
		ID:                id,
		Nombre:            nombre,
		Ubicacion:         ubicacion,
		EstadoVerificacion: estadoVerificacion,
		EstadoActividad:   estadoActividad,
		Reputacion:        reputacion,
		PracticasCultivo:  practicasCultivo,
	}, nil
}

// PuedePublicar determina si el productor puede publicar productos
func (p *Productor) PuedePublicar(minReputacion Reputacion) bool {
	return p.EstadoVerificacion.IsVerificado() && p.Reputacion >= minReputacion && p.EstadoActividad.IsActivo()
}

// ActualizarReputacion permite actualizar la reputacion del productor basándose en cálculos derivados de historial
func (p *Productor) ActualizarReputacion(nuevaReputacion Reputacion) error {
	if (nuevaReputacion < 0 || nuevaReputacion > 5)  && p.EstadoActividad.IsActivo() {
		return errors.New("reputacion fuera de rango permitido")
	}

    reputacionAnterior := p.Reputacion
    p.Reputacion = nuevaReputacion
    
    // Generar evento solo si cambió
    if reputacionAnterior != nuevaReputacion {
        p.addEvent(ReputacionActualizada{
            ProductorID:     p.ID,
            NuevaReputacion: nuevaReputacion,
            At:              time.Now(),
        })
    }
    
    return nil
}


func (p *Productor) IniciarProcesosVerificacion() error {
    if !p.EstadoActividad.IsActivo() {
        return errors.New("el productor no está activo")
    }

    if p.EstadoVerificacion.IsVerificado() {
        return errors.New("el productor ya está verificado")
    }
    if p.EstadoVerificacion.Value == "En Proceso" {
        return errors.New("ya hay un proceso de verificación en curso")
    }
    
    p.EstadoVerificacion = EstadoVerificacion{Value: "En Proceso"}
    
    // Generar evento
    p.addEvent(ProductorEnVerificacion{
        ProductorID: p.ID,
        At:          time.Now(),
    })
    
    return nil
}

func (p *Productor) VerificarProductor() error {
	if !p.EstadoVerificacion.IsEnProceso() {
		return errors.New("el productor no está en proceso de verificación")
	}

	p.EstadoVerificacion = EstadoVerificacion{Value: "Verificado"}

	// Generar evento
	p.addEvent(ProductorVerificado{
		ProductorID: p.ID,
		At:         time.Now(),
	})

	return nil
}


// Métodos para manejar eventos
func (p *Productor) addEvent(event interface{}) {
    p.eventsPending = append(p.eventsPending, event)
}

func (p *Productor) GetPendingEvents() []interface{} {
    return p.eventsPending
}

func (p *Productor) ClearEvents() {
    p.eventsPending = make([]interface{}, 0)
}