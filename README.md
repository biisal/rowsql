<h1 align="center">RowSQL</h1>
<p align="center" style="display: flex; justify-content: center; align-items: center; gap: 10px">
  <a href="https://opensource.org/licenses/MIT"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" /></a>
  <a href="https://golang.org/"><img alt="Go Version" src="https://img.shields.io/badge/Go-1.25.3-00ADD8?style=flat-square&logo=go" /></a>
  <a href="https://reactjs.org/"><img alt="Made with React" src="https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react" /></a>
</p>

> A lightweight, cross-platform database management tool for PostgreSQL, MySQL, and SQLite with a modern web-based interface.

RowSQL provides an intuitive web UI to manage your databases without the complexity of heavy desktop applications. Perfect for developers who need quick access to their data with powerful filtering, editing, and query tracking capabilities.

![Demo](./resources/demo.gif)

## Features

- [x] Represent your sql data in a user-friendly UI.
- [x] Perform various SQL operations – Easily insert, update, or delete records using intuitive web forms.
- [x] Track query history – Maintain a "Recent Activity" log to see every change made to your database using RowSQL.
- [x] Inspect Table Structures
- [x] Shows rows the way you want them - you choose the order, filter

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/biisal/rowsql/refs/heads/main/install | bash
```

## How to Use?

1. Create a .env file in your home directory's rowsql folder:
   - Unix: ~/.rowsql/.env
   - Windows: %USERPROFILE%\\.rowsql\\.env

2. .env file should contain the following

```
DBSTRING=my.db
PORT=8000
LOG_FILE_PATH="~/.rowsql/rowsql.log"
```

3. run `rowsql`
4. open http://127.0.0.1:8000 in your browser

### .env Variables Explanation

**DBSTRING** is the connection string to your database.

- For **SQLite**, this is the path to your `.db` file.
- For **PostgreSQL**/**MySQL**, this is the PostgreSQL/MySQL connection string.
  - `DBSTRING=postgres://admin:admin@127.0.0.1:5432/admin_db`
- RowSQL automatically detects the database type from the connection string.

**PORT**

- The port number that RowSQL listens on. Using that port number, you can access RowSQL's web interface.

**LOG_FILE_PATH**

- The file path where RowSQL writes its Error logs.

## Development

### Prerequisites

- **Go** 1.25.3 or higher
- **Node.js** and **pnpm** (for building the frontend)
- **Air** (optional, for hot-reloading during development)

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/biisal/rowsql.git
   cd rowsql
   ```

2. Install frontend dependencies:
   ```bash
   cd frontend && pnpm install && cd ..
   ```

3. Build the project:
   ```bash
   make build
   ```
   This will build the frontend and compile the Go binary to `bin/rowsql`.

4. Run the binary:
   ```bash
   ./bin/rowsql
   ```

### Development Mode

For development with hot-reloading:

```bash
make dev
```

This runs both the frontend dev server and backend with Air for automatic reloading.

Alternatively, run them separately:
- Frontend only: `make frontend-dev`
- Backend only: `make backend-dev`

### Project Structure

```
rowsql/
├── cmd/server/          # Application entry point
├── internal/            # Internal packages
│   ├── database/        # Database connection & type handling
│   ├── router/          # HTTP routes & handlers
│   ├── service/         # Business logic
│   └── logger/          # Logging utilities
├── frontend/            # React + Vite frontend
├── configs/             # Configuration management
└── resources/           # Static resources
```

## Performance Considerations

### Database Size Recommendations

- **SQLite**: Works well with databases up to several GB. Performance may degrade with very large tables (100M+ rows).
- **PostgreSQL/MySQL**: Suitable for databases of any size, but consider the following:
  - Large result sets (>10,000 rows) are paginated automatically
  - Complex queries on large tables may take time; use filters to narrow results

### Best Practices

1. **Index your tables**: Ensure frequently queried columns have appropriate indexes for faster data retrieval.

2. **Limit result sets**: When working with large tables, use the built-in filtering and sorting features to limit the data loaded.

3. **Connection pooling**: RowSQL uses efficient connection pooling, but avoid running multiple instances against the same database unnecessarily.

4. **Memory usage**: RowSQL loads data on-demand. For very large result sets, pagination prevents excessive memory consumption.

5. **Network latency**: For remote databases (PostgreSQL/MySQL), ensure good network connectivity for optimal performance.

## Contributing

Contributions are welcome! Here's how you can help:

1. **Report bugs**: Open an issue describing the bug and steps to reproduce it.
2. **Suggest features**: Share your ideas for new features or improvements.
3. **Submit pull requests**: Fork the repo, make your changes, and submit a PR.

Please ensure your code follows the existing style and includes appropriate tests where applicable.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

<img src="./frontend/public/logo.png" alt="RowSQL Logo" height="100">
