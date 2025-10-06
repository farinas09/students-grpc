# Sistema de Gestión de Estudiantes y Exámenes con gRPC

Este proyecto implementa un sistema distribuido de gestión de estudiantes y exámenes utilizando gRPC y Go, con PostgreSQL como base de datos.

## 🏗️ Arquitectura

El proyecto está compuesto por dos servicios gRPC independientes:

- **Servicio de Estudiantes** (`server-student`): Puerto 50051
- **Servicio de Exámenes** (`server-test`): Puerto 50052

Ambos servicios comparten la misma base de datos PostgreSQL y utilizan el patrón Repository para el acceso a datos.

## 📋 Características

### Servicio de Estudiantes
- ✅ Crear estudiantes (`SetStudent`)
- ✅ Obtener información de estudiantes (`GetStudent`)

### Servicio de Exámenes
- ✅ Crear exámenes (`SetTest`)
- ✅ Obtener información de exámenes (`GetTest`)
- ✅ Agregar preguntas a exámenes (`SetQuestions`) - Streaming
- ✅ Inscribir estudiantes a exámenes (`EnrollStudents`) - Streaming
- ✅ Obtener estudiantes inscritos en un examen (`GetStudentsPerTest`) - Streaming
- ✅ Realizar examen (`TakeTest`) - Streaming bidireccional

## 🛠️ Tecnologías Utilizadas

- **Go 1.24.3**
- **gRPC** para comunicación entre servicios
- **Protocol Buffers** para definición de APIs
- **PostgreSQL** como base de datos
- **Docker** para containerización de la base de datos
- **lib/pq** como driver de PostgreSQL

## 📁 Estructura del Proyecto

```
go-grpc/
├── database/                 # Configuración de base de datos
│   ├── docker-compose.yml   # Configuración de Docker Compose
│   ├── Dockerfile           # Imagen de PostgreSQL
│   ├── postgres.go          # Conexión y configuración de DB
│   └── up.sql              # Scripts de inicialización
├── models/                  # Modelos de datos
│   └── models.go           # Estructuras de datos
├── repository/              # Patrón Repository
│   └── repository.go       # Interfaces y implementaciones
├── server/                  # Lógica de servidores gRPC
│   ├── server.go           # Implementación de servicios
│   └── tests.go            # Tests unitarios
├── server-student/          # Servidor de estudiantes
│   └── main.go             # Punto de entrada del servicio
├── server-test/            # Servidor de exámenes
│   └── main.go             # Punto de entrada del servicio
├── client/                 # Cliente de ejemplo
│   └── main.go             # Cliente para probar los servicios
├── studentpb/              # Archivos generados de Protocol Buffers
│   ├── student.proto       # Definición del servicio de estudiantes
│   ├── student.pb.go       # Código Go generado
│   └── student_grpc.pb.go  # Servicios gRPC generados
└── testpb/                 # Archivos generados de Protocol Buffers
    ├── test.proto          # Definición del servicio de exámenes
    ├── test.pb.go          # Código Go generado
    └── test_grpc.pb.go     # Servicios gRPC generados
```

## 🚀 Instalación y Configuración

### Prerrequisitos

- Go 1.24.3 o superior
- Docker y Docker Compose
- Protocol Buffers compiler (`protoc`)

### 1. Clonar el repositorio

```bash
git clone https://github.com/farinas09/go-grpc.git
cd go-grpc
```

### 2. Instalar dependencias

```bash
go mod download
```

### 3. Levantar la base de datos

```bash
cd database
docker-compose up -d
```

### 4. Generar código de Protocol Buffers (si es necesario)

```bash
# Para estudiantes
protoc --go_out=. --go-grpc_out=. studentpb/student.proto

# Para exámenes
protoc --go_out=. --go-grpc_out=. testpb/test.proto
```

## 🏃‍♂️ Ejecución

### Ejecutar el servicio de estudiantes

```bash
go run server-student/main.go
```

El servicio estará disponible en `localhost:50051`

### Ejecutar el servicio de exámenes

```bash
go run server-test/main.go
```

El servicio estará disponible en `localhost:50052`

### Ejecutar el cliente de ejemplo

```bash
go run client/main.go
```

El cliente se conecta al servicio de exámenes y demuestra el uso de los métodos gRPC.

## 📊 Base de Datos

### Esquema

El sistema utiliza las siguientes tablas:

- **students**: Información de estudiantes
- **tests**: Información de exámenes
- **questions**: Preguntas asociadas a exámenes
- **enrollments**: Inscripciones de estudiantes a exámenes

### Datos de ejemplo

El script `up.sql` incluye datos de ejemplo:
- 4 estudiantes
- 5 exámenes
- 5 preguntas (una por examen)

## 🔧 Configuración

### Variables de entorno

- **Base de datos**: `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`
- **Puerto estudiantes**: `50051`
- **Puerto exámenes**: `50052`

### Modificar configuración

Para cambiar la configuración de la base de datos, edita las cadenas de conexión en:
- `server-student/main.go` (línea 20)
- `server-test/main.go` (línea 20)

## 📡 API gRPC

### Métodos Principales

#### TakeTest - Realizar Examen
El método `TakeTest` implementa un streaming bidireccional que permite:
- **Entrada**: El cliente envía respuestas (`TakeTestRequest` con `answer` y `test_id`)
- **Salida**: El servidor envía preguntas (`Question` con `id` y `question`)
- **Flujo**: El servidor envía una pregunta, espera la respuesta del cliente, y continúa con la siguiente pregunta

**Estructura del mensaje TakeTestRequest:**
```protobuf
message TakeTestRequest {
    string answer = 1;    // Respuesta del estudiante
    string test_id = 2;   // ID del examen a realizar
}
```

### Servicio de Estudiantes

```protobuf
service StudentService {
    rpc GetStudent(GetStudentRequest) returns (Student);
    rpc SetStudent(Student) returns (SetStudentResponse);
}
```

### Servicio de Exámenes

```protobuf
service TestService {
    rpc GetTest(GetTestRequest) returns (Test);
    rpc SetTest(Test) returns (SetTestResponse);
    rpc SetQuestions(stream Question) returns (SetQuestionResponse);
    rpc EnrollStudents(stream EnrollmentRequest) returns (SetQuestionResponse);
    rpc GetStudentsPerTest(GetStudentsPerTestRequest) returns (stream student.Student);
    rpc TakeTest(stream TakeTestRequest) returns (stream Question);
}
```

## 🐛 Troubleshooting

### Regenerar archivos Protocol Buffers

Si modificas los archivos `.proto`, regenera el código Go:

```bash
# Instalar dependencias de protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Regenerar archivos
protoc --go_out=. --go-grpc_out=. studentpb/student.proto
protoc --go_out=. --go-grpc_out=. testpb/test.proto
```

## 👨‍💻 Autor

**Erick Farinas**
- GitHub: [@farinas09](https://github.com/farinas09)

---

⭐ ¡No olvides darle una estrella al proyecto si te resulta útil!
