// Integracion: la API YA DESPLEGADA, con su base de verdad, sin navegador y sin dobles.
// Es el escalon del medio: ve lo que un test unitario con doble no puede ver, que es
// lo que contesta de verdad Postgres (tipos, fechas, columnas) a lo que le manda el codigo.
import { test, expect } from '@playwright/test'

const API = process.env.API_BASE_URL || 'http://localhost:8080'

// QA es una sola y la comparten todas las corridas: cada dato lleva la hora adentro
// (mas un sufijo al azar) para que ninguna otra corrida use el mismo nombre.
function nombreUnico(prefijo) {
  const hora = new Date().toISOString().replace(/\D/g, '').slice(0, 14)
  return `${prefijo}-${hora}-${Math.random().toString(36).slice(2, 6)}`
}

async function listar(request) {
  const res = await request.get(`${API}/api/gastos`)
  expect(res.status()).toBe(200)
  return res.json()
}

test.describe('API de gastos contra su base real', () => {
  // Red de seguridad: si una asercion falla a mitad de camino, igual se borra lo creado.
  const creados = []
  test.afterEach(async ({ request }) => {
    while (creados.length > 0) {
      await request.delete(`${API}/api/gastos/${creados.pop()}`) // 404 si ya se borro: esta bien
    }
  })

  test('alta, verificacion y borrado de un gasto', async ({ request }) => {
    const descripcion = nombreUnico('int-alta')

    const alta = await request.post(`${API}/api/gastos`, {
      data: { descripcion, monto: 12.5, categoria: 'Otros' },
    })
    expect(alta.status()).toBe(201)
    const creado = await alta.json()
    creados.push(creado.id)
    expect(creado).toMatchObject({ descripcion, monto: 12.5, categoria: 'Otros', estado: 'pendiente' })

    // La base lo guardo de verdad: el GET lo encuentra, con el monto intacto.
    const encontrado = (await listar(request)).find((g) => g.id === creado.id)
    expect(encontrado).toMatchObject({ descripcion, monto: 12.5 })

    const baja = await request.delete(`${API}/api/gastos/${creado.id}`)
    expect(baja.status()).toBe(204)

    // Y lo borro de verdad: ya no esta.
    expect((await listar(request)).find((g) => g.id === creado.id)).toBeUndefined()
  })

  test('un alta invalida se rechaza con 400 y no crea nada', async ({ request }) => {
    const sinDescripcion = await request.post(`${API}/api/gastos`, {
      data: { descripcion: '', monto: 10, categoria: 'Otros' },
    })
    expect(sinDescripcion.status()).toBe(400)
    expect(await sinDescripcion.text()).toContain('obligatorios')

    const categoriaInventada = await request.post(`${API}/api/gastos`, {
      data: { descripcion: nombreUnico('int-invalida'), monto: 10, categoria: 'Inventada' },
    })
    expect(categoriaInventada.status()).toBe(400)

    // No se creo nada: ni un gasto sin descripcion ni el de la categoria inventada.
    const lista = await listar(request)
    expect(lista.filter((g) => g.descripcion === '')).toHaveLength(0)
    expect(lista.filter((g) => g.descripcion.startsWith('int-invalida'))).toHaveLength(0)
  })

  // Tercera, elegida por mi: el cambio de estado. Es un UPDATE en la base, y la regla
  // "un gasto pagado no vuelve a pendiente" tiene que sobrevivir al viaje completo
  // (HTTP -> servicio -> SQL). Si esto se rompe, nadie puede registrar un pago.
  test('marcar como pagado persiste y no se puede volver a pendiente', async ({ request }) => {
    const descripcion = nombreUnico('int-estado')
    const alta = await request.post(`${API}/api/gastos`, {
      data: { descripcion, monto: 7, categoria: 'Otros' },
    })
    expect(alta.status()).toBe(201)
    const { id } = await alta.json()
    creados.push(id)

    const pagar = await request.patch(`${API}/api/gastos/${id}/estado`, { data: { estado: 'pagado' } })
    expect(pagar.status()).toBe(200)
    expect((await listar(request)).find((g) => g.id === id).estado).toBe('pagado')

    const revertir = await request.patch(`${API}/api/gastos/${id}/estado`, { data: { estado: 'pendiente' } })
    expect(revertir.status()).toBe(409)
    expect((await listar(request)).find((g) => g.id === id).estado).toBe('pagado')

    const baja = await request.delete(`${API}/api/gastos/${id}`)
    expect(baja.status()).toBe(204)
    expect((await listar(request)).find((g) => g.id === id)).toBeUndefined()
  })
})
