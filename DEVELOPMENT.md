# Development and Technical Details

This document contains technical information for developers who want to contribute to RowSQL or build it from source.

## Prerequisites

- **Go** 1.25.3 or higher
- **Node.js** and **pnpm** (for building the frontend)
- **Air** (optional, for hot-reloading during development)

## Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/biisal/rowsql.git
   cd rowsql
   ```

2. Install dependencies:
   ```bash
   make install
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

## Development Mode

For development with hot-reloading:

```bash
make dev
```

This runs both the frontend dev server and backend with Air for automatic reloading.

Alternatively, run them separately:
- Frontend only: `make frontend-dev`
- Backend only: `make backend-dev`

> **Note:** When using development tools like `make dev` or `air`, ensure your `config.json` contains only **one** connection. If you have multiple connections, RowSQL will prompt you to select one, which does not work correctly with `air`.

## Configuration

The configuration file is located in the `.rowsql` folder in your home directory:
- **Unix/Mac**: `~/.rowsql/config.json`
- **Windows**: `%USERPROFILE%\.rowsql\config.json`

### Full Example

```json
{
  "connections": [
    {
      "port": 8000,
      "db_string": "my_local_data.db",
      "env": "development"
    }
  ],
  "disable_auto_update": true,
  "max_items_per_page": 100,
  "min_items_per_page": 10,
  "log_file_path": "/path/to/rowsql.log"
}
```

## Performance Considerations

### Database Size Recommendations

- **SQLite**: Works well with databases up to several GB. Performance may degrade with very large tables (100M+ rows).
- **PostgreSQL/MySQL**: Suitable for databases of any size, but consider the following:
  - Large result sets (>10,000 rows) are paginated automatically.
  - Complex queries on large tables may take time; use filters to narrow results.

### Best Practices

1. **Index your tables**: Ensure frequently queried columns have appropriate indexes for faster data retrieval.
2. **Limit result sets**: When working with large tables, use the built-in filtering and sorting features to limit the data loaded.
3. **Connection pooling**: RowSQL uses efficient connection pooling, but avoid running multiple instances against the same database unnecessarily.
4. **Memory usage**: RowSQL loads data on-demand. For very large result sets, pagination prevents excessive memory consumption.
5. **Network latency**: For remote databases (PostgreSQL/MySQL), ensure good network connectivity for optimal performance.
