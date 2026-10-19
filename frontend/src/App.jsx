import { useEffect, useState } from 'react'
import {
  CATEGORIAS,
  calcularTotal,
  claseDePorcentaje,
  enviarGasto,
  formularioValido,
  hoyISO,
} from './lib/gastos'
// Nota: agregar a App.css la clase .excedido { color: #f87171; font-weight: bold; }

export default function App() {
  const [gastos, setGastos] = useState([])
  const [resumen, setResumen] = useState([])
  const [presupuestos, setPresupuestos] = useState([])

  const [descripcion, setDescripcion] = useState('')
  const [monto, setMonto] = useState('')
  const [categoria, setCategoria] = useState(CATEGORIAS[0])
  const [fecha, setFecha] = useState(hoyISO())

  const [presupuestoCategoria, setPresupuestoCategoria] = useState(CATEGORIAS[0])
  const [presupuestoMonto, setPresupuestoMonto] = useState('')

  const [error, setError] = useState('')
  const [cargando, setCargando] = useState(true)

  async function cargar() {
    try {
      const [gastosRes, resumenRes, presupuestosRes] = await Promise.all([
        fetch('/api/gastos'),
        fetch('/api/resumen'),
        fetch('/api/presupuestos'),
      ])
      setGastos(await gastosRes.json())
      setResumen(await resumenRes.json())
      setPresupuestos(await presupuestosRes.json())
    } finally {
      setCargando(false)
    }
  }

  useEffect(() => {
    cargar()
  }, [])

  // Las reglas del formulario viven en lib/gastos.js (y estan testeadas ahi).
  const puedeEnviar = formularioValido({ descripcion, monto, fecha })

  // Aviso visible mientras el formulario es invalido: el boton deshabilitado solo
  // no le dice al usuario QUE esta mal.
  const avisoDeValidacion =
    descripcion !== '' && descripcion.trim() === ''
      ? 'La descripción no puede estar vacía'
      : monto !== '' && !(Number(monto) > 0)
        ? 'El monto tiene que ser mayor a 0'
        : ''

  async function agregarGasto(e) {
    e.preventDefault()
    setError('')
    const resultado = await enviarGasto(fetch, { descripcion, monto, categoria, fecha })
    if (!resultado.ok) {
      setError(resultado.error)
      return
    }
    setDescripcion('')
    setMonto('')
    cargar()
  }

  async function eliminarGasto(id) {
    await fetch(`/api/gastos/${id}`, { method: 'DELETE' })
    cargar()
  }

  async function marcarComoPagado(id) {
    await fetch(`/api/gastos/${id}/estado`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ estado: 'pagado' }),
    })
    cargar()
  }

  async function guardarPresupuesto(e) {
    e.preventDefault()
    setError('')
    const res = await fetch('/api/presupuestos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        categoria: presupuestoCategoria,
        monto_mensual: parseFloat(presupuestoMonto),
      }),
    })
    if (!res.ok) {
      setError(await res.text())
      return
    }
    setPresupuestoMonto('')
    cargar()
  }

  const total = calcularTotal(gastos)

  return (
    <div className="container">
      <h1>Gestión de Gastos Personales · TP6</h1>

      <form onSubmit={agregarGasto} className="form">
        <input
          aria-label="Descripción"
          placeholder="Descripción"
          value={descripcion}
          onChange={(e) => setDescripcion(e.target.value)}
          required
        />
        <input
          aria-label="Monto"
          type="number"
          step="0.01"
          min="0.01"
          placeholder="Monto"
          value={monto}
          onChange={(e) => setMonto(e.target.value)}
          required
        />
        <select
          aria-label="Categoría"
          value={categoria}
          onChange={(e) => setCategoria(e.target.value)}
        >
          {CATEGORIAS.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
        <input
          aria-label="Fecha"
          type="date"
          value={fecha}
          max={hoyISO()}
          onChange={(e) => setFecha(e.target.value)}
        />
        <button type="submit" disabled={!puedeEnviar}>
          Agregar
        </button>
      </form>

      {avisoDeValidacion && (
        <p role="alert" className="error">
          {avisoDeValidacion}
        </p>
      )}
      {error && (
        <p role="alert" className="error">
          {error}
        </p>
      )}

      <details className="presupuestos">
        <summary>Presupuestos mensuales por categoría</summary>
        <form onSubmit={guardarPresupuesto} className="form">
          <select
            aria-label="Categoría del presupuesto"
            value={presupuestoCategoria}
            onChange={(e) => setPresupuestoCategoria(e.target.value)}
          >
            {CATEGORIAS.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </select>
          <input
            aria-label="Monto mensual del presupuesto"
            type="number"
            step="0.01"
            min="0.01"
            placeholder="Monto mensual"
            value={presupuestoMonto}
            onChange={(e) => setPresupuestoMonto(e.target.value)}
            required
          />
          <button type="submit" disabled={!(Number(presupuestoMonto) > 0)}>
            Guardar presupuesto
          </button>
        </form>
        <ul>
          {presupuestos.map((p) => (
            <li key={p.categoria}>
              {p.categoria}: ${Number(p.monto_mensual).toFixed(2)} / mes
            </li>
          ))}
        </ul>
      </details>

      {cargando ? (
        <p>Cargando...</p>
      ) : (
        <>
          <h2>Total: ${total.toFixed(2)}</h2>

          <div className="resumen">
            <h3>Por categoría</h3>
            <ul>
              {resumen.map((r) => (
                <li key={r.categoria}>
                  {r.categoria}: ${Number(r.total).toFixed(2)}
                  {r.porcentaje_usado != null && (
                    <span className={claseDePorcentaje(r.porcentaje_usado)}>
                      {' '}
                      ({r.porcentaje_usado.toFixed(0)}% del presupuesto)
                    </span>
                  )}
                </li>
              ))}
            </ul>
          </div>

          <table>
            <thead>
              <tr>
                <th>Fecha</th>
                <th>Descripción</th>
                <th>Categoría</th>
                <th>Monto</th>
                <th>Estado</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {gastos.map((g) => (
                <tr key={g.id}>
                  <td>{new Date(g.fecha).toLocaleDateString()}</td>
                  <td>{g.descripcion}</td>
                  <td>{g.categoria}</td>
                  <td>${Number(g.monto).toFixed(2)}</td>
                  <td>{g.estado}</td>
                  <td>
                    <button
                      aria-label={`Marcar pagado ${g.descripcion}`}
                      onClick={() => marcarComoPagado(g.id)}
                      disabled={g.estado === 'pagado'}
                    >
                      Marcar pagado
                    </button>
                    <button
                      aria-label={`Eliminar ${g.descripcion}`}
                      onClick={() => eliminarGasto(g.id)}
                    >
                      Eliminar
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  )
}
