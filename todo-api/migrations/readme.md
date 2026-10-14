# This document covers the migration commands

### create migration files
```
migrate create -ext sql -dir <migrate_dir> -seq <file_name>

e.g.
migrate create -ext sql -dir migrations -seq init_todos_table
```