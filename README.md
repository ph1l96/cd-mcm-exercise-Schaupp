# Exercise 4: Vulnerability Scanning & Kubernetes Deployment

**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 24

| Exercise | Topic | Branch |
|----------|-------|--------|
| 1 | Git Basics: PRs, Interactive Rebase, Unit Tests | `exercise/01-git-basics` |
| 2 | Microservice Architecture, Docker & GitHub Actions | `exercise/02-microservice-docker` |
| 3 | CI Pipeline: SonarCloud, Matrix Builds, Linting | `exercise/03-ci-pipeline` |
| 4 | Vulnerability Scanning & Kubernetes Deployment | `exercise/04-security-k8s` |

- Integrate vulnerability scanning into the CI/CD pipeline
- Scan Docker images and Go dependencies for known vulnerabilities
- Deploy a multi-tier application to Kubernetes (Minikube)
- Understand Kubernetes concepts: Deployments, Services, Secrets, Probes

- **Language:** Go 1.24+
- **Web Framework:** Gorilla Mux
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **CI/CD:** GitHub Actions
- **Code Quality:** SonarCloud, golangci-lint
- **Security:** Trivy, govulncheck
- **Deployment:** Kubernetes (Minikube)

## Project: Product Catalog API

Throughout the four exercises you will build and evolve a **Product Catalog API** -- a RESTful web service for managing products (create, read, update, delete). The API is written in Go and grows in complexity with each exercise.

### What the Application Does

The Product Catalog API exposes the following HTTP endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/products` | List all products |
| POST | `/products` | Create a new product |
| GET | `/products/{id}` | Get a product by ID |
| PUT | `/products/{id}` | Update a product |
| DELETE | `/products/{id}` | Delete a product |

A product has three fields: `id`, `name`, and `price`.

### Project Structure

```
cmd/api/main.go                # Application entry point -- starts the HTTP server
internal/
  model/product.go             # Product data model and validation
  store/
    memory.go                  # In-memory store (Exercise 1-2)
    postgres.go                # PostgreSQL store (from Exercise 2)
  handler/handler.go           # HTTP request handlers (routing, JSON encoding)
Dockerfile                     # Multi-stage Docker build (from Exercise 2)
docker-compose.yml             # Orchestrates API + PostgreSQL (from Exercise 2)
.github/workflows/ci.yml       # CI/CD pipeline (from Exercise 2, extended in 3-4)
k8s/                           # Kubernetes manifests (Exercise 4)
```

### What You Build in Each Exercise

| Exercise | What You Do |
|----------|-------------|
| **1 -- Git Basics** | Fork the repo, write unit tests for the in-memory store, create your first Pull Request, and practice interactive rebase to clean up commit history. |
| **2 -- Microservice & Docker** | Understand the microservice architecture, complete a GitHub Actions CI pipeline with a Docker build job, analyze the Dockerfile and Docker Compose setup, and add HTTP handler tests. |
| **3 -- CI Pipeline** | Extend the pipeline with matrix builds (multiple Go versions and OS), integrate golangci-lint for code quality, set up SonarCloud for static analysis, and improve test coverage to ≥ 80%. |
| **4 -- Security & K8s** | Scan the Docker image with Trivy, scan Go dependencies with govulncheck, deploy the application to a local Kubernetes cluster (Minikube), and configure production-readiness features (probes, resource limits). |

- Completed Exercise 3 (CI pipeline with quality gates)
- Docker Desktop installed
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) installed
- [kubectl](https://kubernetes.io/docs/tasks/tools/) installed
- [Trivy](https://aquasecurity.github.io/trivy/) installed (optional, for local scanning)

## What's New in This Exercise

- **Kubernetes manifests** (`k8s/`) -- Deployment, Service, Secret, PVC
- **Trivy scanning** -- container image vulnerability scanning
- **Dependency scanning** -- Go module vulnerability checks
- **Complete CD pipeline** -- from code to running in Kubernetes

---

## Tasks

### Task 1: Matrix Builds (4 Points)

1. **Build the Docker image locally:**
   ```bash
   docker build -t product-catalog:latest .
   ```

2. **Scan the image with Trivy:**
   ```bash
   trivy image product-catalog:latest
   ```

3. **Analyze the results:**
   - How many vulnerabilities were found? Categorize by severity (CRITICAL, HIGH, MEDIUM, LOW).
   - Which base image contributes the most vulnerabilities?
   - Can you reduce vulnerabilities by changing the base image? Try switching to `scratch` or `distroless`.

4. **Add a Trivy scan job to the CI pipeline** (see the TODO in `ci.yml`) that:
   - Runs after the `docker-build` job
   - Scans the built Docker image using `aquasecurity/trivy-action@master`
   - Fails the build if CRITICAL or HIGH vulnerabilities are found
   - Outputs results in `table` format

   > **Hint:** The Trivy action needs `image-ref`, `format`, `exit-code`, and `severity` parameters.

5. **Upload the Trivy scan results as a build artifact:**
   - Generate a JSON report (use `format: 'json'` and `output` parameter)
   - Upload it using `actions/upload-artifact@v4`
   - Use `if: always()` so the report is uploaded even if the scan finds vulnerabilities

**Deliverable:** Trivy scan output (before and after base image optimization). Updated CI workflow. Trivy JSON report downloadable as artifact from the Actions run.

---

### Task 2: Vulnerability Scanning -- Dependencies (4 Points)

1. **Scan Go dependencies:**
   ```bash
   # Using govulncheck (official Go vulnerability checker)
   go install golang.org/x/vuln/cmd/govulncheck@latest
   govulncheck ./...
   ```

2. **Add a `vulnerability-scan` job to the CI pipeline** (see the TODO in `ci.yml`) that:
   - Runs after the `test` job
   - Installs `govulncheck` and runs it against the codebase
   - Fails if known vulnerabilities are found

   > **Hint:** Use `go install golang.org/x/vuln/cmd/govulncheck@latest` to install the tool.

3. **If vulnerabilities are found:**
   - Update the affected dependencies (`go get -u <module>`)
   - Document the CVEs and how you resolved them

**Deliverable:** govulncheck output. Updated `go.mod` if changes were needed.

---

### Task 3: Kubernetes Deployment with Minikube (8 Points)

1. **Start Minikube:**
   ```bash
   minikube start
   ```

2. **Build the image inside Minikube's Docker daemon:**
   ```bash
   eval $(minikube docker-env)
   docker build -t product-catalog:latest .
   ```

3. **Deploy the application:**
   ```bash
   kubectl apply -f k8s/namespace.yml
   kubectl apply -f k8s/postgres-deployment.yml
   kubectl apply -f k8s/api-deployment.yml
   ```

4. **Verify the deployment:**
   ```bash
   kubectl get all -n product-catalog
   kubectl logs deployment/product-catalog-api -n product-catalog
   ```

5. **Access the API:**
   ```bash
   minikube service product-catalog-api -n product-catalog --url
   # Use the returned URL to test the API
   curl <URL>/health
   curl <URL>/products
   ```

6. **Test CRUD operations** against the Kubernetes-deployed API.

**Deliverable:** Screenshots of:
- `kubectl get all -n product-catalog` output
- Successful API calls to the Kubernetes-hosted service
- Pod logs showing healthy operation

---

### Task 4: Production Readiness (6 Points)

1. **Scaling:** Scale the API deployment to 3 replicas and verify all pods are running:
   ```bash
   kubectl scale deployment product-catalog-api --replicas=3 -n product-catalog
   kubectl get pods -n product-catalog
   ```

2. **Health Checks:** The Kubernetes manifests include `readinessProbe` and `livenessProbe`. Explain:
   - What is the difference between a readiness and a liveness probe?
   - What happens if the readiness probe fails? What about the liveness probe?
   - Why are different `initialDelaySeconds` values used?

3. **Resource Limits:** The API deployment specifies CPU and memory limits. Explain:
   - What happens if a pod exceeds its memory limit?
   - What happens if it exceeds its CPU limit?
   - Why are requests and limits both specified?

**Deliverable:** Add a `K8S.md` file with your answers and screenshots.

---

## Kubernetes Manifest Overview

| File | Contents |
|------|----------|
| `k8s/namespace.yml` | Namespace `product-catalog` |
| `k8s/postgres-deployment.yml` | PostgreSQL Deployment, Service, Secret, PVC |
| `k8s/api-deployment.yml` | API Deployment (2 replicas), NodePort Service |

---

## Useful Commands

```bash
# Minikube
minikube start / stop / delete
minikube dashboard                    # Open Kubernetes dashboard
eval $(minikube docker-env)           # Use Minikube's Docker daemon

# kubectl
kubectl get pods -n product-catalog
kubectl describe pod <name> -n product-catalog
kubectl logs <pod-name> -n product-catalog
kubectl exec -it <pod-name> -n product-catalog -- /bin/sh
kubectl port-forward svc/product-catalog-api 8080:8080 -n product-catalog

# Trivy
trivy image <image>
trivy fs .                            # Scan filesystem/dependencies
```

---

Each exercise branch contains a detailed `README.md` with instructions.

| Task | Points |
|------|--------|
| Vulnerability Scanning -- Docker Image | 6 |
| Vulnerability Scanning -- Dependencies | 4 |
| Kubernetes Deployment with Minikube | 8 |
| Production Readiness | 6 |
| **Total** | **24** |

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)

