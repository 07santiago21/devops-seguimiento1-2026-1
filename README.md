# DevOps - Seguimiento #1 (2026-1)

API REST construida en Go, desplegada en AWS Lambda con API Gateway, base de datos PostgreSQL en RDS y acceso SSH mediante bastion host. Toda la infraestructura se gestiona con Terraform.

---

## Stack tecnológico

| Capa | Tecnología |
|------|-----------|
| Lenguaje | Go |
| Router HTTP | gorilla/mux |
| ORM | GORM |
| Base de datos | PostgreSQL 15 (AWS RDS) |
| Cómputo | AWS Lambda (`provided.al2023`) |
| API | AWS API Gateway REST (proxy) |
| Red | VPC, subnets públicas/privadas, bastion EC2 |
| IaC | Terraform >= 1.6 |
| Testing | Testify, go-sqlmock |

---

## Arquitectura

```
Internet
   │
   ▼
API Gateway  ──────────────────►  Lambda (Go)
                                      │
                     ┌────────────────┘
                     ▼
              RDS PostgreSQL (subred privada)
                     ▲
              Bastion EC2 (subred pública)
                     ▲
                   SSH
                     ▲
              Tu máquina local
```

- **API Gateway** recibe todas las peticiones HTTP y las reenvía a Lambda vía proxy (`/` y `/{proxy+}`, método `ANY`).
- **Lambda** ejecuta el binario Go. Detecta el entorno Lambda con `AWS_LAMBDA_RUNTIME_API` y usa el adaptador `gorillamux`.
- **RDS** está en subred privada, sin IP pública. Solo accesible desde Lambda y desde el bastion.
- **Bastion** es un EC2 en subred pública que permite abrir un túnel SSH hacia RDS.

---

## Endpoints de la API

### Estudiantes

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/students` | Crear estudiante |
| `GET` | `/students` | Listar todos |
| `GET` | `/students/{id}` | Obtener por ID |
| `PATCH` | `/students/{id}` | Actualización parcial |
| `PUT` | `/students/{id}` | Actualización completa |
| `DELETE` | `/students/{id}` | Eliminar |
| `POST` | `/api/v2/students` | Crear estudiante (v2) |

### Cursos

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/courses` | Crear curso |
| `GET` | `/courses` | Listar todos |
| `GET` | `/courses/{id}` | Obtener por ID |
| `PATCH` | `/courses/{id}` | Actualización parcial |
| `PUT` | `/courses/{id}` | Actualización completa |
| `DELETE` | `/courses/{id}` | Eliminar |

### Matrículas

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/enrollments` | Matricular estudiante |
| `GET` | `/enrollments` | Listar todas |
| `GET` | `/enrollments/{id}` | Obtener por ID |
| `PATCH` | `/enrollments/{id}` | Actualización parcial |
| `PUT` | `/enrollments/{id}` | Actualización completa |
| `DELETE` | `/enrollments/{id}` | Eliminar |

---

## Requisitos previos

- [Go](https://go.dev/) >= 1.21
- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.6
- [AWS CLI](https://aws.amazon.com/cli/) configurado con credenciales
- [psql](https://www.postgresql.org/download/) (cliente PostgreSQL)
- Par de llaves SSH en `~/.ssh/id_rsa` y `~/.ssh/id_rsa.pub`

Para generar el par de llaves si no lo tienes:
```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa
```

---

## Despliegue

### 1. Configurar variables de Terraform

```bash
cp terraform/terraform.tfvars.example terraform/terraform.tfvars
```

Edita `terraform/terraform.tfvars` y define al menos:

```hcl
database_password = "TuContraseñaSegura123!"
aws_access_key    = "AKIA..."      # o usa aws configure
aws_secret_key    = "..."
```

> `terraform.tfvars` está en `.gitignore` — nunca lo subas al repositorio.

La llave SSH pública (`~/.ssh/id_rsa.pub`) se detecta automáticamente. Terraform la importa a AWS y la asocia al bastion.

### 2. Compilar el binario Lambda

El runtime `provided.al2023` requiere un binario Linux llamado `bootstrap` empaquetado en un ZIP.

**Con Make (Git Bash / WSL / Linux):**
```bash
make build
```

**Sin Make (Git Bash directo):**
```bash
mkdir -p dist
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bootstrap .
cd dist && zip function.zip bootstrap && cd ..
```

**Sin Make (PowerShell):**
```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"
go build -o dist/bootstrap .
Push-Location dist; Compress-Archive -Path bootstrap -DestinationPath function.zip -Force; Pop-Location
```

Verifica que el artefacto existe:
```bash
ls -lh dist/function.zip
```

### 3. Desplegar con Terraform

```bash
cd terraform
terraform init
terraform plan
terraform apply   # escribe 'yes' para confirmar
```

Duración aproximada: **8-12 minutos** (RDS tarda más).

### 4. Obtener los outputs

```bash
cd terraform
terraform output
```

Ejemplo de salida:
```
api_invoke_url    = "https://9awbstv2oj.execute-api.us-east-1.amazonaws.com/prod"
bastion_public_ip = "3.89.87.222"
rds_endpoint      = "devops-seguimiento3-prod-postgres.xxxx.us-east-1.rds.amazonaws.com"
database_name     = "academia"
database_username = "postgres"
```

---

## Acceso SSH a la base de datos

RDS está en subred privada. Para accederlo se usa un túnel SSH a través del bastion.

### Verificar que el bastion responde

```bash
ssh -i ~/.ssh/id_rsa -o StrictHostKeyChecking=no ec2-user@<BASTION_PUBLIC_IP> "echo SSH OK"
```

Respuesta esperada: `SSH OK`

### Abrir el túnel SSH

Deja esta terminal abierta mientras trabajas con la base de datos:

```bash
ssh -i ~/.ssh/id_rsa \
    -L 5433:<RDS_ENDPOINT>:5432 \
    ec2-user@<BASTION_PUBLIC_IP> \
    -N -o StrictHostKeyChecking=no
```

El cursor queda en blanco — eso es correcto, el túnel está activo.

### Conectarse con psql

En otra terminal:

```bash
psql -h 127.0.0.1 -p 5433 -U postgres -d academia
```

| Parámetro | Valor | Descripción |
|-----------|-------|-------------|
| `-h` | `127.0.0.1` | Tu máquina local (entrada del túnel) |
| `-p` | `5433` | Puerto local del túnel |
| `-U` | `postgres` | Usuario de la base de datos |
| `-d` | `academia` | Nombre de la base de datos |

Si ves `academia=>` la conexión es exitosa.

---

## Scripts SQL

Los scripts están en la carpeta `sql/`. Con el túnel activo, ejecutarlos desde la raíz del proyecto:

### Crear tablas

```bash
psql -h 127.0.0.1 -p 5433 -U postgres -d academia -f sql/create.sql
```

Crea: `students`, `courses`, `enrollments` con sus índices y foreign keys.

Output esperado:
```
CREATE TABLE
CREATE TABLE
CREATE TABLE
CREATE INDEX
CREATE INDEX
```

### Insertar datos de prueba y validar

```bash
psql -h 127.0.0.1 -p 5433 -U postgres -d academia -f sql/insert_select.sql
```

Inserta 5 estudiantes, 4 cursos y 5 matrículas. Luego ejecuta SELECTs de validación.

### Eliminar tablas

```bash
psql -h 127.0.0.1 -p 5433 -U postgres -d academia -f sql/drop.sql
```

Elimina `enrollments`, `students` y `courses` en orden (respeta constraints).

### Comandos útiles dentro de psql

```sql
\dt                     -- listar tablas
\d students             -- ver estructura de una tabla
SELECT * FROM students;
SELECT * FROM courses;
SELECT * FROM enrollments;
SELECT COUNT(*) FROM students;
\q                      -- salir
```

---

## Probar la API

Con la infraestructura desplegada, reemplaza `<API_URL>` con el valor de `api_invoke_url`:

```bash
# Listar estudiantes
curl https://<API_URL>/students

# Crear estudiante
curl -X POST https://<API_URL>/students \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan","last_name":"García","age":21}'

# Listar cursos
curl https://<API_URL>/courses
```

---

## Variables de entorno (ejecución local)

Para correr la API localmente (sin Lambda), crea un archivo `.env` en la raíz:

```env
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=tu_password
DATABASE_NAME=academia
```

```bash
go run main.go
# Servidor disponible en http://localhost:8000
```

---

## Troubleshooting

| Problema | Causa | Solución |
|----------|-------|----------|
| `make: command not found` | `make` no instalado | Instala con `choco install make` o usa los comandos directos |
| `filebase64sha256` falla en destroy | ZIP no existe | Normal — el destroy funciona igual |
| `Connection refused` en psql | Túnel no está abierto | Verifica que el comando SSH con `-N` está corriendo |
| `Permission denied (publickey)` | Llave incorrecta | Usa `~/.ssh/id_rsa` (la privada, no la `.pub`) |
| `Connection timed out` al SSH | Security group | Verifica `allowed_ssh_cidr` en tfvars |
| Lambda no conecta a RDS | Variables de entorno | Confirma que el `terraform apply` terminó sin errores |
| Warning `stage_name is deprecated` | Ya corregido | Usar la versión actual del código |

---

## Limpieza de infraestructura

Para destruir todos los recursos y evitar costos:

```bash
cd terraform
terraform destroy   # escribe 'yes' para confirmar
```

Elimina: VPC, subnets, bastion EC2, RDS, Lambda, API Gateway, IAM roles, Key Pair y Security Groups.

---

## Estructura del proyecto

```
.
├── Makefile                        # Build del artefacto Lambda
├── main.go                         # Entry point (Lambda + HTTP local)
├── go.mod / go.sum
├── sql/
│   ├── create.sql                  # Crear tablas
│   ├── drop.sql                    # Eliminar tablas
│   └── insert_select.sql           # Datos de prueba + consultas
├── internal/
│   ├── database/postgres.go        # Conexión PostgreSQL
│   ├── student/                    # Dominio, repositorio, servicio, handler
│   ├── course/                     # Dominio, repositorio, servicio, handler
│   └── enrollment/                 # Dominio, repositorio, servicio, handler
└── terraform/
    ├── main.tf                     # Orquestación de módulos
    ├── variables.tf
    ├── outputs.tf
    ├── terraform.tfvars.example    # Plantilla de variables
    └── modules/
        ├── network/                # VPC, subnets, bastion, security groups
        ├── database/               # RDS PostgreSQL
        └── compute/                # Lambda + API Gateway
```
