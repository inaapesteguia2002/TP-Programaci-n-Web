# Trabajo Práctico 1: Mi Primera Aplicación Web

Aplicación: De Rabona (Resultados y seguimiento de partidos de fútbol en tiempo real).

# Definición del Dominio
La aplicación permitirá registrar y actualizar partidos de fútbol. Para representar esta información, la entidad principal será Partido, la cual almacenará los siguientes datos:
-Equipo Local(texto)
-Equipo Visitante(texto)
-Goles Local(número)
-Goles Visitante(número)
-Estado(texto - ej: "en vivo", "terminado", "programado")
-Fecha y Hora(fecha)

# Cómo ejecutar la aplicación:
1. Abre una terminal y posiciónate en el directorio raíz del proyecto (donde se encuentra `main.go`). Tener en cuenta que se debe tener una version Go 1.20 o superior.

2. Ejecuta el servidor con el siguiente comando: go run main.go
    Tener en cuenta: al ejecutar el comando debería aparecer: 
    " Servidor ESTÁTICO escuchando en http://localhost:8080
    Sirviendo archivos desde: ./static "

3. Ingresar desde un navegador web a http://localhost:8080

