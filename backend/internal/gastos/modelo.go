package gastos

import (
	"errors"
	"time"
)

// Gasto representa un gasto personal cargado por el usuario.
type Gasto struct {
	ID          int       `json:"id"`
	Descripcion string    `json:"descripcion"`
	Monto       float64   `json:"monto"`
	Categoria   string    `json:"categoria"`
	Fecha       time.Time `json:"fecha"`
	Estado      string    `json:"estado"`
}

var ErrNoEncontrado = errors.New("no existe ese gasto")
