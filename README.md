# Sistema de Votación en Go con Múltiples Bases de Datos

Este proyecto implementa un sistema de votación completo usando Go con soporte para **MySQL**, **PostgreSQL** y **MongoDB**. Incluye API REST, autenticación por tokens y documentación Swagger.

## 📋 Características

### Soporte Multi-Base de Datos

- ✅ **MySQL** (predeterminado)
- ✅ **PostgreSQL**
- ✅ **MongoDB**
- Cambio de base de datos mediante variable de entorno
- Interfaz unificada para todas las bases de datos

### API REST

- Registro de usuarios con contraseñas encriptadas (bcrypt)
- Login/Logout con autenticación por tokens
- Listado de candidatos
- Sistema de votación (un voto por usuario)
- Visualización de resultados
- Documentación interactiva con Swagger UI

### Interfaz Web

- Interfaz bilingüe (Inglés/Español)
- Bootstrap 5 responsive
- Registro e inicio de sesión
- Visualización de candidatos
- Sistema de votación intuitivo
- Resultados en tiempo real con gráficos

## Instalación Rápida

### Opción 1: Script Automático (Recomendado)

```bash
# Dar permisos de ejecución al script
chmod +x setup_database.sh

# Ejecutar el script
./setup_database.sh
```

El script te permitirá elegir:

1. MySQL
2. PostgreSQL
3. MongoDB
4. Todas las bases de datos

El script automáticamente:

- Detecta tu sistema operativo (macOS/Linux)
- Instala la(s) base(s) de datos seleccionada(s)
- Crea base de datos y usuario
- Instala dependencias de Go
- Crea archivo `.env` con la configuración
- Crea script `run.sh` para ejecutar la aplicación

### Opción 2: Instalación Manual

#### Prerequisitos

- Go 1.16+
- Una de las siguientes bases de datos:
  - MySQL 5.7+ / MariaDB 10.3+
  - PostgreSQL 12+
  - MongoDB 4.4+

#### Pasos

1. **Instalar dependencias de Go**

```bash
go mod init voting-system
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get go.mongodb.org/mongo-driver/mongo
go get golang.org/x/crypto/bcrypt
```

2. **Configurar base de datos** (ver sección de configuración específica abajo)

3. **Crear archivo .env**

```bash
DB_TYPE=mysql          # mysql, postgresql, o mongodb
DB_USER=voting_user
DB_PASSWORD=voting_password
DB_HOST=localhost
DB_PORT=3306          # 3306 para MySQL, 5432 para PostgreSQL, 27017 para MongoDB
DB_NAME=voting_system
```

4. **Ejecutar aplicación**

```bash
source .env
go run main.go
```

## 💾 Configuración por Base de Datos

### MySQL

**Instalación:**

```bash
# macOS
brew install mysql
brew services start mysql

# Ubuntu/Debian
sudo apt-get install mysql-server
sudo systemctl start mysql
```

**Configuración:**

```bash
mysql -u root -p

CREATE DATABASE voting_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'voting_user'@'localhost' IDENTIFIED BY 'voting_password';
GRANT ALL PRIVILEGES ON voting_system.* TO 'voting_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

**Variables de entorno:**

```bash
export DB_TYPE=mysql
export DB_PORT=3306
```

### PostgreSQL

**Instalación:**

```bash
# macOS
brew install postgresql@14
brew services start postgresql@14

# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**Configuración:**

```bash
# macOS
createdb voting_system
psql -d voting_system

# Linux
sudo -u postgres psql

CREATE DATABASE voting_system;
CREATE USER voting_user WITH PASSWORD 'voting_password';
GRANT ALL PRIVILEGES ON DATABASE voting_system TO voting_user;
\c voting_system
GRANT ALL ON SCHEMA public TO voting_user;
\q
```

**Variables de entorno:**

```bash
export DB_TYPE=postgresql
export DB_PORT=5432
```

### MongoDB

**Instalación:**

```bash
# macOS
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community

# Ubuntu/Debian
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt-get update
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
```

**Configuración:**

```bash
# MongoDB no requiere configuración adicional
# La base de datos se crea automáticamente
```

**Variables de entorno:**

```bash
export DB_TYPE=mongodb
export DB_PORT=27017
# DB_USER y DB_PASSWORD son opcionales para MongoDB sin autenticación
```

## 🏃 Ejecución

### Usando el script (si usaste setup_database.sh)

```bash
./run.sh
```

### Manual

```bash
# Cargar variables de entorno
source .env

# O exportar manualmente
export DB_TYPE=mysql
export DB_USER=voting_user
export DB_PASSWORD=voting_password
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=voting_system

# Ejecutar
go run main.go
```

La aplicación:

- Se conectará a la base de datos configurada
- Creará las tablas/colecciones automáticamente
- Poblará candidatos iniciales (Alice Johnson, Bob Smith, Charlie Brown)
- Iniciará el servidor en el puerto 8000

## 📚 Acceso a la Aplicación

### Swagger UI

http://127.0.0.1:8000/swagger/

### Interfaz Web

Abrir `index.html` en el navegador

### API Base URL

http://127.0.0.1:8000/api/

## 🔌 Endpoints de la API

### Autenticación

**Registrar Usuario:**

```bash
curl -X POST http://127.0.0.1:8000/api/register/ \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password123"}'
```

**Iniciar Sesión:**

```bash
curl -X POST http://127.0.0.1:8000/api/login/ \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'
```

**Cerrar Sesión:**

```bash
curl -X POST http://127.0.0.1:8000/api/logout/ \
  -H "Authorization: Token YOUR_TOKEN_HERE"
```

### Votación

**Listar Candidatos:**

```bash
curl -X GET http://127.0.0.1:8000/api/candidates/ \
  -H "Authorization: Token YOUR_TOKEN_HERE"
```

**Emitir Voto:**

```bash
curl -X POST http://127.0.0.1:8000/api/vote/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Token YOUR_TOKEN_HERE" \
  -d '{"candidate":1}'
```

**Ver Resultados:**

```bash
curl -X GET http://127.0.0.1:8000/api/results/ \
  -H "Authorization: Token YOUR_TOKEN_HERE"
```

## 🔄 Cambiar de Base de Datos

Para cambiar entre bases de datos, simplemente modifica la variable `DB_TYPE`:

```bash
# Cambiar a PostgreSQL
export DB_TYPE=postgresql
export DB_PORT=5432

# Cambiar a MongoDB
export DB_TYPE=mongodb
export DB_PORT=27017

# Cambiar a MySQL
export DB_TYPE=mysql
export DB_PORT=3306

# Ejecutar
go run main.go
```

O edita el archivo `.env`:

```bash
DB_TYPE=postgresql  # mysql, postgresql, o mongodb
```

## 🏗️ Arquitectura

### Estructura del Proyecto

```
voting-system/
├── main.go              # Código principal con soporte multi-DB
├── setup_database.sh    # Script de configuración automática
├── run.sh              # Script para ejecutar la aplicación
├── index.html          # Interfaz web bilingüe
├── .env                # Variables de entorno
├── go.mod              # Dependencias de Go
├── go.sum              # Checksums de dependencias
└── README.md           # Esta documentación
```

### Tecnologías

**Backend:**

- Go (Golang)
- MySQL Driver: `github.com/go-sql-driver/mysql`
- PostgreSQL Driver: `github.com/lib/pq`
- MongoDB Driver: `go.mongodb.org/mongo-driver/mongo`
- bcrypt para encriptación

**Frontend:**

- HTML5
- Bootstrap 5
- JavaScript (Vanilla)
- Sistema bilingüe (EN/ES)

## 📊 Esquemas de Base de Datos

### MySQL / PostgreSQL

**users**

```sql
id            INT/SERIAL PRIMARY KEY
username      VARCHAR(100) UNIQUE NOT NULL
email         VARCHAR(255) NOT NULL
password      VARCHAR(255) NOT NULL
has_voted     BOOLEAN DEFAULT FALSE
created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

**candidates**

```sql
id            INT/SERIAL PRIMARY KEY
name          VARCHAR(255) NOT NULL
created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

**votes**

```sql
id            INT/SERIAL PRIMARY KEY
user_id       INT NOT NULL (FK -> users.id)
candidate_id  INT NOT NULL (FK -> candidates.id)
created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
UNIQUE(user_id)
```

**tokens**

```sql
id            INT/SERIAL PRIMARY KEY
user_id       INT NOT NULL (FK -> users.id)
token         VARCHAR(255) UNIQUE NOT NULL
created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
INDEX(token)
```

### MongoDB

**Colecciones:**

- `users` - Usuarios con índice único en username
- `candidates` - Candidatos
- `votes` - Votos con índice único en user_id
- `tokens` - Tokens de autenticación con índice único en token

## 🐳 Docker Setup (Opcional)

### docker-compose.yml

```yaml
version: '3.8'

services:
  # MySQL
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: voting_system
      MYSQL_USER: voting_user
      MYSQL_PASSWORD: voting_password
    ports:
      - '3306:3306'
    volumes:
      - mysql_data:/var/lib/mysql

  # PostgreSQL
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: voting_system
      POSTGRES_USER: voting_user
      POSTGRES_PASSWORD: voting_password
    ports:
      - '5432:5432'
    volumes:
      - postgres_data:/var/lib/postgresql/data

  # MongoDB
  mongodb:
    image: mongo:6.0
    ports:
      - '27017:27017'
    volumes:
      - mongodb_data:/data/db

volumes:
  mysql_data:
  postgres_data:
  mongodb_data:
```

### Ejecutar con Docker

```bash
# Iniciar base de datos específica
docker-compose up -d mysql      # Solo MySQL
docker-compose up -d postgres   # Solo PostgreSQL
docker-compose up -d mongodb    # Solo MongoDB

# O todas
docker-compose up -d

# Ejecutar aplicación
export DB_TYPE=mysql  # o postgresql, mongodb
go run main.go
```

## 🧪 Testing Completo

### Script de prueba automatizado

```bash
#!/bin/bash

API_URL="http://127.0.0.1:8000"

echo "1. Registrando usuario..."
curl -s -X POST $API_URL/api/register/ \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"test123"}'

echo -e "\n2. Iniciando sesión..."
RESPONSE=$(curl -s -X POST $API_URL/api/login/ \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}')

TOKEN=$(echo $RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token: $TOKEN"

echo -e "\n3. Listando candidatos..."
curl -s -X GET $API_URL/api/candidates/ \
  -H "Authorization: Token $TOKEN"

echo -e "\n4. Votando..."
curl -s -X POST $API_URL/api/vote/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Token $TOKEN" \
  -d '{"candidate":1}'

echo -e "\n5. Viendo resultados..."
curl -s -X GET $API_URL/api/results/ \
  -H "Authorization: Token $TOKEN"

echo -e "\n6. Cerrando sesión..."
curl -s -X POST $API_URL/api/logout/ \
  -H "Authorization: Token $TOKEN"
```

## 🔐 Seguridad

- ✅ Contraseñas encriptadas con bcrypt
- ✅ Autenticación basada en tokens
- ✅ Tokens criptográficamente seguros
- ✅ Prepared statements (prevención SQL injection)
- ✅ Transacciones ACID para integridad
- ✅ Un voto por usuario (constraint)
- ✅ CORS habilitado para desarrollo

## 📈 Comparación de Bases de Datos

| Característica           | MySQL           | PostgreSQL             | MongoDB               |
| ------------------------ | --------------- | ---------------------- | --------------------- |
| Tipo                     | Relacional      | Relacional             | NoSQL                 |
| Transacciones            | ✅              | ✅                     | ✅                    |
| ACID                     | ✅              | ✅                     | ✅                    |
| Joins                    | ✅              | ✅                     | ❌                    |
| Esquema flexible         | ❌              | ❌                     | ✅                    |
| Rendimiento lectura      | Alto            | Alto                   | Muy alto              |
| Rendimiento escritura    | Alto            | Medio                  | Muy alto              |
| Escalabilidad horizontal | Media           | Media                  | Alta                  |
| Mejor para               | Web tradicional | Aplicaciones complejas | Big data, tiempo real |

## 🛠️ Troubleshooting

### Error de conexión

```bash
# Verificar que la BD esté corriendo
# MySQL
brew services list | grep mysql
sudo systemctl status mysql

# PostgreSQL
brew services list | grep postgresql
sudo systemctl status postgresql

# MongoDB
brew services list | grep mongodb
sudo systemctl status mongod
```

### Error de autenticación

```bash
# Verificar credenciales en .env
cat .env

# Recrear usuario (MySQL)
mysql -u root -p
DROP USER 'voting_user'@'localhost';
CREATE USER 'voting_user'@'localhost' IDENTIFIED BY 'voting_password';
```

### Puerto en uso

```bash
# Cambiar puerto en .env o matar proceso
lsof -ti:8000 | xargs kill -9
```

## 🎯 Próximas Mejoras

- [ ] Tests unitarios e integración
- [ ] Migrations automáticas
- [ ] Rate limiting
- [ ] Panel de administración
- [ ] Estadísticas avanzadas
- [ ] Exportación de resultados (CSV, PDF)
- [ ] Autenticación OAuth
- [ ] WebSockets para resultados en tiempo real
- [ ] CI/CD pipeline
- [ ] Kubernetes deployment

## 📝 Licencia

Este proyecto es de código abierto bajo la Licencia MIT.

## 👤 Autor

Desarrollado como proyecto educativo de sistema de votación con soporte multi-base de datos.

---

**¡Sistema de votación listo para usar! 🗳️**

Para comenzar:

1. Ejecuta `./setup_database.sh`
2. Sigue las instrucciones
3. Ejecuta `./run.sh`
4. Abre http://127.0.0.1:8000/swagger/
