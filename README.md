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
    image: cmstore
    command: ["cmstore", "-init", "-namespace", "default", "-name", "configmap1", "-dir", "/config"]
    volumeMounts:
    - name: config
      mountPath: /config
  containers:
  - name: main
    # ...
  - name: cmstore
    image: cmstore
    command: ["cmstore", "-sidecar", "-dir", "/config"]
    volumeMounts:
    - name: config
      mountPath: /config
  volumes:
  - name: config
    emptyDir: {}
```
