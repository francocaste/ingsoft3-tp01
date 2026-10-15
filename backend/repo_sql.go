package main

import (
	"database/sql"
	"time"

	"gastos-backend/internal/gastos"
)

// repoSQL implementa gastos.Repositorio sobre Postgres. Es la unica parte
// que habla SQL para crear gastos y cambiar su estado; se prueba con
// tests de integracion (con base real), no con los tests unitarios.
type repoSQL struct {
	db *sql.DB
}

func (r *repoSQL) PresupuestoDe(categoria string) (float64, bool, error) {
	return presupuestoDeCategoria(categoria)
}

func (r *repoSQL) TotalDelMes(categoria string, fecha time.Time) (float64, error) {
	var total sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT SUM(monto) FROM gastos
		WHERE categoria = $1
		  AND date_trunc('month', fecha) = date_trunc('month', $2::date)
	`, categoria, fecha).Scan(&total)
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return total.Float64, nil
}

func (r *repoSQL) Insertar(g *gastos.Gasto) error {
	return r.db.QueryRow(
		`INSERT INTO gastos (descripcion, monto, categoria, fecha, estado) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		g.Descripcion, g.Monto, g.Categoria, g.Fecha, g.Estado,
	).Scan(&g.ID)
}

func (r *repoSQL) EstadoDe(id int) (string, error) {
	var estado string
	err := r.db.QueryRow(`SELECT estado FROM gastos WHERE id = $1`, id).Scan(&estado)
	if err == sql.ErrNoRows {
		return "", gastos.ErrNoEncontrado
	}
	return estado, err
}

func (r *repoSQL) ActualizarEstado(id int, estado string) error {
	_, err := r.db.Exec(`UPDATE gastos SET estado = $1 WHERE id = $2`, estado, id)
	return err
}
