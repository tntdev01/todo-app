#!/bin/sh

kubectl exec \
  todo-app-db-0 \
  -n todo-app-namespace \
  -- psql \
  -U todo-app-db-user \
  -d todo-app-db \
  -c "
SELECT 'CREATE DATABASE todo_app-test-db'
WHERE NOT EXISTS (
  SELECT FROM pg_database
  WHERE datname='todo_app-test-db'
)\gexec
"