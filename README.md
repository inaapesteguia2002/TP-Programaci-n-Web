# Trabajo Práctico 2: Persistiendo el Dominio

## Descripción del Proyecto
Este proyecto implementa la capa de persistencia utilizando Go, PostgreSQL, Docker y la herramienta sqlc. Gestiona las entidades principales del sistema: Ligas, equipos, partidos y tabala de posiciones. 

## Análisis de la base de datos
1. Ligas: entidad que almacena el ID, nombre de la liga y el país al que pertenece.
2. Equipos: almacena el Id y el nombre de los clubes registrados. Tiene una restriccion 'UNIQUE' para evitar la duplicacion de equipos.
3. Tablas de Posiciones: es una tabla relacional entre ligas y equipos. Con la restriccion UNIQUE asegura que un equipo no pertenezca dos veces a una unica tabla. Las estadisticas se inicializan en DEFAULT 0. Implementamos 'ON DELETE CASCADE' para que si se elimina un equipo o una liga tambien lo hagan las estadisticas.
4. Partidos: identifica a las ligas y equipos(local y visitante) mediante FOREIGN KEYS. Con 'ON DELETE CASCADE' aseguramos que al borrar una liga se borren automaticamente todos sus partidos. El estado se inicializa con DEFAULT'pendiente' para que se vea el partido antes de que se juegue. 


## Instrucciones de ejecución y testing
Para evaluar este trabajo práctico, el proyecto cuenta con un script de automatización (Makefile) que se encarga de preparar todo el entorno de forma transparente sin requerir configuraciones manuales.

## Requisitos
* Go: version 1.20 o superior
* Docker y Docker Compose
* sqlc
* Make instalado en el sistema

## Pasos para probar el proyecto
1. Clonar el repositorio y abrir una terminal en el directorio raiz.
2. Posicionarse en la rama correspondiente a la entrega:
   ```bash
   git checkout tp2
   ```
3. En la terminal escribir el comando: make test
