# Trabajo Práctico 4: Desarrollo de Interfaz de Usuario y Consumo de API REST

## Descripción del Proyecto
Este proyecto extiende la arquitectura del sistema **"De Rabona"** incorporando una **interfaz de usuario web (Frontend)** completa. Se desarrolló un Panel de Administración interactivo utilizando HTML5, CSS3 moderno (CSS Grid/Flexbox) y Vanilla JavaScript. 
La interfaz actúa como cliente de la API REST desarrollada en Go, consumiendo los endpoints (`/ligas`, `/equipos`, `/partidos`, `/tabla`) mediante peticiones `fetch` asíncronas. El sistema permite realizar operaciones CRUD completas, cargar selectores dinámicamente relacionando entidades por ID y actualizar el DOM en tiempo real sin recargar la página.

## Componentes principales
- `static/index.html`: Estructura semántica del panel de administración (Dashboard), dividida en tarjetas modulares para cada entidad.
- `static/styles.css`: Hojas de estilo responsivas que proveen una experiencia de usuario limpia, con efectos de estado y diseño adaptativo.
- `static/app.js`: Lógica del lado del cliente. Implementa el manejo de eventos del DOM, delegación de eventos, peticiones HTTP (`GET`, `POST`, `DELETE`) y el renderizado dinámico de la interfaz.
- `main.go`: Orquestador principal. Además de inyectar dependencias y conectarse a la base de datos, expone el servidor HTTP nativo (`net/http`) configurado para servir los archivos estáticos de la interfaz web.
- `Handler/`, `Logic/`, `DataBase/`: Capas de presentación, negocio y persistencia del backend (PostgreSQL + `sqlc`), que procesan las peticiones originadas por la nueva interfaz.

## Requisitos Previos
- **Go**: versión 1.20 o superior
- **Docker y Docker Compose** (para levantar la base de datos PostgreSQL)
- **Navegador Web Moderno** (Opera GX, Chrome, Edge, etc.) para visualizar la interfaz.

## Instrucciones de ejecución
Para probar la interfaz y el funcionamiento integral del sistema, sigue estos pasos:

1. **Clonar el repositorio** y abrir una terminal en el directorio raíz.
2. **Posicionarse en la rama correspondiente a la entrega (`tp4`)**:
   ```bash
   git fetch origin
   git checkout tp4
   ```
3. **Levantar la base de datos PostgreSQL**:
   ```bash
   docker-compose up -d
   ```
4. **Ejecutar el servidor backend en Go**:
   ```bash
   go run main.go
   ```
5. **Acceder a la Interfaz de Usuario**:
   - Abrir un navegador web.
   - Ingresar a **`http://localhost:8080/`**
   - Desde el panel, podrás interactuar visualmente con las Ligas, Equipos, Partidos y la Tabla de Posiciones, reflejándose los cambios directamente en la base de datos.
