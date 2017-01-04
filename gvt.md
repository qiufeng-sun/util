# 工程依赖管理

## 使用工具: [gvt](https://github.com/FiloSottile/gvt)

```sh
# 安装gvt
$ go get -u github.com/FiloSottile/gvt

# 加入新的依赖
$ gvt fetch github.com/golang/protobuf/proto

# 重新获取已加入的依赖代码
$ gvt restore
```
* 注: 
> 1.gvt被安装在$GOPATH/bin下
> 2.加入依赖后, 可通过修改./vendor/manifest文件, 指定代码源路径, 目的路径(import路径), 及下载的版本
