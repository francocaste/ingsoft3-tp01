# --- Los dos servicios de PREPROD, corriendo la imagen del TP7 ---

variable "image_repo_api" {
  default = "ghcr.io/francocaste/ingsoft3-tp01-backend" # SIN la etiqueta
}

variable "image_repo_front" {
  default = "ghcr.io/francocaste/ingsoft3-tp01-frontend" # SIN la etiqueta
}

# Preprod corre la imagen de UN commit concreto, escrita en el codigo. Es el commit
# de la release v7.0.0 (el que verifico QA y se aprobo en PROD). Para correr otra:
# cambiar esta linea por PR y recrear el servicio (destroy/apply o -replace), porque
# Render no deja actualizar un servicio del plan gratuito.
variable "image_tag" {
  default = "sha-ed172f47447a7346a2864c774e9db041f158fe0f"
}

variable "region" {
  default = "oregon"
}

locals {
  # La cadena de conexion se ARMA con lo que devolvio Neon: nadie la copia ni la pega.
  # La contrasena es un atributo sensible; urlencode por si trae caracteres reservados.
  database_url = "postgres://${neon_role.preprod.name}:${urlencode(neon_role.preprod.password)}@${neon_project.preprod.database_host}/${neon_database.preprod.name}?sslmode=require"
}

resource "render_web_service" "api" {
  name   = "gastos-api-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_api # sin etiqueta
      tag       = var.image_tag      # la etiqueta va en su propio campo
    }
  }

  # No se declara GIT_SHA: el commit que informa /health viene horneado en la imagen
  # (argumento de build del pipeline) y una variable del servicio lo pisaria.
  env_vars = {
    DATABASE_URL = { value = local.database_url }
  }
}

resource "render_web_service" "front" {
  name   = "gastos-front-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_front
      tag       = var.image_tag
    }
  }

  env_vars = {
    BACKEND_URL  = { value = render_web_service.api.url } # sale de la api, sin copiar nada
    DNS_RESOLVER = { value = "8.8.8.8" }
  }
}

output "api_url" {
  value = render_web_service.api.url
}

output "front_url" {
  value = render_web_service.front.url
}
