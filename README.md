# Product Catalog Microservice

Microservicio de catálogo de productos agroecológicos desarrollado en Go. La idea principal del proyecto es practicar Domain-Driven Design (DDD) aplicando conceptos como agregados, entidades, objetos de valor, servicios de dominio y eventos.

## Objetivo del proyecto

- Practicar y afianzar DDD en un contexto realista (catálogo de productos de productores rurales).
- Diseñar límites claros entre capas (domain, application/handlers, infraestructura/repositorios).
- Proveer una API HTTP mínima para publicar y consultar productos.

## Stack técnico

- Go 1.21+
- Gin (HTTP framework)
- uuid (github.com/google/uuid)

## Estructura del proyecto

```
Product_Catalog_Microservice/
├─ cmd/
│  └─ app/
│     └─ main.go                 # bootstrap del servicio y wiring
├─ internal/
│  ├─ domain/
│  │  ├─ producto/
│  │  │  ├─ producto.go          # entidad/aggregate ProductoAgroecologico
│  │  │  ├─ valueobjects.go      # objetos de valor del producto
│  │  │  ├─ events.go            # eventos de dominio del producto
│  │  │  └─ infrastructure.go    # contratos/puertos desde dominio
│  │  ├─ productor/
│  │  │  ├─ productor.go         # entidad/aggregate Productor
│  │  │  ├─ valueobjects.go      # objetos de valor del productor
│  │  │  ├─ events.go            # eventos de dominio del productor
│  │  │  └─ infrastructure.go    # contratos/puertos desde dominio
│  │  └─ service/
│  │     └─ catalogoService.go   # servicio de dominio/orquestación
│  ├─ handlers/
│  │  └─ ProductoHandler.go      # controladores HTTP (Gin)
│  └─ repository/
│     ├─ ProductoRepostioryInterface.go  # repo en memoria para productos
│     └─ ProductorRepository.go          # repo en memoria para productores
├─ go.mod
├─ go.sum
└─ README.md
```

## Diseño de dominio (DDD)

### Agregados y raíces de agregado

- ProductoAgroecologico (Aggregate Root)
	- Identidad: ProductoID
	- Invariante: siempre asociado a un Productor existente, con categoría/tipo válidos y temporada consistente (inicio <= fin).

- Productor (Aggregate Root)
	- Identidad: ProductorID
	- Invariante: estados de verificación y actividad válidos; reputación en rango [0..5].

Nota: Catálogo se modela como un servicio de dominio que orquesta la publicación y consulta; no es un agregado persistente.

### Entidades

- Productor
	- Campos típicos: ID, Nombre, Ubicacion, EstadoVerificacion, EstadoActividad, Reputacion, PracticasCultivo.
	- Reglas: solo productores verificados/aptos pueden publicar; puede actualizar reputación y estado.

- ProductoAgroecologico
	- Campos típicos: ID, Nombre, Descripcion, Categoria, TipoProduccion, Temporada, EstadoDisponibilidad, Ubicacion, Imagen, ProductorID, PublicadoEn.
	- Reglas: puede marcarse como Excedente/Agotado; calcula disponibilidad por temporada/fecha.

### Objetos de valor (Value Objects)

- Del agregado ProductoAgroecologico
	- NombreProducto, DescripcionProducto
	- Categoria (p. ej., Fruta, Hortaliza, Tubérculo, PlantaMedicinal, Lácteo)
	- TipoProduccion (Agroecologico, Organico, Tradicional)
	- TemporadaLocal { Inicio, Fin }
	- EstadoDisponibilidad (Disponible, Agotado, Excedente)
	- Ubicacion { ZonaVeredal, Finca }
	- Imagen { URL, Descripcion }

- Del agregado Productor
	- NombreProductor
	- Ubicacion { ZonaVeredal, Finca }
	- EstadoVerificacion (Pendiente, EnProceso, Verificado)
	- EstadoActividad (Activo, Inactivo, Suspendido)
	- Reputacion (float32 [0..5])
	- PracticasCultivo (colección tipada)

### Eventos de dominio (ejemplos)

- ProductoPublicado, ProductoMarcadoComoExcedente, ProductoAgotado
- ProductorEnVerificacion, ProductorVerificado, ReputacionActualizada

## Endpoints (HTTP)

Los paths exactos pueden variar según el router, pero desde los handlers se desprenden los siguientes endpoints:

- POST /productos/publicar
	- Publica un nuevo producto.
	- Request JSON (ejemplo):
		```json
		{
			"productor_id": "123e4567-e89b-12d3-a456-426614174000",
			"nombre": "Tomate Orgánico",
			"descripcion": "Tomates frescos cultivados sin pesticidas.",
			"categoria": "Hortaliza",
			"tipo_produccion": "Agroecologica",
			"temporada_inicio": "2025-09-01",
			"temporada_fin": "2025-12-01",
			"zona_veredal": "Vereda El Paraíso",
			"finca": "Finca La Esperanza",
			"imagen_url": "https://ejemplo.com/tomate.jpg",
			"imagen_desc": "Tomates recién cosechados",
			"min_reputacion": 4.5
		}
		```

- POST /productos/excedente
	- Marca un producto como excedente en una fecha.
	- Request JSON (ejemplo):
		```json
		{
			"producto_id": "123e4567-e89b-12d3-a456-426614174000",
			"fecha": "2025-09-07"
		}
		```

- PUT /productos/disponibilidad
	- Recalcula/actualiza la disponibilidad según temporada y fecha.

- GET /catalogo (o similar)
	- Retorna el catálogo completo.

- GET /productos/listar (temporal)
	- Endpoint temporal para listar productos desde el repositorio en memoria.

## Cómo funciona (resumen de flujo)

1. El handler HTTP (Gin) recibe la petición y valida/transforma el JSON a los objetos de valor requeridos.
2. Se generan IDs (con uuid) en el controlador o en capa de dominio, según la política elegida.
3. El servicio de dominio `CatalogoService` orquesta la creación del agregado `ProductoAgroecologico`, aplica reglas de negocio y emite eventos.
4. El repositorio en memoria persiste el agregado (mapa protegido con mutex) y permite consultas simples.

## Ejecución local

Requisitos: Go 1.21+

```bash
go mod tidy
go run ./cmd/app
```

Luego invoca los endpoints con tu cliente HTTP favorito (curl, Postman, VS Code REST).

## Repositorios en memoria

- ProductoRepository: `map[ProductoID]*ProductoAgroecologico` con `sync.RWMutex`.
	- Métodos típicos: Save, GetByID, Update, GetAll, GetByCategoria, GetByEstado, GetByUbicacion, GetAvailableProducts, GetProductsInSeason, UpdateEstadoDisponibilidad.

- ProductorRepository: `map[ProductorID]*Productor` con `sync.RWMutex`.
	- Métodos típicos: Save, GetByID, Delete, GetAll, GetByUbicacion, GetVerificados, UpdateReputacion, UpdateEstadoVerificacion.

## Buenas prácticas DDD aplicadas

- Lógica de negocio en el dominio; handlers delgados.
- Objetos de valor para validar y encapsular invariantes.
- Servicios de dominio para orquestación (p. ej., `CatalogoService`).
- Eventos de dominio para estados relevantes.

## Desarrolladores

Añade aquí a los integrantes del equipo (reemplaza los placeholders):

- Julian David Meneses Daza (@Crownyny) — Desarrollador / juliandavidm@unicauca.edu.co

- Esteban Santiago Escandon Causaya (@santiago1scan) — Desarrollador / estebanescandon@unicauca.edu.co
- Miguel Angel Calambas Vivas (@maxskaink) — Desarrollador / mangelcvivas@unicauca.edu.co

Sugerencia: incluye a `@Crownyny` (owner del repo) y demás colaboradores.

## Próximos pasos

- Persistencia real (base de datos) y repositorios concretos.
- Cobertura de tests unitarios y de integración.
- Documentación OpenAPI/Swagger.
- Observabilidad (logs estructurados, métricas, tracing).
