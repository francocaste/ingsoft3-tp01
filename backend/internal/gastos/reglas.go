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
