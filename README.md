# cmstore

```yaml
apiVersion: v1
kind: Pod
metadata:
  namespace: default
  name: cmstore-example
spec:
  initContainers:
  - name: init-cmstore
    image: ghcr.io/ryotarai/cmstore:master
    command: ["cmstore", "init", "--namespace", "default", "--name", "configmap1", "--dir", "/config", "--create-if-not-found"]
    volumeMounts:
    - name: config
      mountPath: /config
  containers:
  - name: main
    image: ubuntu:20.04
    command: ["sleep", "infinity"]
    volumeMounts:
    - name: config
      mountPath: /config
  - name: cmstore
    image: ghcr.io/ryotarai/cmstore:master
    command: ["cmstore", "watch", "--namespace", "default", "--name", "configmap1", "--dir", "/config"]
    volumeMounts:
    - name: config
      mountPath: /config
  volumes:
  - name: config
    emptyDir: {}
```
