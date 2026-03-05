package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task one",
		Content: "Some content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	// 1. Abrimos la conexión a la base de datos
	conexion := conectarBD()
	// Buena práctica en Go: defer asegura que la conexión se cierre al final de la función
	defer conexion.Close()

	// 2. Ejecutamos la consulta SQL
	registros, err := conexion.Query("SELECT id, name, content FROM tasks")
	if err != nil {
		// Si hay error en la BD, devolvemos un código 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer registros.Close()

	// 3. Creamos una variable vacía basada en tu estructura para guardar lo que llegue
	var tareasBD allTasks

	// 4. Recorremos fila por fila lo que nos devolvió MySQL
	for registros.Next() {
		var t task
		// Scan mapea las columnas de la BD directamente a las propiedades de tu struct
		err := registros.Scan(&t.ID, &t.Name, &t.Content)
		if err != nil {
			log.Println("Error al leer la fila:", err)
			continue
		}
		tareasBD = append(tareasBD, t)
	}

	// 5. Respondemos enviando el slice convertido a JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tareasBD)
}

func createTasks(w http.ResponseWriter, r *http.Request) {
	var newTask task

	// 1. Leer el cuerpo de la petición (Payload)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	// 2. Deserialización (Unmarshal)
	err = json.Unmarshal(reqBody, &newTask)
	if err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Conexión a la base de datos
	conexion := conectarBD()
	defer conexion.Close()

	// 4. Ejecución de la consulta SQL (Sentencia Preparada / Exec)
	resultado, err := conexion.Exec("INSERT INTO tasks (name, content) VALUES (?, ?)", newTask.Name, newTask.Content)
	if err != nil {
		http.Error(w, "Error al insertar en la base de datos", http.StatusInternalServerError)
		return
	}

	// 5. Obtener el ID autoincremental
	ultimoID, err := resultado.LastInsertId()
	if err != nil {
		http.Error(w, "Error al recuperar el ID", http.StatusInternalServerError)
		return
	}

	// 6. Actualizar nuestro objeto y preparar la respuesta
	newTask.ID = int(ultimoID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// 7. Serialización (Encode)
	json.NewEncoder(w).Encode(newTask)
}
func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	conexion := conectarBD()
	defer conexion.Close()

	var t task
	// QueryRow ejecuta la consulta y de una vez hace el Scan
	err = conexion.QueryRow("SELECT id, name, content FROM tasks WHERE id = ?", taskID).Scan(&t.ID, &t.Name, &t.Content)

	if err != nil {
		// sql.ErrNoRows es un error especial de Go para cuando el SELECT viene vacío
		if err == sql.ErrNoRows {
			http.Error(w, "Tarea no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func deleteTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	conexion := conectarBD()
	defer conexion.Close()

	// Ejecutamos el DELETE
	resultado, err := conexion.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		http.Error(w, "Error al eliminar la tarea", http.StatusInternalServerError)
		return
	}

	// Verificamos si realmente se borró algo
	filasAfectadas, _ := resultado.RowsAffected()
	if filasAfectadas == 0 {
		http.Error(w, "Tarea no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "La tarea con ID %v fue eliminada exitosamente de la base de datos", taskID)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	// 1. Capturamos el ID de la URL
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Leemos el cuerpo de la petición (Payload)
	var updatedTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	// 3. Deserializamos el JSON entrante
	err = json.Unmarshal(reqBody, &updatedTask)
	if err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	// 4. Abrimos la conexión a MySQL
	conexion := conectarBD()
	defer conexion.Close()

	// 5. Ejecutamos el UPDATE con una Sentencia Preparada
	resultado, err := conexion.Exec("UPDATE tasks SET name = ?, content = ? WHERE id = ?", updatedTask.Name, updatedTask.Content, taskID)
	if err != nil {
		http.Error(w, "Error al actualizar en la base de datos", http.StatusInternalServerError)
		return
	}

	// 6. Verificamos si realmente se actualizó algo
	filasAfectadas, _ := resultado.RowsAffected()
	if filasAfectadas == 0 {
		http.Error(w, "Tarea no encontrada para actualizar", http.StatusNotFound)
		return
	}

	// 7. Le asignamos el ID a la tarea actualizada y respondemos
	updatedTask.ID = taskID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi REST api")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTasks).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTasks).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// Función para conectar a la base de datos MySQL en el puerto 3308
func conectarBD() *sql.DB {
	conexion, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3308)/sistema")
	if err != nil {
		log.Fatal("Error de conexión a la BD:", err)
	}
	return conexion
}
