# API de Bandas de Metal - Go HTTP
## Hugo Méndez - 241265

API REST construida con la librería estándar de Go para gestionar un catálogo de bandas de metal. No utiliza frameworks externos.

---

## Tema

Catálogo de bandas de metal extremo (Metalcore, Deathcore, Beatdown Hardcore). Cada banda tiene id, nombre, género, país de origen, año de formación y álbum destacado.

---

## Estructura del proyecto

```
.
├── main.go
├── data/
│   └── bandas.json
├── capturas/
│   └── delete.jpeg
│   └── elementos.jpeg
│   └── error.jpeg
│   └── filtros_combinados.jpeg
│   └── path_parameters.jpeg
│   └── persistencia_datos.jpeg
│   └── post.jpeg
│   └── puerto.jpeg
│   └── query_parametro.jpeg
├── Dockerfile
└── docker-compose.yml
```

---

## Ejecución

```bash
docker compose up --build
```

El servidor queda disponible en `http://localhost:41265`.

---

## Endpoints

### GET /api/ping
Verificar que el servidor está corriendo.

```
GET http://localhost:41265/api/ping
```

Respuesta:
```json
{ "message": "pong" }
```

---

### GET /api/bandas
Obtener todas las bandas. Soporta filtros combinados por query params.

```
GET http://localhost:41265/api/bandas
GET http://localhost:41265/api/bandas?id=1
GET http://localhost:41265/api/bandas?genero=deathcore
GET http://localhost:41265/api/bandas?pais_origen=suiza
GET http://localhost:41265/api/bandas?anio_desde=2014&anio_hasta=2019
GET http://localhost:41265/api/bandas?genero=deathcore&anio_desde=2014
```

Query params disponibles:

| Parámetro | Descripción |
|---|---|
| `id` | Filtrar por id exacto |
| `nombre` | Filtrar por nombre (contiene, sin distinción de mayúsculas) |
| `genero` | Filtrar por género (contiene) |
| `pais_origen` | Filtrar por país (contiene) |
| `anio_desde` | Año de formación mínimo |
| `anio_hasta` | Año de formación máximo |

Respuesta exitosa `200 OK`:
```json
[
  {
    "id": 5,
    "nombre": "Slaughter to Prevail",
    "genero": "Deathcore",
    "pais_origen": "Rusia / Estados Unidos",
    "anio_formacion": 2014,
    "album_destacado": "Kostolom"
  }
]
```

---

### GET /api/bandas/{id}
Obtener una banda por id usando path parameter.

```
GET http://localhost:41265/api/bandas/1
```

Respuesta exitosa `200 OK`:
```json
{
  "id": 1,
  "nombre": "Kublai Khan TX",
  "genero": "Metalcore / Beatdown",
  "pais_origen": "Estados Unidos",
  "anio_formacion": 2009,
  "album_destacado": "Absolute"
}
```

Error `404 Not Found`:
```json
{
  "error": "NOT_FOUND",
  "code": 404,
  "message": "No se encontró una banda con id 99"
}
```

---

### POST /api/bandas
Crear una nueva banda. Se persiste en `data/bandas.json`.

```
POST http://localhost:41265/api/bandas
Content-Type: application/json
```

Body:
```json
{
  "nombre": "Knocked Loose",
  "genero": "Metalcore / Hardcore",
  "pais_origen": "Estados Unidos",
  "anio_formacion": 2013,
  "album_destacado": "A Different Shade of Blue"
}
```

Respuesta exitosa `201 Created`:
```json
{
  "id": 11,
  "nombre": "Knocked Loose",
  "genero": "Metalcore / Hardcore",
  "pais_origen": "Estados Unidos",
  "anio_formacion": 2013,
  "album_destacado": "A Different Shade of Blue"
}
```

Campos requeridos: `nombre`, `genero`, `pais_origen`, `anio_formacion`, `album_destacado`.

---

### PUT /api/bandas/{id}
Reemplazar una banda completa por id.

```
PUT http://localhost:41265/api/bandas/1
Content-Type: application/json
```

Body (todos los campos requeridos):
```json
{
  "nombre": "Kublai Khan TX",
  "genero": "Metalcore",
  "pais_origen": "Estados Unidos",
  "anio_formacion": 2009,
  "album_destacado": "Absolute"
}
```

Respuesta exitosa `200 OK`: la banda actualizada.

---

### PATCH /api/bandas/{id}
Actualizar parcialmente una banda (solo los campos enviados).

```
PATCH http://localhost:41265/api/bandas/1
Content-Type: application/json
```

Body (solo los campos a modificar):
```json
{
  "album_destacado": "Napalmbuddhism"
}
```

Respuesta exitosa `200 OK`: la banda actualizada.

---

### DELETE /api/bandas/{id}
Eliminar una banda por id.

```
DELETE http://localhost:41265/api/bandas/1
```

Respuesta exitosa `200 OK`:
```json
{
  "message": "Banda eliminada correctamente"
}
```

---

## Errores estructurados

Todos los errores devuelven JSON con el mismo formato:

```json
{
  "error": "NOT_FOUND",
  "code": 404,
  "message": "No se encontró una banda con id 99"
}
```

| Código | error | Causa |
|---|---|---|
| 400 | `BAD_REQUEST` | Body inválido o campo faltante |
| 400 | `INVALID_ID` | El id en la URL no es un número |
| 404 | `NOT_FOUND` | No existe una banda con ese id |
| 405 | `METHOD_NOT_ALLOWED` | Método HTTP no soportado |

---

## Docker

El `Dockerfile` usa la imagen oficial `golang:1.22-alpine`. No requiere compilación previa — Go compila y ejecuta directamente con `go run main.go`.

```dockerfile
FROM golang:1.22-alpine
WORKDIR /app
CMD ["go", "run", "main.go"]
```

El volumen `.:/app` permite que los cambios en `data/bandas.json` persistan fuera del contenedor.
