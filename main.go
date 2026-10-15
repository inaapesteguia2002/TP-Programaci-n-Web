package main

import (
	"context"      // Manejo de tiempos límite (timeouts) y cancelación de peticiones
	"database/sql" // Interfaz genérica para interactuar con bases de datos relacionales
	"fmt"          // Entrada y salida de texto (impresión en consola)
	"log"          // Registro de eventos y manejo de errores críticos (Fatalf)
	"net/http"     // Creación del servidor web y enrutamiento de peticiones HTTP
	"os"           // Interacción con el sistema (lectura de variables de entorno)
	"time"         // Manejo de duraciones y fechas (los 5 segundos del timeout)

	_ "github.com/jackc/pgx/v5/stdlib" // Driver de PostgreSQL (se registra internamente)

	"DeRabona/DataBase" // Repositorio de consultas generado por sqlc
	"DeRabona/Handler"  // Capa de Presentación (Handlers HTTP)
	"DeRabona/Logic"    // Capa de Lógica de Negocio (Services)
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:tu_contraseña@localhost:5432/deRabona_db?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Error al abrir la conexion: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("La base de datos no esta respondiendo: %v", err)
	}
	fmt.Println("Conexion a la base de datos exitosa.")

	// 1. Instanciamos el repositorio de sqlc
	repo := DataBase.New(db)

	// 2. Instanciamos los Services inyectándoles el repo
	equipoService := &Logic.EquiposService{Repo: repo}
	ligaService := &Logic.LigasService{Repo: repo}
	partidoService := &Logic.PartidosService{Repo: repo}
	tablaService := &Logic.TablaService{Repo: repo}

	// 3. Instanciamos los Handlers inyectándoles los services
	equipoHandler := &Handler.EquiposHandler{Service: equipoService}
	ligaHandler := &Handler.LigasHandler{Service: ligaService}
	partidoHandler := &Handler.PartidosHandler{Service: partidoService}
	tablaHandler := &Handler.TablaHandler{Service: tablaService}

	// 4. Enrutamiento HTTP oficial con los nuevos handlers
	// Equipos
	http.HandleFunc("/equipos", equipoHandler.HandleEquipos)
	http.HandleFunc("/equipos/", equipoHandler.HandleEquipoByID)

	// Ligas
	http.HandleFunc("/ligas", ligaHandler.HandleLigas)
	http.HandleFunc("/ligas/", ligaHandler.HandleLigaByID)

	// Partidos
	http.HandleFunc("/partidos", partidoHandler.HandlePartidos)
	http.HandleFunc("/partidos/", partidoHandler.HandlePartidoByID)

	// Tabla de Posiciones
	http.HandleFunc("/tabla", tablaHandler.HandleTabla)
	http.HandleFunc("/tabla/", tablaHandler.HandleTablaByLiga)

	// 5. Servidor de archivos estáticos (Frontend)
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
