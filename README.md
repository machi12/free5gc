## 修改

### 1. AMF

- NFs/amf/internal/context/amf_ue.go：增加随机数N，便于对它的存储。
- NFs/amf/internal/gmm/handler.go：通过对注册请求消息的处理来对N进行赋值。
- NFs/amf/internal/gmm/message/build.go：增加对SNMAC的计算。

### 2. AUSF

- NFs/ausf/internal/sbi/producer/ue_authentication.go：从UDM的响应中获取HNMAC并赋值到发送给AMF的认证向量AV中。

### 3. UDM

- NFs/udm/internal/sbi/producer/generate_auth_data.go：将5G-AKA协议的相关计算修改为5G-ESAKA协议的计算过程。
