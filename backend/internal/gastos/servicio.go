package gastos

import "time"

// Repositorio es lo que el servicio necesita de la base de datos. Es una
// interfaz a proposito: en produccion la implementa Postgres y en los tests
// la reemplazamos por un doble, asi probamos la logica sin base de datos.
type Repositorio interface {
	PresupuestoDe(categoria string) (monto float64, tiene bool, err error)
	TotalDelMes(categoria string, fecha time.Time) (float64, error)
	Insertar(g *Gasto) error
	EstadoDe(id int) (string, error) // devuelve ErrNoEncontrado si no existe
	ActualizarEstado(id int, estado string) error
}

// Servicio orquesta las reglas y el acceso a datos.
type Servicio struct {
	repo  Repositorio
	ahora func() time.Time // inyectable para controlar "hoy" en los tests
}

func NuevoServicio(repo Repositorio, ahora func() time.Time) *Servicio {
	if ahora == nil {
		ahora = time.Now
	}
	return &Servicio{repo: repo, ahora: ahora}
}

// CrearGasto aplica las reglas 1 a 4 en orden y recien despues inserta.
// Si alguna regla falla, NO se toca la base.
func (s *Servicio) CrearGasto(g Gasto) (Gasto, error) {
	if err := ValidarCampos(g.Descripcion, g.Categoria, g.Monto); err != nil {
		return Gasto{}, err
	}
	if err := ValidarCategoria(g.Categoria); err != nil {
		return Gasto{}, err
	}

	ahora := s.ahora()
	if g.Fecha.IsZero() {
		g.Fecha = ahora
	}
	if err := ValidarFecha(g.Fecha, ahora); err != nil {
		return Gasto{}, err
	}

	g.Estado = EstadoPendiente

	presupuesto, tiene, err := s.repo.PresupuestoDe(g.Categoria)
	if err != nil {
		return Gasto{}, err
	}
	if tiene {
		total, err := s.repo.TotalDelMes(g.Categoria, g.Fecha)
		if err != nil {
			return Gasto{}, err
		}
		if err := ValidarPresupuesto(total, g.Monto, presupuesto, tiene); err != nil {
			return Gasto{}, err
		}
	}

	if err := s.repo.Insertar(&g); err != nil {
		return Gasto{}, err
	}
	return g, nil
}

// CambiarEstado aplica la regla 5 (transicion de estado).
func (s *Servicio) CambiarEstado(id int, nuevo string) error {
	if nuevo != EstadoPendiente && nuevo != EstadoPagado {
		return ErrEstadoInvalido
	}
	actual, err := s.repo.EstadoDe(id)
	if err != nil {
		return err
	}
	if err := ValidarTransicion(actual, nuevo); err != nil {
		return err
	}
	return s.repo.ActualizarEstado(id, nuevo)
}
