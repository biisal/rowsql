<h1 align="center">RowSQL</h1>

<p align="center">
  <a href="https://opensource.org/licenses/MIT"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" /></a>
</p>

RowSQL is a simple tool to view and manage your databases. It provides a clean web interface to look at your information, edit records, and track your activity without needing complex desktop software.

![Demo](./resources/demo.gif)

## Features

- **Easy Viewing**: See your data in a clear, organized table.
- **Direct Editing**: Add, update, or delete records using simple forms.
- **Activity Log**: See a history of the changes you've made.
- **Filtering**: Find exactly what you need with easy search and sort options.
- **Multiple Databases**: Works with SQLite, MySQL, and PostgreSQL.

## Quick Install

Run this command in your terminal to install RowSQL:

```bash
curl -fsSL https://raw.githubusercontent.com/biisal/rowsql/refs/heads/main/install | bash
```

## How to use

Follow these steps to get started:

1.  **Start RowSQL**: Type `rowsql` in your terminal. 
2.  **Setup Configuration**: If it's your first run, RowSQL will ask to create a default configuration file. Type `y` and press Enter.
3.  **Find your config**: The settings are stored in a file named `config.json` in your home directory:
    - **Unix/Mac**: `~/.rowsql/config.json`
    - **Windows**: `%USERPROFILE%\.rowsql\config.json`
4.  **Connect your database**: Open `config.json` and add your database details. Here is a full example with multiple databases:
    ```json
    {
      "connections": [
        {
          "port": 8000,
          "db_string": "my_local_data.db"
        },
        {
          "port": 8081,
          "db_string": "postgres://admin:password@localhost:5432/customers"
        }
      ],
      "disable_auto_update": false,
      "max_items_per_page": 50,
      "log_file_path": "/home/user/.rowsql/rowsql.log"
    }
    ```
5.  **Open in browser**: Visit `http://localhost:8000` to start managing your data.

### Configuration Details

Here is what each setting in your `config.json` does:

- **`connections`**: A list of your database settings. You can add more than one database here.
  - **`port`**: The number you type into your browser to see your data (for example, `8000`).
  - **`db_string`**: The location of your database (either a file path or a web address).
- **`disable_auto_update`**: If set to `true`, RowSQL will stop checking for new versions automatically.
- **`max_items_per_page`**: The maximum number of rows to show on one page (default is 100).
- **`log_file_path`**: The file where RowSQL saves error messages if something goes wrong.

### Connection Examples

RowSQL uses a "connection string" (`db_string`) to find your database. Here are common formats:

- **SQLite**: Use the path to your database file (for example, `my_data.db`).
- **PostgreSQL or MySQL**: Use a standard connection address (for example, `postgres://user:password@localhost:5432/my_db`).

## Contributing

We welcome your help! If you find a bug or have a suggestion, please open an issue or submit a pull request.

For technical details and how to build from source, see [DEVELOPMENT.md](./DEVELOPMENT.md).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

<img src="./frontend/public/logo.png" alt="RowSQL Logo" height="100">
