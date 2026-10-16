package gastos

import (
	"errors"
	"time"
)

var ErrCamposObligatorios = errors.New("descripcion, categoria y monto (mayor a 0) son obligatorios")

func ValidarCampos(descripcion, categoria string, monto float64) error {
	if descripcion == "" || categoria == "" || monto <= 0 {
		return ErrCamposObligatorios
	}
	return nil
}

var CategoriasValidas = map[string]bool{
	"Comida":     true,
	"Transporte": true,
	"Ocio":       true,
	"Servicios":  true,
	"Salud":      true,
	"Otros":      true,
}

var ErrCategoriaInvalida = errors.New("categoria invalida")

func ValidarCategoria(categoria string) error {
	if !CategoriasValidas[categoria] {
		return ErrCategoriaInvalida
	}
	return nil
}

var ErrFechaFutura = errors.New("la fecha del gasto no puede ser futura")

// TruncarADia deja sólo la fecha (hora 00:00:00).
func TruncarADia(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func ValidarFecha(fecha, ahora time.Time) error {
	if TruncarADia(fecha).After(TruncarADia(ahora)) {
		return ErrFechaFutura
	}
	return nil
}

var ErrPresupuestoExcedido = errors.New("el gasto excede el presupuesto mensual de la categoria")

func ValidarPresupuesto(totalDelMes, monto, presupuesto float64, tienePresupuesto bool) error {
	if !tienePresupuesto {
		return nil
	}
	if totalDelMes+monto >= presupuesto {
		return ErrPresupuestoExcedido
	}
	return nil
}

const (
	EstadoPendiente = "pendiente"
	EstadoPagado    = "pagado"
)

var (
	ErrEstadoInvalido     = errors.New("estado invalido: debe ser 'pendiente' o 'pagado'")
	ErrTransicionInvalida = errors.New("un gasto pagado no puede volver a pendiente")
)

func ValidarTransicion(actual, nuevo string) error {
	if nuevo != EstadoPendiente && nuevo != EstadoPagado {
		return ErrEstadoInvalido
	}
	if actual == EstadoPagado && nuevo == EstadoPendiente {
		return ErrTransicionInvalida
	}
	return nil
}

// PorcentajeUsado devuelve cuanto del presupuesto se uso (0 a 100+).
// Con presupuesto 0 devuelve 0 para no dividir por cero.
func PorcentajeUsado(total, presupuesto float64) float64 {
	if presupuesto <= 0 {
		return 0
	}
	return (total / presupuesto) * 100
}
