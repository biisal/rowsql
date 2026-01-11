# RowSQL

## A Multi SQL database management tool.

- RowSQL supports PostgreSQL, MySQL, and SQLite.

![](./resources/demo.gif)

## What does it do?

- [x] Represnt your sql data in a user-friendly UI.
- [x] Perform various SQL operations – Easily insert, update, or delete records using intuitive web forms.
- [x] Track query history – Maintain a "Recent Activity" log to see every change made to your database using RowSQL.
- [x] Inspect Table Structures
- [x] Shows rows the way you want them - you choose the order, filer

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

- The port number that RowSQL listens on.Using that port number, you can access RowSQL's web interface.

**LOG_FILE_PATH**

- The file path where RowSQL writes its Error logs.
