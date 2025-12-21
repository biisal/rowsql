-- name: ListTables :many
SELECT
  table_schema::text AS table_schema,
  table_name::text   AS table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
  AND table_schema NOT IN ('pg_catalog', 'information_schema')
ORDER BY table_schema, table_name;


-- name: GetTableFields :many
SELECT
  column_name::text AS column_name,
  data_type::text   AS data_type
FROM information_schema.columns
WHERE table_name::text = $1
ORDER BY ordinal_position;

