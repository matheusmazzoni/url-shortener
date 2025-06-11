# Go URL Shortener API 🚀

A RESTful URL shortener API built in Go, focusing on production-ready practices like structured logging, middleware, and Docker.

## ✨ Key Features

  * **Backend:** 100% Go.
  * **Storage:** SQLite, with the data layer abstracted via a `Storage` interface for easy replacement.
  * **API:** RESTful with endpoints to create (`POST /shorten`) and redirect (`GET /{shortKey}`) URLs.
  * **Architecture:** Follows the Standard Go Project Layout (`/cmd`, `/internal`).
  * **Observability:** Structured, contextual logging with `zerolog`, including a `request_id` for every request.
  * **Robustness:** Middleware chain for panic recovery and request logging.
  * **Reliability:** Graceful Shutdown implementation to prevent data loss and ensure clean exits.
  * **Configuration:** Flexible and 12-Factor App ready via environment variables.
  * **DevOps:**
      * Fully containerized with an optimized, multi-stage **Dockerfile**.
      * **Makefile** for automating common tasks (build, run, test, docker).

## 🏁 Getting Started

Follow the instructions below to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You will need the following tools installed on your machine:

  * [Go](https://go.dev/doc/install) (version 1.22 or higher)
  * [Docker](https://docs.docker.com/get-docker/)
  * [Make](https://www.gnu.org/software/make/) (typically pre-installed on Linux and macOS systems)

### Installation & Running

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/your-username/url-shortener.git
    cd url-shortener
    ```

2.  **Run the application:**
    There are two primary ways to run the project, thanks to the `Makefile`.

      * **Option A: Run locally with Go**
        This command will download dependencies and start the server. It's ideal for rapid development.

        ```sh
        make run
        ```

      * **Option B: Run with Docker (Recommended)**
        This command will build the Docker image and start a container. This method most closely resembles a production environment.

        ```sh
        make docker-run
        ```

    In either case, the API will be available at `http://localhost:8080`.

## ⚙️ API Usage

### 1\. Shorten a New URL

  * **Endpoint:** `POST /shorten`
  * **Body:** JSON containing the URL to be shortened.
  * **cURL Example:**
    ```sh
    curl -X POST -H "Content-Type: application/json" \
    -d '{"url": "https://github.com/your-username"}' \
    http://localhost:8080/shorten
    ```
  * **Success Response (201 Created):**
    ```json
    {
        "short_url": "http://localhost:8080/aB3xYz"
    }
    ```

### 2\. Redirect to the Original URL

  * **Endpoint:** `GET /{shortKey}`

  * **Example:**
    Access the `short_url` returned in the previous step directly in your browser:
    `http://localhost:8080/aB3xYz`

    You will be redirected to the original URL with a `301 Moved Permanently` status.

## 🔧 Configuration

The application is configured via environment variables, following the 12-Factor App principles.

| Variable        | Description                                                          | Default                  |
| :-------------- | :------------------------------------------------------------------- | :----------------------- |
| `SERVER_ADDRESS`| Address and port where the server will run.                          | `:8080`                  |
| `DB_PATH`       | Path to the SQLite database file.                                    | `./data/urls.db`         |
| `APP_BASE_URL`  | Public base URL used to construct the short URLs in the response.    | `http://localhost:8080`  |

## 🏗️ Project Structure

The project uses the standard Go community layout to separate concerns:

```
.
├── cmd/                # Application entry points (main.go)
│   └── url-shortener-api/
├── internal/           # Private application code (cannot be imported by other projects)
│   ├── api/            # Handlers, router, and middlewares
│   ├── config/         # Configuration logic
│   ├── shortener/      # Core business logic (key generation)
│   └── storage/        # Database interaction logic
├── Dockerfile          # Application container definition
├── Makefile            # Automation commands
└── go.mod              # Dependency management
```

## 🧰 Development Tools

The `Makefile` provides a series of commands to simplify the development workflow:

  * `make run`: Runs the application locally.
  * `make build`: Compiles the application binary.
  * `make test`: Runs all tests in the project.
  * `make tidy`: Tidies up go module dependencies.
  * `make docker-build`: Builds the application's Docker image.
  * `make docker-run`: Starts a container from the built image.
  * `make help`: Displays all available commands.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.
