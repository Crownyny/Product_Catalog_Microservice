package productor

type ProductorRepositoryInterface interface {
    Save(productor *Productor) error
    GetByID(id ProductorID) (*Productor, error)
    Delete(id ProductorID) error // Establece al productor como inactivo

    GetByUbicacion(ubicacion Ubicacion) ([]*Productor, error)
    GetByEstadoVerificacion(estado EstadoVerificacion) ([]*Productor, error)
    GetByReputacionMinima(minReputacion Reputacion) ([]*Productor, error)
    GetVerificados() ([]*Productor, error)
    GetPendientesVerificacion() ([]*Productor, error)
    GetAll() ([]*Productor, error)
    
    UpdateReputacion(id ProductorID, nuevaReputacion Reputacion) error
    UpdateEstadoVerificacion(id ProductorID, nuevoEstado EstadoVerificacion) error
}