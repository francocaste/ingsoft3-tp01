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

// Presupuesto es el tope mensual definido para una categoria.
type Presupuesto struct {
	Categoria    string  `json:"categoria"`
	MontoMensual float64 `json:"monto_mensual"`
}

// ResumenCategoria es el total gastado en una categoria, con su presupuesto si tiene.
type ResumenCategoria struct {
	Categoria          string   `json:"categoria"`
	Total              float64  `json:"total"`
	PresupuestoMensual *float64 `json:"presupuesto_mensual,omitempty"`
	PorcentajeUsado    *float64 `json:"porcentaje_usado,omitempty"`
}

var ErrNoEncontrado = errors.New("no existe ese gasto")
