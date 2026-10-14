package gastos

import (
	"errors"
	"testing"
	"time"
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

func TestValidarFecha(t *testing.T) {
	// Arrange: "ahora" fijo, así el test da lo mismo hoy que en un año
	ahora := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

	casos := []struct {
		nombre    string
		fecha     time.Time
		quiereErr error
	}{
		{"hoy mismo", ahora, nil},
		{"hoy a ultima hora (misma fecha)", time.Date(2026, 10, 9, 23, 59, 0, 0, time.UTC), nil},
		{"ayer", ahora.AddDate(0, 0, -1), nil},
		{"manana", ahora.AddDate(0, 0, 1), ErrFechaFutura},
		{"el anio que viene", ahora.AddDate(1, 0, 0), ErrFechaFutura},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Act
			err := ValidarFecha(c.fecha, ahora)

			// Assert
			if !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
		})
	}
}

func TestValidarPresupuesto(t *testing.T) {
	casos := []struct {
		nombre      string
		totalDelMes float64
		monto       float64
		presupuesto float64
		tiene       bool
		quiereErr   error
	}{
		{"sin presupuesto definido, cualquier monto", 0, 999999, 0, false, nil},
		{"dentro del presupuesto", 1000, 500, 2000, true, nil},
		{"justo en el tope (se permite)", 1500, 500, 2000, true, nil},
		{"se pasa por un centavo", 1500, 500.01, 2000, true, ErrPresupuestoExcedido},
		{"ya estaba al tope y suma algo", 2000, 1, 2000, true, ErrPresupuestoExcedido},
		{"un solo gasto mas grande que el tope", 0, 5000, 2000, true, ErrPresupuestoExcedido},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Act
			err := ValidarPresupuesto(c.totalDelMes, c.monto, c.presupuesto, c.tiene)

			// Assert
			if !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
		})
	}
}
func TestValidarTransicion(t *testing.T) {
	casos := []struct {
		actual, nuevo string
		quiereErr     error
	}{
		{EstadoPendiente, EstadoPagado, nil},
		{EstadoPendiente, EstadoPendiente, nil},
		{EstadoPagado, EstadoPagado, nil},
		{EstadoPagado, EstadoPendiente, ErrTransicionInvalida},
		{EstadoPendiente, "cancelado", ErrEstadoInvalido},
		{EstadoPagado, "", ErrEstadoInvalido},
	}

	for _, c := range casos {
		t.Run(c.actual+"->"+c.nuevo, func(t *testing.T) {
			// Act
			err := ValidarTransicion(c.actual, c.nuevo)

			// Assert
			if !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
		})
	}
}

func TestPorcentajeUsado(t *testing.T) {
	casos := []struct {
		nombre             string
		total, presupuesto float64
		quiere             float64
	}{
		{"mitad", 500, 1000, 50},
		{"completo", 1000, 1000, 100},
		{"excedido", 1500, 1000, 150},
		{"nada gastado", 0, 1000, 0},
		{"presupuesto cero no divide por cero", 100, 0, 0},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := PorcentajeUsado(c.total, c.presupuesto); got != c.quiere {
				t.Fatalf("quise %v y obtuve %v", c.quiere, got)
			}
		})
	}
}
