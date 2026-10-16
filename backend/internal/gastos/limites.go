package gastos

import (
	"errors"
	"strings"
)

const (
	MontoMaximoPorGasto    = 1_000_000
	LargoMaximoDescripcion = 120
)

var (
	ErrMontoMuyAlto     = errors.New("el monto supera el maximo permitido por gasto")
	ErrDescripcionLarga = errors.New("la descripcion es demasiado larga")
)

// ValidarLimites: un solo gasto no puede superar un monto maximo y la
// descripcion no puede ser kilometrica.
func ValidarLimites(descripcion string, monto float64) error {
	if monto > MontoMaximoPorGasto {
		return ErrMontoMuyAlto
	}
	if len([]rune(descripcion)) > LargoMaximoDescripcion {
		return ErrDescripcionLarga
	}
	return nil
}

// CategoriaSugerida propone una categoria a partir de la descripcion.
func CategoriaSugerida(descripcion string) string {
	d := strings.ToLower(descripcion)
	switch {
	case strings.Contains(d, "nafta") || strings.Contains(d, "colectivo") || strings.Contains(d, "uber"):
		return "Transporte"
	case strings.Contains(d, "super") || strings.Contains(d, "almuerzo") || strings.Contains(d, "cena"):
		return "Comida"
	case strings.Contains(d, "cine") || strings.Contains(d, "salida"):
		return "Ocio"
	case strings.Contains(d, "luz") || strings.Contains(d, "gas") || strings.Contains(d, "internet"):
		return "Servicios"
	default:
		return "Otros"
	}
}
