# 🚀 Go REST API - Task Management

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-005C84?style=for-the-badge&logo=mysql&logoColor=white)
![Gorilla Mux](https://img.shields.io/badge/Gorilla_Mux-FF9E0F?style=for-the-badge&logo=go&logoColor=white)

## 📌 Descripción

Esta es una API RESTful escalable y eficiente construida con **Go (Golang)**. El proyecto implementa operaciones CRUD completas (Crear, Leer, Actualizar y Eliminar) para un sistema de gestión de tareas.

La arquitectura incluye un enrutamiento robusto manejado a través de `gorilla/mux` y persistencia de datos real conectada a una base de datos relacional **MySQL**. Este proyecto demuestra el manejo de peticiones HTTP, serialización/deserialización de formato JSON, y la ejecución segura de consultas SQL mediante sentencias preparadas.

## 🛠️ Tecnologías Utilizadas

- **Backend:** Go (Golang)
- **Enrutador:** `gorilla/mux`
- **Base de Datos:** MySQL
- **Driver DB:** `go-sql-driver/mysql`

## 🔌 Endpoints de la API

La API se ejecuta localmente en el puerto `3000` (`http://localhost:3000`).

| Método   | Ruta          | Descripción                                             |
| :------- | :------------ | :------------------------------------------------------ |
| `GET`    | `/`           | Ruta de bienvenida.                                     |
| `GET`    | `/tasks`      | Obtiene la lista de todas las tareas.                   |
| `GET`    | `/tasks/{id}` | Obtiene una tarea específica por su ID.                 |
| `POST`   | `/tasks`      | Crea una nueva tarea.                                   |
| `PUT`    | `/tasks/{id}` | Actualiza el nombre y contenido de una tarea existente. |
| `DELETE` | `/tasks/{id}` | Elimina una tarea por su ID.                            |

### Ejemplo de Payload (POST / PUT)

Al enviar datos a la API, asegúrate de usar formato JSON en el cuerpo (Body) de la petición:

```json
{
  "Name": "Aprender Arquitectura Backend",
  "Content": "Entender la diferencia entre Unmarshal y Encode en Go."
}
```

⚙️ Configuración e Instalación Local

1. Requisitos Previos
   Tener instalado Go.

Tener un servidor MySQL corriendo (ej. XAMPP). Nota: Por defecto, este proyecto apunta al puerto 3308.

### 2. Configurar la Base de Datos y Variables de Entorno

**Base de Datos:**
Crea una base de datos llamada `sistema` y ejecuta el siguiente script SQL:

```sql
CREATE DATABASE IF NOT EXISTS sistema;
USE sistema;

CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    content TEXT NOT NULL
);
```

Variables de Entorno (.env):
Por seguridad, las credenciales no están incluidas en el código fuente. Debes crear tu propio archivo de configuración:

En la raíz del proyecto, copia el archivo de ejemplo:
cp .env.example .env

Abre el nuevo archivo .env y configura tus credenciales locales de MySQL (usuario, contraseña y puerto).

3. Clonar y Ejecutar
   Abre tu terminal y ejecuta los siguientes comandos:

```
git add .
Bash

# Clonar el repositorio

git clone [https://github.com/tu-usuario/Go-REST-API.git](https://github.com/tu-usuario/Go-REST-API.git)

# Entrar a la carpeta del proyecto

cd Go-REST-API

# Instalar las dependencias (gorilla/mux y el driver de MySQL)

go mod tidy

# Iniciar el servidor

go run main.go

```

```
👨‍💻 Autor:
Johan Sebastian Vasquez Diaz - Ingeniero de Sistemas / Full-Stack Developer
```
