# OpenAI API Proxy with Key Pool

`openai-api-proxy-key-pool` 是一个反向代理服务，用于代理 OpenAI API 的请求。同时，它还提供了轮询池 key 调用功能，使您能够在访问 OpenAI API 时轮询使用不同的 API 密钥。

## 功能

- 反向代理 OpenAI API 请求
- 轮询池 key 调用，实现对不同 API 密钥的轮询使用
- 配置文件支持，方便管理 API 密钥

## 安装

### 方法 1：使用 Go 编译

1. 克隆项目到本地：

```bash
git clone https://github.com/surenkid/openai-api-proxy-key-pool.git
```

2. 进入项目目录：

```bash
cd openai-api-proxy-key-pool
```

3. 使用 Go 编译项目（确保已安装 [Go](https://golang.org/doc/install)）：

```bash
go build
```

### 方法 2：使用 Docker 部署

确保已安装 [Docker](https://docs.docker.com/engine/install/)。

运行以下命令，将代理服务部署为 Docker 容器：

```bash
docker run -p 8124:8124 -v /home/volume/openai-api-proxy-key-pool/config.json:/config/config.json surenkid/openai-api-proxy-key-pool:latest
```

### 方法 3：使用 Docker Compose 部署

确保已安装 [Docker Compose](https://docs.docker.com/compose/install/)。

1. 在项目根目录创建一个名为 `docker-compose.yml` 的文件，并填充以下内容：

```yaml
version: '3'

services:
  openai-api-proxy-key-pool:
    image: surenkid/openai-api-proxy-key-pool:latest
    ports:
      - "8124:8124"
    volumes:
      - ./config/config.json:/config/config.json
```

2. 运行以下命令，使用 Docker Compose 启动代理服务：

```bash
docker-compose up
```

## 使用

1. 更新 `config/config.json` 文件，添加您的 API 密钥。参考配置文件示例：

```json
{
  "keys": {
    "ai-001abc": [
      "sk-abcdef1234567890abcdefghijklmnopqrstuvwxyz000001",
      "sk-abcdef1234567890abcdefghijklmnopqrstuvwxyz000002",
      "sk-abcdef1234567890abcdefghijklmnopqrstuvwxyz000003"
    ],
    "ai-002def": [
      "sk-abcdef1234567890abcdefghijklmnopqrstuvwxyz000004"
    ]
  ,
  "helicone": "sk-123456"
}
```

2. 如果使用 Go 编译方法，运行编译好的二进制文件：

```bash
./openai-api-proxy-key-pool
```

如果使用 Docker 部署或 Docker Compose 部署方法，容器已在上一步启动。

3. 代理服务将在端口 `8124` 上启动，您可以将您的请求发送到 `http://localhost:8124`。
