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

	// 3. Instanciamos los Handlers inyectándoles los services correspondientes
	equipoHandler := &Handler.EquiposHandler{Service: equipoService}
	ligaHandler := &Handler.LigasHandler{Service: ligaService}
	partidoHandler := &Handler.PartidosHandler{
		PartidoService: partidoService,
		LigaService:    ligaService,
		EquipoService:  equipoService,
	}
	tablaHandler := &Handler.TablaHandler{
		TablaService:  tablaService,
		LigaService:   ligaService,
		EquipoService: equipoService,
	}

	// 4. Creamos el Multiplexor (ServeMux) y configuramos las rutas
	mux := http.NewServeMux()

	// Equipos
	mux.HandleFunc("/equipos", equipoHandler.HandleEquipos)
	mux.HandleFunc("/equipos/", equipoHandler.HandleEquipos)

	// Ligas
	mux.HandleFunc("/ligas", ligaHandler.HandleLigas)
	mux.HandleFunc("/ligas/", ligaHandler.HandleLigas)

	// Partidos
	mux.HandleFunc("/partidos", partidoHandler.HandlePartidos)
	mux.HandleFunc("/partidos/", partidoHandler.HandlePartidos)

	// Tabla de Posiciones
	mux.HandleFunc("/tabla", tablaHandler.HandleTabla)

	// 5. Servidor de archivos estáticos y redirección raíz
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Redirigir la raíz ("/") directamente a "/ligas"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/ligas", http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	})

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
