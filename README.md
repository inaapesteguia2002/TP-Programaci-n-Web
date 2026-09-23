# Trabajo Práctico 3: Desarrollo de API REST y Capa de Logica

## Descripción del Proyecto
Este proyecto implementa una **API REST completa en Go** utilizando arquitectura modular con separación de responsabilidades. Se conecta a una base de datos PostgreSQL 16 mediante consultas seguras generadas con `sqlc` y un driver optimizado (`pgx/v5`). 
El sistema gestiona las entidades principales de nuestra web: **Ligas, Equipos, Partidos y Tabla de Posiciones**

## Componentes principales
- `main.go`: Configura la conexión a la base de datos con un timeout de 5 segundos, inyecta el repositorio en la capa de lógica y levanta el servidor HTTP nativo (`net/http`).
- `Logic/`: Contiene los handlers (`equipos`, `ligas`, `partidos`, `tabla`) y los DTOs para validar las peticiones JSON de forma limpia.
- `DataBase/`: Modelos y código SQL generado automáticamente por `sqlc`.
- `requests.hurl`: Pruebas funcionales automatizadas para validar todos los endpoints.

## Endpoints de la API
- **Equipos**: 
  - `GET /equipos` / `POST /equipos`
  - `GET /equipos/{id}` / `PUT /equipos/{id}` / `DELETE /equipos/{id}`
- **Ligas**: 
  - `GET /ligas` / `POST /ligas`
  - `GET /ligas/{id}` / `PUT /ligas/{id}` / `DELETE /ligas/{id}`
- **Partidos**: 
  - `GET /partidos` / `POST /partidos`
  - `GET /partidos/{id}` / `PUT /partidos/{id}` / `DELETE /partidos/{id}`
- **Tabla de Posiciones**: 
  - `POST /tabla` (registrar equipo en la tabla)
  - `GET /tabla/{liga_id}` (ver posiciones ordenadas con JOIN a equipos)
  - `PUT /tabla-equipo` (actualizar puntos y estadísticas)
  - `DELETE /tabla-equipo` (quitar un equipo de la tabla)

## Instrucciones de ejecución y testing
Para evaluar este trabajo práctico, el proyecto cuenta con un script de automatización (Makefile) que se encarga de preparar todo el entorno de forma transparente sin requerir configuraciones manuales.

## Requisitos
- Go: version 1.20 o superior
- Docker y Docker Compose
- sqlc
- Make instalado en el sistema
- Hurl (para correr las pruebas)

## Pasos para probar el proyecto
1. Clonar el repositorio (>git:Clone y luego se pega la url del repositorio) y abrir una terminal en el directorio raiz.
2. Posicionarse en la rama correspondiente a la entrega:
   ```bash
   git fetch origin
   git checkout tp2
   ```
3. Levanta la base de datos PostgreSQL 16 y la API en contenedores aislados:
   ```bash
   docker compose up --build
   ```