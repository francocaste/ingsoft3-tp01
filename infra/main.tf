# TP8 - Infraestructura como Codigo.
#
# Este archivo declara PREPROD: un entorno NUEVO (base en Neon + api y front en Render)
# que no existe hasta que `terraform apply` lo crea y desaparece con `terraform destroy`.
# QA y PROD no se declaran aca ni se tocan: siguen siendo los del TP6/TP7.
#
# Las credenciales NO estan en el codigo: los providers las leen del entorno
# (NEON_API_KEY, RENDER_API_KEY, RENDER_OWNER_ID).

terraform {
  required_version = ">= 1.5"

  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "0.18.0" # fijada: el repo del provider esta archivado
    }
    render = {
      source  = "render-oss/render"
      version = "1.9.1"
    }
  }
}

provider "neon" {}   # lee NEON_API_KEY
provider "render" {} # lee RENDER_API_KEY y RENDER_OWNER_ID

# --- La base de PREPROD: un proyecto de Neon propio, que se puede destruir entero ---

variable "neon_org_id" {
  type        = string
  description = "Organizacion de Neon dueña del proyecto (org-...). Se pasa por TF_VAR_neon_org_id."
}

resource "neon_project" "preprod" {
  name   = "gastos-preprod"
  org_id = var.neon_org_id

  # Sin esta linea el primer apply falla: el default del provider (1 dia) supera el
  # tope del plan gratuito (6 horas).
  history_retention_seconds = 21600
}

resource "neon_role" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "gastos_preprod_owner"
}

resource "neon_database" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "gastos_preprod"
  owner_name = neon_role.preprod.name
}
