# 🧰 Go API — Local Setup Guide

## 📁 Project Structure

```
.
├── cmd/
│   └── api/
│       ├── main.go
│       └── app.go
├── internal/
│   ├── controllers/
│   ├── services/
│   ├── repositories/
│   └── models/
├── .env
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🚀 Getting Started

### 1️⃣ Prerequisites

Make sure you have the following installed:

- [Go 1.21+](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/)

---

### 2️⃣ Environment Variables

Create a `.env` file in the root directory:

```bash
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=mydb
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

APP_PORT=8080
GIN_MODE=release
```

---

### 3️⃣ Docker Setup

Build and run the containers:

```bash
docker-compose up --build
```

This will:

- Start a **PostgreSQL** container.
- Start your **Go API** container.
- Automatically connect both services via Docker network.

---

### 4️⃣ Access the Application

- API runs at:  
  👉 **http://localhost:8080**

- PostgreSQL runs at:  
  👉 **localhost:5432**

You can connect with:

```bash
host=localhost
port=5432
user=myuser
password=mypassword
dbname=mydb
```

---

### 5️⃣ Database Health Check

Run this to verify the DB is ready:

```bash
docker exec -it postgres_db pg_isready -U myuser -d mydb -h localhost
```

If healthy, it will return:

```
mydb: accepting connections
```

---

### 6️⃣ Running Locally (Without Docker)

If you prefer to run Go directly:

```bash
# Start PostgreSQL manually or use Docker
export $(cat .env | xargs)
go run ./cmd/api
```

---

### 7️⃣ Common Commands

| Command | Description |
|----------|--------------|
| `docker-compose up` | Run all containers |
| `docker-compose down -v` | Stop and remove containers & volumes |
| `docker ps` | Check running containers |
| `docker logs <container_name>` | View container logs |
| `go run ./cmd/api` | Run Go app locally |
| `go mod tidy` | Sync dependencies |

---

### 🧩 Useful Tips

- To rebuild without cache:

  ```bash
  docker-compose build --no-cache
  ```

- To import CSV data manually:

  ```bash
  docker cp transactions.csv postgres_db:/tmp/
  docker exec -it postgres_db psql -U myuser -d mydb -c "\copy transactions FROM '/tmp/transactions.csv' DELIMITER ',' CSV HEADER;"
  ```

---

### ✅ Verification

Once containers are up:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status": "ok"}
```
