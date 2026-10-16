package gastos

import (
	"errors"
	"strings"
	"testing"
)

func TestValidarLimites(t *testing.T) {
	casos := []struct {
		nombre      string
		descripcion string
		monto       float64
		quiereErr   error
	}{
		{"gasto normal", "Almuerzo", 1500, nil},
		{"justo en el maximo de monto", "Auto", MontoMaximoPorGasto, nil},
		{"monto sobre el maximo", "Auto", MontoMaximoPorGasto + 1, ErrMontoMuyAlto},
		{"descripcion justo en el largo maximo", strings.Repeat("a", LargoMaximoDescripcion), 10, nil},
		{"descripcion demasiado larga", strings.Repeat("a", LargoMaximoDescripcion+1), 10, ErrDescripcionLarga},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if err := ValidarLimites(c.descripcion, c.monto); !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
		})
	}
}

func TestCategoriaSugerida(t *testing.T) {
	casos := []struct {
		descripcion string
		quiere      string
	}{
		{"Nafta YPF", "Transporte"},
		{"Colectivo", "Transporte"},
		{"Uber al aeropuerto", "Transporte"},
		{"Super Coto", "Comida"},
		{"Almuerzo con amigos", "Comida"},
		{"Cena", "Comida"},
		{"Cine", "Ocio"},
		{"Salida del sabado", "Ocio"},
		{"Factura de luz", "Servicios"},
		{"Internet", "Servicios"},
		{"Regalo", "Otros"},
	}
	for _, c := range casos {
		t.Run(c.descripcion, func(t *testing.T) {
			if got := CategoriaSugerida(c.descripcion); got != c.quiere {
				t.Fatalf("quise %q y obtuve %q", c.quiere, got)
			}
		})
	}
}
