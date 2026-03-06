package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConectarBD() *sql.DB {
	// 1. Cargar las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: Error al cargar el archivo .env")
	}

	// 2. Leer las variables usando el paquete 'os' (Operating System)
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 3. String de Conexión (Data Source Name - DSN)
	// Formateamos la ruta dinámicamente inyectando las variables de entorno
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// 4. Abrir la conexión (Pool de conexiones)
	conexion, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error fatal al conectar a la BD:", err)
	}

	// 5. Ping (Health Check)
	// sql.Open realmente no se conecta en ese instante, solo prepara la configuración.
	// Usamos Ping() para forzar a Go a verificar si MySQL realmente está vivo y respondiendo.
	err = conexion.Ping()
	if err != nil {
		log.Fatal("La base de datos no responde (Ping fallido):", err)
	}

	fmt.Println("Conexión exitosa a MySQL en el puerto", dbPort)

	return conexion
}
