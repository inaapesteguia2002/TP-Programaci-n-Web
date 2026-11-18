# Trabajo Práctico N° 5: Migración a Server-Side Rendering (SSR) con Templ

## Descripción del Proyecto
Este proyecto evoluciona el sistema **"De Rabona"** migrando su arquitectura desde una API REST con respuestas JSON hacia un enfoque de **Server-Side Rendering (SSR)** en Go. Se utilizan componentes HTML fuertemente tipados mediante **`templ`**, estilos con **Pico.css** complementados por una hoja de estilos propia (`static/styles.css`), y **formularios HTML clásicos (`GET` y `POST`)**.

El sistema cumple estrictamente con la restricción de **cero archivos de JavaScript (`.js`)**, resolviendo toda la interacción mediante envíos de formularios estándar (`application/x-www-form-urlencoded`) y aplicando el patrón **Post/Redirect/Get (PRG)** con redirecciones `303 See Other` para evitar reenvíos accidentales y mantener la interfaz sincronizada. La persistencia se gestiona en PostgreSQL mediante consultas tipadas generadas con **`sqlc`** y transacciones SQL, manteniendo la arquitectura en capas (`Handler`, `Logic`, `DataBase`, `Views`).

## Componentes principales
- `Views/`: Contiene los componentes `.templ` (`layout.templ`, `entity_form.templ`, `entity_list.templ`) que estructuran visualmente el panel de Ligas, Equipos, Partidos y Tabla de Posiciones, compilados nativamente a código Go.
- `Handler/`: Capa de presentación HTTP que procesa la vista principal (`GET /`) y las peticiones de formularios clásicos (`POST`), aplicando el patrón Post/Redirect/Get (PRG) y propagando errores de validación a la vista cuando corresponde.
- `Logic/`: Capa de lógica de negocio y servicios que validan las reglas del dominio (inscripción de equipos, cálculo automático de puntos al finalizar o eliminar partidos, y manejo de transacciones).
- `DataBase/`: Repositorio de acceso a datos generado automáticamente por `sqlc` a partir del esquema y consultas SQL de PostgreSQL.
- `static/`: Directorio de recursos estáticos que aloja la hoja de estilos personalizada (`styles.css`).
- `main.go`: Punto de entrada de la aplicación encargado de la conexión a PostgreSQL, inyección de dependencias, registro de rutas en `http.ServeMux` y servidor de archivos estáticos.

## Requisitos Previos
- **Go**: versión `1.27.0` (según lo especificado en `go.mod`)
- **Docker y Docker Compose** (para la ejecución de PostgreSQL y contenedores de prueba)
- **Herramientas CLI**: `templ`, `sqlc` y `hurl` instaladas en el entorno
- **Make**
- **Navegador Web Moderno** (Chrome, Firefox, Edge, Opera GX, etc.)

## Flujo de uso del Panel de Administración
1. **Ligas**: Crear una competición indicando su nombre y país, o eliminar ligas existentes.
2. **Equipos**: Registrar clubes indicando su nombre, o eliminarlos del sistema.
3. **Tabla de Posiciones (Inscripción y Consulta)**: Inscribir equipos en una liga específica y seleccionar una competición mediante el selector (`GET /?liga_id=...`) para visualizar la tabla ordenada por puntos.
4. **Partidos**: Registrar encuentros seleccionando la liga, equipo local, equipo visitante, goles y estado (`pendiente` o `finalizado`). Si el partido se registra como `finalizado`, el sistema actualiza automáticamente los partidos jugados, ganados, empatados, perdidos y puntos en la tabla de posiciones correspondiente.

## Instrucciones de ejecución

1. **Clonar el repositorio y posicionarse en la rama `tp5`**:
   ```bash
   git fetch origin
   git checkout tp5
