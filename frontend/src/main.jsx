import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import './App.css'
import * as Sentry from '@sentry/react'

// El DSN de un front no es un secreto: solo sirve para ENVIAR datos, y viaja en el JavaScript
// que baja el navegador. En el back, en cambio, va por variable de entorno.
Sentry.init({
  dsn: 'https://656bd85da287391f6f34bcdf8811556d@o4512285309861888.ingest.us.sentry.io/4512285417603072',
  environment: window.location.hostname, // de que sitio vino: QA, produccion o localhost
  integrations: [Sentry.browserTracingIntegration()],
  tracesSampleRate: 1.0,
})

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
