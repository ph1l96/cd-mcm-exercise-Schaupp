# Kubernetes Production Readiness

## Scaling

The API deployment can be scaled to three replicas with:

```bash
minikube kubectl -- scale deployment product-catalog-api --replicas=3 -n product-catalog
```

Verification command:

```bash
kubectl get pods -n product-catalog
```

```
NAME                                       READY   STATUS    RESTARTS   AGE
pod/postgres-58c46544bc-lqnqc              1/1     Running   0          26m
pod/product-catalog-api-7d85d9445c-hdw9k   1/1     Running   0          17m
pod/product-catalog-api-7d85d9445c-mk7xj   1/1     Running   0          17m
pod/product-catalog-api-7d85d9445c-wbkzq   1/1     Running   0          17s

NAME                          TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
service/postgres              ClusterIP   10.100.28.11   <none>        5432/TCP         26m
service/product-catalog-api   NodePort    10.99.61.13    <none>        8080:30080/TCP   26m

NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/postgres              1/1     1            1           26m
deployment.apps/product-catalog-api   3/3     3            3           26m

NAME                                             DESIRED   CURRENT   READY   AGE
replicaset.apps/postgres-58c46544bc              1         1         1       26m
replicaset.apps/product-catalog-api-77c8d8bb78   0         0         0       26m
replicaset.apps/product-catalog-api-7d85d9445c   3         3         3       17m
```

## Health Checks

### Readiness vs. liveness probe: whats the difference?

A readiness probe checks whether a container is ready to receive traffic. In this project, the API readiness probe calls `/health` on port `8080`. For the PostgreSQL-backed API, this endpoint also checks whether the database can be reached. If the readiness probe succeeds, Kubernetes includes the pod as an endpoint behind the `product-catalog-api` service.

A liveness probe checks whether the container is still healthy enough to keep running. In this project, the liveness probe also calls `/health` on port `8080`, but it is used for a different decision: whether Kubernetes should restart the container.

### What happens when each probe fails?

If the readiness probe fails, Kubernetes keeps the pod running but removes it from service traffic. The pod is not restarted only because readiness failed. This is useful when the application is still starting, temporarily overloaded, or cannot reach a dependency such as PostgreSQL.

If the liveness probe fails repeatedly, Kubernetes assumes the container is broken and restarts it. This helps recover from deadlocks, stuck application processes, or other states where the process is still running but no longer healthy.

### Why different initialDelaySeconds values?

The readiness probe uses a shorter delay:

```yaml
readinessProbe:
  initialDelaySeconds: 5
```

This allows Kubernetes to start sending traffic soon after the API is able to answer requests.

The liveness probe uses a longer delay:

```yaml
livenessProbe:
  initialDelaySeconds: 15
```

This gives the application more time to start before Kubernetes begins making restart decisions. Without a longer liveness delay, Kubernetes might restart the container while it is still performing normal startup work, especially because this API connects to PostgreSQL before it can serve successfully.

## Resource Limits

### What happens if memory or CPU limits are exceeded?

If a container exceeds its memory limit, Kubernetes can terminate it with an out-of-memory error. The pod then restarts according to the deployment's restart policy. In this project, the API memory limit is:

```yaml
limits:
  memory: "128Mi"
```

If a container tries to use more CPU than its CPU limit, it is throttled instead of immediately killed. The application keeps running, but requests may become slower because Kubernetes restricts how much CPU time the container can use. In this project, the API CPU limit is:

```yaml
limits:
  cpu: "250m"
```

### Why specify both requests and limits?

Requests tell Kubernetes how many resources the container needs for scheduling. The scheduler uses these values to decide which node has enough capacity for the pod. In this project, the API requests:

```yaml
requests:
  memory: "64Mi"
  cpu: "100m"
```

Limits define the maximum resources the container is allowed to use. They protect the cluster from one container consuming too much memory or CPU and affecting other workloads.
Using both requests and limits gives Kubernetes enough information to place pods reliably while also enforcing an upper boundary for runtime resource usage.
