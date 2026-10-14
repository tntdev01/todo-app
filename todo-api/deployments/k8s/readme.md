# This document covers the 8KS commands

### Get a list of contexts
```
kubectl config get-contexts
```

### Get current context
```
kubectl config current-context
```

### Switch current context
```
kubectl config use-context <context-name>

e.g.
kubectl config use-context docker-desktop
```

### Remove context
```
kubectl config delete-context '<context-name>'
```

### Remove cluster
```
kubectl config unset clusters.<cluster-name>
```

### Remove user
```
kubectl config unset users.<user-name>
```

### Get nodes for the current context
```
kubectl get nodes
```

### Apply change
```
kubectl apply -f <yaml-file-path>
```

### Get current namespace
```
kubectl config view --minify | grep namespace
```

### Switch namespace
```
kubectl config set-context --current --namespace=<namespace>
```

### Get secrets
```
kubectl get secret <secret-name>
```

### Get secrets
```
kubectl get secret <secret-name>
kubectl get secret <secret-name> -o json
kubectl get secret <secret-name> -o json | jq -r '(.data // {})'
kubectl get secret <secret-name> -o json | jq -r '(.data // {}) | to_entries[]
kubectl get secret <secret-name> -o json | jq -r '(.data // {}) | to_entries[] | "\(.key)=\(.value | @base64d)"'
```

### Port forwarding
```
kubectl port-forward -n <namespace> service/<database-service> 5432:5432

e.g.
kubectl port-forward -n todo-app-namespace service/todo-app-db-service 5432:5432
```

### To build an image for the api
```
docker build -t todo-api:latest .
```

### Install IGINX Ingress Controller
```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

### Verify Ingress is installed
```
kubectl get pods -n ingress-nginx
```

