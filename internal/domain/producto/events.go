package producto

import "time"

type ProductoPublicado struct {
    ProductoID ProductoID
    At         time.Time
}

type ProductoMarcadoComoExcedente struct {
    ProductoID ProductoID
    At         time.Time
}

type ProductoAgotado struct {
    ProductoID ProductoID
    At         time.Time
}
