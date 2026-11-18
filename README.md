# Trabajo Práctico N° 6: Interfaces Reactivas con HTMX y Templ

## Descripción del Proyecto
Este proyecto evoluciona el sistema **"De Rabona"** desde un enfoque de Server-Side Rendering (SSR) con recargas completas hacia una **interfaz web reactiva con experiencia de Single-Page Application (SPA)**. Se utilizan componentes HTML fuertemente tipados en Go mediante **`templ`**, estilos con **Pico.css** complementados por una hoja de estilos propia (`static/styles.css`), e **HTMX** para realizar actualizaciones parciales del DOM sin recargar la página.

El sistema cumple estrictamente con la restricción de **cero archivos de JavaScript (`.js`)** en el repositorio, resolviendo toda la interactividad mediante atributos declarativos de HTMX (`hx-post`, `hx-delete`, `hx-get`, `hx-target`, `hx-swap` y `hx-swap-oob`). En lugar de redirigir tras cada operación, el servidor devuelve únicamente los fragmentos de HTML actualizados (o respuestas vacías `200 OK` al eliminar filas) e incluye actualizaciones fuera de banda (*Out-of-Band Swaps*) para mantener sincronizados los menús desplegables en tiempo real. La persistencia se gestiona en PostgreSQL mediante consultas tipadas generadas con **`sqlc`** y transacciones SQL, manteniendo la arquitectura en capas (`Handler`, `Logic`, `DataBase`, `Views`).

## Componentes principales
- `Views/`: Contiene los componentes `.templ` (`layout.templ`, `entity_form.templ`, `entity_list.templ`) que cargan HTMX en el `<head>`, definen los formularios reactivos con `hx-post`, los selectores sincronizados con `hx-swap-oob="true"` y las tablas identificadas por ID con botones de borrado `hx-delete`.
- `Handler/`: Capa de presentación HTTP que procesa la carga inicial (`GET /`), las consultas parciales mediante la cabecera `HX-Request`, las creaciones reactivas (`POST /<entidades>`) devolviendo solo el fragmento de la lista actualizada, y las bajas (`DELETE /<entidades>/{id}`) devolviendo `200 OK` vacío para remover la fila del DOM.
- `Logic/`: Capa de lógica de negocio y servicios que validan las reglas del dominio (inscripción de equipos, cálculo automático de puntos al finalizar o eliminar partidos, y manejo de transacciones).
- `DataBase/`: Repositorio de acceso a datos generado automáticamente por `sqlc` a partir del esquema y consultas SQL de PostgreSQL.
- `static/`: Directorio exclusivo para recursos estáticos que aloja la hoja de estilos personalizada (`styles.css`), sin archivos `.js`.
- `main.go`: Punto de entrada de la aplicación encargado de la conexión a PostgreSQL, inyección de dependencias, registro de rutas (`GET`, `POST` y `DELETE`) en `http.ServeMux` y servidor de archivos estáticos.

## Requisitos Previos
- **Go**: versión `1.27.0` (según lo especificado en `go.mod`)
- **Docker y Docker Compose** (para la ejecución de PostgreSQL y contenedores de prueba)
- **Herramientas CLI**: `templ`, `sqlc` y `hurl` instaladas en el entorno
- **Make**
- **Navegador Web Moderno** (Chrome, Firefox, Edge, Opera GX, etc.)

## Flujo de uso del Panel de Administración
1. **Ligas**: Crear una competición indicando su nombre y país (`hx-post="/ligas"`). La tabla de ligas y todos los menús desplegables de la página se actualizan instantáneamente sin recargar. Permite eliminar ligas con confirmación (`hx-delete="/ligas/{id}"`), quitando la fila de la tabla en el acto.
2. **Equipos**: Registrar clubes indicando su nombre (`hx-post="/equipos"`), actualizando al instante el listado y los selectores de partidos e inscripciones, o eliminarlos dinámicamente (`hx-delete="/equipos/{id}"`).
3. **Tabla de Posiciones (Inscripción y Consulta)**: Inscribir equipos en una liga (`hx-post="/tabla"`), consultar las posiciones de cualquier competición sin mover la pantalla (`hx-get="/"` hacia `#lista-tabla`) o quitar equipos de una liga (`hx-delete="/tabla/{liga_id}/{equipo_id}"`).
4. **Partidos**: Registrar encuentros (`hx-post="/partidos"`) seleccionando la liga, equipo local, equipo visitante, goles y estado (`pendiente` o `finalizado`). Si el partido se registra como `finalizado`, el sistema actualiza automáticamente los partidos jugados, ganados, empatados, perdidos y puntos en la tabla de posiciones, o los revierte al borrar el partido (`hx-delete="/partidos/{id}"`).

## Instrucciones de ejecución

1. **Clonar el repositorio y posicionarse en la rama `tp6`**:
   ```bash
   git fetch origin
   git checkout tp6
