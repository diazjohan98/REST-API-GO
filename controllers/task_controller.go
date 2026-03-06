package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"Go-REST-API/config"
	"Go-REST-API/models"

	"github.com/gorilla/mux"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	conexion := config.ConectarBD()
	defer conexion.Close()

	registros, err := conexion.Query("SELECT id, name, content FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer registros.Close()

	var tareasBD models.AllTasks // Usamos el alias de la carpeta models

	for registros.Next() {
		var t models.Task // Usamos el Struct de la carpeta models
		err := registros.Scan(&t.ID, &t.Name, &t.Content)
		if err != nil {
			continue
		}
		tareasBD = append(tareasBD, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tareasBD)
}

// CreateTask maneja la petición POST
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &newTask)

	conexion := config.ConectarBD()
	defer conexion.Close()

	resultado, err := conexion.Exec("INSERT INTO tasks (name, content) VALUES (?, ?)", newTask.Name, newTask.Content)
	if err != nil {
		http.Error(w, "Error al insertar en la base de datos", http.StatusInternalServerError)
		return
	}

	ultimoID, _ := resultado.LastInsertId()
	newTask.ID = int(ultimoID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// GetTask maneja la petición GET para un solo ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	conexion := config.ConectarBD()
	defer conexion.Close()

	var t models.Task
	err = conexion.QueryRow("SELECT id, name, content FROM tasks WHERE id = ?", taskID).Scan(&t.ID, &t.Name, &t.Content)

	if err != nil {
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

// UpdateTask maneja la petición PUT
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var updatedTask models.Task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &updatedTask)

	conexion := config.ConectarBD()
	defer conexion.Close()

	resultado, err := conexion.Exec("UPDATE tasks SET name = ?, content = ? WHERE id = ?", updatedTask.Name, updatedTask.Content, taskID)
	if err != nil {
		http.Error(w, "Error al actualizar", http.StatusInternalServerError)
		return
	}

	filasAfectadas, _ := resultado.RowsAffected()
	if filasAfectadas == 0 {
		http.Error(w, "Tarea no encontrada para actualizar", http.StatusNotFound)
		return
	}

	updatedTask.ID = taskID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// DeleteTask maneja la petición DELETE
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	conexion := config.ConectarBD()
	defer conexion.Close()

	resultado, err := conexion.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		http.Error(w, "Error al eliminar", http.StatusInternalServerError)
		return
	}

	filasAfectadas, _ := resultado.RowsAffected()
	if filasAfectadas == 0 {
		http.Error(w, "Tarea no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "La tarea con ID %v fue eliminada exitosamente", taskID)
}
