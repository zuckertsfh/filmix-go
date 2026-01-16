# Filmix Backend

Backend API for Filmix application, built with **Go** and **Fiber**.

## 🚀 Getting Started

### Prerequisites
*   [Go 1.20+](https://go.dev/dl/)
*   [PostgreSQL](https://www.postgresql.org/download/)

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/senatroxx/filmix-backend.git
    cd filmix-backend
    ```

2.  **Setup Environment Variables**
    Copy `.env.example` to `.env` and adjust the database credentials.
    ```bash
    cp .env.example .env
    ```
    *Make sure your Postgres server is running and the database exists.*

3.  **Run Migrations & Seeder**
    Initialize your database schema and dummy data.
    ```bash
    # Run migrations
    go run main.go migrate up

    # Seed initial data (Users, Cinemas, Movies from TMDB)
    go run main.go seed
    ```

4.  **Run the Server**
    ```bash
    go run main.go serve
    ```
    The server will start at `http://localhost:8080` (or the port defined in `.env`).

5.  **Run with Live Reload (Optional)**
    If you have `air` installed for hot-reloading:
    ```bash
    air
    ```

---

## 🏗️ Project Structure & Flow

This project follows a **Scalable Monolith** structure (Clean Architecture inspired).
If you are coming from a **Mobile Development** background (Flutter/Android/iOS), here is an analogy map:

### 1. API Flow Lifecycle 🔄

1.  **Router (`internal/http/routes`)** 🚦
    *   *Analogi Mobile:* **Navigation Graph**.
    *   Defines the URL endpoints (e.g., `/api/v1/login`) and maps them to Handlers.

2.  **Handler (`internal/http/handlers`)** 🎮
    *   *Analogi Mobile:* **ViewModel / Controller**.
    *   Parses incoming JSON requests.
    *   Validates input.
    *   Calls the **Service**.
    *   Returns JSON response to the client.

3.  **Service (`internal/services`)** 🧠
    *   *Analogi Mobile:* **UseCase / Interactor**.
    *   Contains purely **Business Logic**.
    *   Example: "Check if password matches", "Calculate transaction total".
    *   Coordinates data flow between Handler and Repository.

4.  **Repository (`internal/repositories`)** 🗄️
    *   *Analogi Mobile:* **Repository / DAO**.
    *   Responsible for **Data Access** only (SQL Queries).
    *   The only layer allowed to talk to the Database.

5.  **Database (`internal/database`)** 💾
    *   *Analogi Mobile:* **Local DB (Room/SQLite)**.
    *   Contains Entities (Data Models) and Migrations.

### 2. Key Directories 📂

| Directory | Description | Mobile Analogy |
| :--- | :--- | :--- |
| `cmd/` | Application entry points (Main) | `Application` / `main()` |
| `internal/config/` | Environment configuration | `Config` / `Flavor` |
| `internal/http/handlers` | Request processors | `ViewModel` |
| `internal/http/middleware` | Pre-request logic (Auth, Logging) | `Interceptor` |
| `internal/services` | Business rules | `UseCase` |
| `internal/repositories` | Database queries | `RepositoryImpl` |
| `internal/database/entities`| Database structs | `DataModel` |

### 3. Development Workflow 🛠️

When adding a new feature (e.g., "Add Comment"), work **Inside-Out**:

1.  **Entity**: Define the `Comment` struct in `database/entities`.
2.  **Repository**: Create `CommentRepository` with `Create` method.
3.  **Service**: Create `CommentService` to handle business logic.
4.  **Handler**: Create `CommentHandler` to parse input.
5.  **Route**: Register the endpoint in `routes`.

---
