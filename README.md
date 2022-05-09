## 100行实现UDP加密转发

### Example

路径：node1->node2

node1:
```
forward.Start("0.0.0.0:1194", "node2_ip:11111", "aes_key", "aes_iv")
```
node2:
```
forward.Start("0.0.0.0:11111", "127.0.0.1:1194", "aes_key", "aes_iv")
```
加密方式为AES-256-CTR，具体例子请看`cmd/node1/main.go`和`cmd/node2/main.go`
