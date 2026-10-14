package gastos

import (
	"errors"
	"testing"
	"time"
)

// repoFalso es un doble de prueba del Repositorio. Hace dos cosas:
//   - STUB: devuelve datos que nosotros fijamos (presupuesto, total, estado).
//   - SPY: cuenta cuantas veces lo llamaron, para verificar COMPORTAMIENTO.
type repoFalso struct {
	presupuesto float64
	tienePresup bool
	totalDelMes float64
	estado      string

	// errores que le podemos pedir que devuelva, para simular una base caída
	errAlConsultar  error
	errAlTotal      error
	errAlInsertar   error
	errAlActualizar error

	// espía: cuántas veces lo llamaron y con qué
	llamadasPresupuesto  int
	llamadasInsertar     int
	llamadasActualizar   int
	ultimoEstadoGuardado string
}

func (r *repoFalso) PresupuestoDe(string) (float64, bool, error) {
	r.llamadasPresupuesto++
	return r.presupuesto, r.tienePresup, r.errAlConsultar
}

func (r *repoFalso) TotalDelMes(string, time.Time) (float64, error) {
	return r.totalDelMes, r.errAlTotal
}

func (r *repoFalso) Insertar(g *Gasto) error {
	r.llamadasInsertar++
	g.ID = 42
	return r.errAlInsertar
}

func (r *repoFalso) EstadoDe(int) (string, error) {
	if r.estado == "" {
		return "", ErrNoEncontrado
	}
	return r.estado, nil
}

func (r *repoFalso) ActualizarEstado(_ int, estado string) error {
	r.llamadasActualizar++
	r.ultimoEstadoGuardado = estado
	return r.errAlActualizar
}

var hoyFijo = time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

func servicioConRepo(r *repoFalso) *Servicio {
	return NuevoServicio(r, func() time.Time { return hoyFijo })
}

func gastoValido() Gasto {
	return Gasto{Descripcion: "Almuerzo", Monto: 1500, Categoria: "Comida", Fecha: hoyFijo}
}

func TestCrearGasto_Valido_InsertaUnaVezYQuedaPendiente(t *testing.T) {
	// Arrange
	repo := &repoFalso{}
	svc := servicioConRepo(repo)
	// Act
	creado, err := svc.CrearGasto(gastoValido())
	// Assert
	if err != nil {
		t.Fatalf("no esperaba error y obtuve %v", err)
	}
	if repo.llamadasInsertar != 1 {
		t.Fatalf("debia insertar 1 vez y inserto %d", repo.llamadasInsertar)
	}
	if creado.Estado != EstadoPendiente || creado.ID != 42 {
		t.Fatalf("estado/ID inesperados: %+v", creado)
	}
}

func TestCrearGasto_ReglaIncumplida_NoTocaLaBase(t *testing.T) {
	casos := []struct {
		nombre    string
		modificar func(*Gasto)
		quiereErr error
	}{
		{"monto cero", func(g *Gasto) { g.Monto = 0 }, ErrCamposObligatorios},
		{"categoria invalida", func(g *Gasto) { g.Categoria = "Inventada" }, ErrCategoriaInvalida},
		{"fecha futura", func(g *Gasto) { g.Fecha = hoyFijo.AddDate(0, 0, 3) }, ErrFechaFutura},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			repo := &repoFalso{}
			g := gastoValido()
			c.modificar(&g)
			// Act
			_, err := servicioConRepo(repo).CrearGasto(g)
			// Assert: el error correcto Y la base ni se enteró
			if !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
			if repo.llamadasPresupuesto != 0 || repo.llamadasInsertar != 0 {
				t.Fatalf("no debia tocar la base (presupuesto=%d, insertar=%d)",
					repo.llamadasPresupuesto, repo.llamadasInsertar)
			}
		})
	}
}

func TestCrearGasto_ExcedePresupuesto_NoInserta(t *testing.T) {
	repo := &repoFalso{presupuesto: 2000, tienePresup: true, totalDelMes: 1800}

	_, err := servicioConRepo(repo).CrearGasto(gastoValido()) // 1800 + 1500 > 2000

	if !errors.Is(err, ErrPresupuestoExcedido) {
		t.Fatalf("quise ErrPresupuestoExcedido y obtuve %v", err)
	}
	if repo.llamadasInsertar != 0 {
		t.Fatal("un gasto que excede el presupuesto no debe insertarse")
	}
}

func TestCrearGasto_ErroresDeLaBase_SePropagan(t *testing.T) {
	errBase := errors.New("la base se cayo")
	casos := []struct {
		nombre string
		repo   *repoFalso
	}{
		{"falla al consultar el presupuesto", &repoFalso{errAlConsultar: errBase}},
		{"falla al consultar el total del mes", &repoFalso{tienePresup: true, presupuesto: 9999, errAlTotal: errBase}},
		{"falla al insertar", &repoFalso{errAlInsertar: errBase}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := servicioConRepo(c.repo).CrearGasto(gastoValido())
			if !errors.Is(err, errBase) {
				t.Fatalf("quise propagar %v y obtuve %v", errBase, err)
			}
		})
	}
}

func TestCambiarEstado_PendienteAPagado_Actualiza(t *testing.T) {
	repo := &repoFalso{estado: EstadoPendiente}

	err := servicioConRepo(repo).CambiarEstado(1, EstadoPagado)

	if err != nil {
		t.Fatalf("no esperaba error y obtuve %v", err)
	}
	if repo.llamadasActualizar != 1 || repo.ultimoEstadoGuardado != EstadoPagado {
		t.Fatalf("debia guardar 'pagado' una vez (llamadas=%d, estado=%q)",
			repo.llamadasActualizar, repo.ultimoEstadoGuardado)
	}
}

func TestCambiarEstado_PagadoAPendiente_SeRechazaYNoActualiza(t *testing.T) {
	repo := &repoFalso{estado: EstadoPagado}

	err := servicioConRepo(repo).CambiarEstado(1, EstadoPendiente)

	if !errors.Is(err, ErrTransicionInvalida) {
		t.Fatalf("quise ErrTransicionInvalida y obtuve %v", err)
	}
	if repo.llamadasActualizar != 0 {
		t.Fatal("una transicion invalida no debe escribir en la base")
	}
}

func TestCambiarEstado_CasosBorde(t *testing.T) {
	errBase := errors.New("la base se cayo")
	casos := []struct {
		nombre    string
		repo      *repoFalso
		nuevo     string
		quiereErr error
	}{
		{"estado inexistente", &repoFalso{estado: EstadoPendiente}, "cancelado", ErrEstadoInvalido},
		{"gasto que no existe", &repoFalso{estado: ""}, EstadoPagado, ErrNoEncontrado},
		{"falla al guardar", &repoFalso{estado: EstadoPendiente, errAlActualizar: errBase}, EstadoPagado, errBase},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := servicioConRepo(c.repo).CambiarEstado(1, c.nuevo)
			if !errors.Is(err, c.quiereErr) {
				t.Fatalf("quise %v y obtuve %v", c.quiereErr, err)
			}
		})
	}
}

func TestNuevoServicio_SinReloj_UsaTimeNow(t *testing.T) {
	svc := NuevoServicio(&repoFalso{}, nil)

	if svc.ahora == nil {
		t.Fatal("debia usar time.Now por defecto")
	}
}
