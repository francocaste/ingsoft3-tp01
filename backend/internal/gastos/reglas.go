package gastos

import "errors"

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
