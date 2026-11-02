# Trabajo Práctico N° 5: Migración a Server-Side Rendering con Templ y HTMX

## Descripción del Proyecto
Este proyecto evoluciona el sistema **"De Rabona"** migrando su arquitectura desde una API REST con cliente separado hacia un enfoque robusto de **Server-Side Rendering (SSR)**. Se utilizaron componentes tipados en Go mediante **`templ`** y **HTMX** para manejar la interactividad dinámica y asíncrona en el cliente.
El sistema cumple estrictamente con la restricción de **cero archivos de JavaScript (`.js`)**, delegando la generación del DOM y las actualizaciones parciales directamente en el servidor. Mantiene la persistencia robusta utilizando PostgreSQL, consultas tipadas con **`sqlc`** y una organización en capas (`Handler`, `Logic`, `DataBase`).

## Componentes principales
- `views/`: Contiene los componentes `.templ` que estructuran visualmente las Ligas, Equipos, Partidos y la Tabla de Posiciones, compilados nativamente a Go.
- `Handler/`: Capa de presentación HTTP que maneja peticiones clásicas y parciales mediante headers de HTMX (`HX-Request`), aplicando el patrón Post/Redirect/Get (PRG) cuando corresponde.
- `Logic/`: Capa de lógica de negocio y servicios que procesan las reglas del sistema de fútbol.
- `DataBase/`: Repositorios de datos generados automáticamente por `sqlc` a partir de esquemas SQL y PostgreSQL.
- `static/`: Directorio exclusivo para hojas de estilo y recursos visuales estáticos (sin lógica de scripts).
- `main.go`: Orquestador principal, encargado de la conexión a la base de datos, inyección de dependencias, configuración del enrutamiento HTTP y el servidor de archivos estáticos.

## Requisitos Previos
- **Go**: versión 1.20 o superior
- **Docker y Docker Compose** (para levantar la base de datos PostgreSQL)
- **Navegador Web Moderno** (Opera GX, Chrome, Edge, etc.) para visualizar la interfaz.
- **Herramientas CLI**: `templ` y `sqlc` instaladas en el entorno
- **Make**

## Instrucciones de ejecución
Para probar la interfaz y el funcionamiento integral del sistema, sigue estos pasos:

1. **Clonar el repositorio** y abrir una terminal en el directorio raíz.
2. **Posicionarse en la rama correspondiente a la entrega (`tp5`)**:
   ```bash
   git fetch origin
   git checkout tp5
   ```
3. **Ejecutar el sistema mediante el Makefile**:
   Este comando realiza automáticamente la limpieza de contenedores previos, levanta la base de datos, espera a que esté activa, ejecuta sqlc, compila los componentes de templ y levanta el servidor Go:
   ```bash
   make run-app
   ```
4. **Acceder a la Interfaz de Usuario**:
   - Abrir un navegador web.
   - Ingresar a **`http://localhost:8080/`**
   - Desde el panel, podrás interactuar visualmente con las Ligas, Equipos, Partidos y la Tabla de Posiciones, reflejándose los cambios directamente en la base de datos.
