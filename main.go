package main

import (
	"fmt"
	"log"
	"net/http"

	"Go-REST-API/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Dejamos la ruta de bienvenida aquí porque es muy sencilla
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi REST API con Arquitectura por Capas")
}

func main() {
	// 1. Instanciamos el enrutador principal
	router := mux.NewRouter().StrictSlash(true)

	// 2. Registramos la ruta base
	router.HandleFunc("/", indexRoute)

	// 3. Delegamos el registro de las rutas a nuestra capa de Enrutamiento
	routes.SetTaskRoutes(router)

	// 4. Configuramos el Middleware de CORS
	// allowAll permite peticiones de cualquier puerto (ideal para desarrollo local)
	c := cors.AllowAll()

	// Envolvemos nuestro enrutador principal dentro del portero CORS
	handlerConCORS := c.Handler(router)

	// 5. Levantamos el servidor usando el nuevo handler
	fmt.Println("Servidor backend corriendo en el puerto 3000 con CORS habilitado...")
	log.Fatal(http.ListenAndServe(":3000", handlerConCORS))
}
