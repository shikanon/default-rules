# 规范

## 命名空间节点亲和性

namespace 的 label 添加 `shikanon.com/node-affinity` 字段。

|----:|----:|
| label字段 | 使用说明 |
|----:|----:|
| `shikanon.com/node-affinity` | 一个键值对的json |

比如：
`shikanon.com/node-affinity: '{"dev-groups": ["rcmd", "media"]}'`

将被转换为：
```yaml
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: dev-groups
            operator: In
            values:
            - rcmd
            - media
```