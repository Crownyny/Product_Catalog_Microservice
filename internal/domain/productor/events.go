package productor

import "time"

type ProductorEnVerificacion struct {
    ProductorID ProductorID
    At         time.Time
}

type ProductorVerificado struct{
	ProductorID ProductorID
    At         time.Time
}

type ReputacionActualizada struct {
    ProductorID    ProductorID
    NuevaReputacion Reputacion
    At             time.Time
}

