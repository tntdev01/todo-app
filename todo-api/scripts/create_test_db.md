# Create Test DB

We assume that the postgres_stateful.yaml file has already been applied, and there is already a DB pod running.  We will just connect to that pod to run some SQL commands to create a test DB.

### Connecting to the POD
```
kubectl exec -it <DB pod name> -n <namespace> -- psql -U <DB user> -d <DB name>

e.g.
kubectl exec -it todo-app-db-0 -n todo-app-namespace -- psql -U todo-app-db-user -d todo-app-db
```

### Create DB
```
CREATE DATABASE <DB name>

e.g.
CREATE DATABASE "todo-app-test-db";
```

### Create Test DB user
```
CREATE USER <New user name> WITH PASSWORD <New user password>;
```

### Grant new user access
```
GRANT ALL PRIVILEGES
ON DATABASE <DB name>
TO <New user name>;

GRANT ALL ON SCHEMA public TO <New user name>;
```
