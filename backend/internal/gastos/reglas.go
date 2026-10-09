package gastos

import "errors"

var ErrCamposObligatorios = errors.New("descripcion, categoria y monto (mayor a 0) son obligatorios")

func ValidarCampos(descripcion, categoria string, monto float64) error {
	if descripcion == "" || categoria == "" || monto <= 0 {
		return ErrCamposObligatorios
	}
	return nil
}
