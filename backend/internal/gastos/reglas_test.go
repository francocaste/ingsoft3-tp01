package gastos

import (
	"errors"
	"testing"
)

func TestValidarCampos_TodoCompleto_NoDevuelveError(t *testing.T) {
	// Arrange: preparo los datos de un gasto bien cargado
	descripcion := "Almuerzo"
	categoria := "Comida"
	monto := 1500.0

	// Act: ejecuto la regla
	err := ValidarCampos(descripcion, categoria, monto)

	// Assert: espero que NO haya error
	if err != nil {
		t.Fatalf("no esperaba error y obtuve: %v", err)
	}
}

func TestValidarCampos_MontoCero_DevuelveError(t *testing.T) {
	// Arrange
	descripcion := "Almuerzo"
	categoria := "Comida"
	monto := 0.0

	// Act
	err := ValidarCampos(descripcion, categoria, monto)

	// Assert: espero justamente ErrCamposObligatorios
	if !errors.Is(err, ErrCamposObligatorios) {
		t.Fatalf("esperaba ErrCamposObligatorios y obtuve: %v", err)
	}
}
