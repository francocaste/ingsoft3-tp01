package gastos

import (
	"errors"
	"testing"
)

func TestValidarCampos(t *testing.T) {
	// Arrange: la tabla. Cada fila es un caso: datos de entrada + resultado esperado.
	casos := []struct {
		nombre      string
		descripcion string
		categoria   string
		monto       float64
		quiereErr   error
	}{
		{"todo completo", "Almuerzo", "Comida", 1500, nil},
		{"monto minimo valido", "Chicle", "Otros", 0.01, nil},
		{"sin descripcion", "", "Comida", 1500, ErrCamposObligatorios},
		{"sin categoria", "Almuerzo", "", 1500, ErrCamposObligatorios},
		{"monto cero", "Almuerzo", "Comida", 0, ErrCamposObligatorios},
		{"monto negativo", "Almuerzo", "Comida", -50, ErrCamposObligatorios},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Act
			err := ValidarCampos(c.descripcion, c.categoria, c.monto)

			// Assert
			if !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
		})
	}
}
func TestValidarCategoria(t *testing.T) {
	// Arrange: la tabla de casos
	casos := []struct {
		categoria string
		quiereErr error
	}{
		{"Comida", nil},
		{"Transporte", nil},
		{"Ocio", nil},
		{"Servicios", nil},
		{"Salud", nil},
		{"Otros", nil},
		{"Inventada", ErrCategoriaInvalida},
		{"", ErrCategoriaInvalida},
		{"comida", ErrCategoriaInvalida}, // distingue mayusculas
	}

	for _, c := range casos {
		t.Run("categoria="+c.categoria, func(t *testing.T) {
			// Act
			err := ValidarCategoria(c.categoria)

			// Assert
			if !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
		})
	}
}
