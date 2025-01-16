---
title: "Kubernetes"
---

Gravity can be run in Kubernetes, either as a single node or as a cluster.

## Single node

For a single node instance, the Gravity container can be deployed as-is.

Note that to use DHCP when running in Kubernetes, it is recommended to deploy a DHCP on the router and have it forward requests to Gravity's LoadBalancer IP.

See the example deployment below:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: gravity
---
apiVersion: v1
kind: Service
metadata:
  name: gravity
  namespace: gravity
  labels:
    app.kubernetes.io/name: gravity
spec:
  selector:
    app.kubernetes.io/name: gravity
  ports:
    - protocol: UDP
      port: 53
      targetPort: 53
      name: dns
    - protocol: TCP
      port: 53
      targetPort: 53
      name: dnstcp
    - port: 67
      targetPort: 67
      name: dhcp
      protocol: UDP
    - port: 69
      targetPort: 69
      name: tftp
      protocol: UDP
    - port: 8008
      targetPort: 8008
      name: http
      protocol: TCP
  type: LoadBalancer
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: gravity
  namespace: gravity
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: gravity
    spec:
      containers:
        - name: gravity
          image: ghcr.io/beryju/gravity:stable
          env:
            - name: BOOTSTRAP_ROLES
              value: dns;dhcp;api;discovery;backup;monitoring;tsdb;tftp
            - name: INSTANCE_IP
              value: <insert loadbalancer IP of gravity>
            - name: INSTANCE_LISTEN
              value: 0.0.0.0
          livenessProbe:
            httpGet:
              path: /healthz/live
              port: http-metrics
          readinessProbe:
            httpGet:
              path: /healthz/ready
              port: http-metrics
          ports:
            - containerPort: 53
              name: dns-tcp
            - containerPort: 53
              protocol: UDP
              name: dns-udp
            - containerPort: 67
              protocol: UDP
              name: dhcp
            - containerPort: 68
              protocol: UDP
              name: dhcp-alt
            - containerPort: 69
              protocol: UDP
              name: tftp
            - containerPort: 8008
              name: http
            - containerPort: 8009
              name: http-metrics
          volumeMounts:
            - name: data
              mountPath: /data
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 1Gi
```

## Multi-node cluster

When running a multi-node Gravity cluster in Kubernetes, it is recommended to not use the built-in etcd server in Gravity. This is because the built-in cluster management is optimized for static nodes and simplicity and as such doesn't offer enough dynamic options that Kubernetes requires.

It is recommended to use either an etcd Operator like [this](https://etcd.aenix.io/) or [this](https://operatorhub.io/operator/etcd), deploy etcd manually, or use Kubernetes' etcd with Gravity.
