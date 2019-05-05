工程依赖管理
================

* 使用的工具: gvp + gpm + gpm-git
---------------------------------

* [gvp](https://github.com/pote/gvp)用于设置GOPATH

```bash
$ git clone https://github.com/pote/gvp.git && cd gvp
$ git checkout v0.2.0 # You can ignore this part if you want to install HEAD.
$ ./configure
$ make install
```

* [gpm](https://github.com/pote/gpm)

```bash
$ git clone https://github.com/pote/gpm.git && cd gpm
$ git checkout v1.3.1 # You can ignore this part if you want to install HEAD.
$ ./configure
$ make install
```

* [gpm-git](https://github.com/technosophos/gpm-git)

```sh
$ git clone https://github.com/technosophos/gpm-git.git
$ git checkout v1.0.1
$ make install
```

* 获取工程, 并编译(以goproject/mining为例)
-----------------------------------------

```sh
$ mkdir work && cd ~/work # work为工作目录
$ git clone git@git.n.xiaomi.com:sunxiguang/goproject.git && cd goproject # 克隆goproject工程到work目录里
$ . gvp # 设置$GOPATH(go env可查看). 当前的为"/home/sxg/work/goproject/.godeps:/home/sxg/work/goproject"
$ cd src/mining # 进入子工程目录
$ gpm-git # 读取当前目录下Godeps-Git下载相关依赖
$ go install # 生成可执行文件(会生成在$GOBIN中)
$ cp $GOBIN/mining ~/work/goproject/bin/ # 从$GOBIN拷贝可执行文件到bin目录下
$ cp -r conf/ ~/work/goproject/bin/ # 拷贝配置文件到bin目录下
```
