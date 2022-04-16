
controller-gen:
```bash
controller-gen "crd:trivialVersions=true,preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```

```bash
controller-gen object:headerFile="hack\\boilerplate.go.txt" paths="./..."
```

## 测试工具

安装
```bash
go get github.com/onsi/ginkgo/ginkgo
```

第一次创建：
```bash
ginkgo bootstrap
```

运行：
