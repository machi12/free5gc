## 本项目的修改

### 1. AMF

- NFs/amf/internal/context/amf_ue.go：增加随机数N，便于对它的存储。
- NFs/amf/internal/gmm/handler.go：通过对注册请求消息的处理来对N进行赋值。
- NFs/amf/internal/gmm/message/build.go：增加对SNMAC的计算。

### 2. AUSF

- NFs/ausf/internal/sbi/producer/ue_authentication.go：从UDM的响应中获取HNMAC并赋值到发送给AMF的认证向量AV中。

### 3. UDM

- NFs/udm/internal/sbi/producer/generate_auth_data.go：将5G-AKA协议的相关计算修改为5G-ESAKA协议的计算过程。

## 对依赖项目的修改

- [nas]{https://github.com/machi12/nas}
- [openapi]{https://github.com/machi12/openapi}
- [util]{https://github.com/machi12/util}

## 构建docker镜像

1. 克隆free5gc-compose项目
```
git clone https://github.com/free5gc/free5gc-compose.git
cd free5gc-compose
```

2. 在base目录下克隆修改后的free5gc项目
```
cd base
git clone https://github.com/machi12/free5gc.git
```

3. 添加go代理，在base目录下的Dockerfile文件中添加如下代码
```
#下面三条配置用于保证在不同go语言版本中配置代理一定生效，建议都写上  
RUN export GOPROXY=https://goproxy.io  
RUN export GO111MODULE=on  
RUN go env -w GOPROXY=https://goproxy.io 
```

4. 构建docker镜像（需要切换到free5gc-compose目录下）
```
cd ..
make all
docker compose -f docker-compose-build.yaml build
```

## 报错解决

- go.sum中对aper包的哈希检查失败，由于重新生成go.sum文件后该包的哈希值仍无法与下载的对应，只能暂时将对应的哈希值改为报错信息中下载包的哈希值。
```
go: downloading github.com/free5gc/aper v1.0.5
verifying github.com/free5gc/aper@v1.0.5: checksum mismatch
downloaded: h1:oErgVDTYmOYmPDCERQosVg7dTBIqhloX1vnYko5xgDY=
go.sum:     h1:sUYFFmOXDLjyL4rU6zFnq81M4YluqP90Pso5e/J4UhA=
```

- 构建docker镜像时出现磁盘空间不足问题，通过``docker system prune``命令来删除所有未使用的资源。或者给虚拟机分配更多的磁盘空间。
```
ERROR: failed to update builder last activity time: write /home/machi/.docker/buildx/activity/.tmp-default1617546388: no space left on device
```