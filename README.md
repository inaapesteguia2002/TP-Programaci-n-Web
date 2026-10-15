# Trabajo Práctico 3: Desarrollo de API REST y Capa de Lógica

## Descripción del Proyecto
Este proyecto implementa una **API REST completa en Go** utilizando una arquitectura modular con separación de responsabilidades estricta (Handler - Service - Repository). Se conecta a una base de datos PostgreSQL 16 mediante consultas seguras generadas con `sqlc` y un driver optimizado (`pgx/v5`). 
El sistema gestiona las entidades principales del torneo: **Ligas, Equipos, Partidos y Tabla de Posiciones**.

## Componentes principales
- `main.go`: Configura la conexión a la base de datos con un timeout de 5 segundos, realiza la inyección de dependencias entre la capa de negocio (`Logic`), la de presentación (`Handler`) y el repositorio (`DataBase`), y levanta el servidor HTTP nativo (`net/http`) junto con archivos estáticos.
- `Logic/`: Capa de lógica de negocio que contiene los servicios y reglas de validación de cada entidad de forma pura.
- `Handler/`: Capa de presentación encargada de procesar las peticiones HTTP, parsear los JSON, validar códigos de estado y comunicarse con los servicios.
- `DataBase/`: Modelos y código SQL generado automáticamente por `sqlc`.
- `Makefile`: Automatización integral del ciclo de pruebas, contenedores, compilación y limpieza.

## Requisitos Previos
- **Go**: versión 1.20 o superior
- **Docker y Docker Compose**
- **sqlc**
- **Hurl** (herramienta para pruebas funcionales de los endpoints)

## Instrucciones de ejecución y testing

El proyecto cuenta con un automatizador en el `Makefile` que realiza todas las tareas necesarias de forma transparente (limpieza previa, levantamiento de contenedores, espera activa de la base de datos, compilación, ejecución de tests y limpieza posterior).

### Pasos para probar el proyecto:

1. **Clonar el repositorio** y abrir una terminal en el directorio raíz.
2. **Posicionarse en la rama correspondiente a la entrega (`tp3`)**:
   ```bash
   git fetch origin
   git checkout tp3
3. **Correr los test**
    ```bash
    make test-api