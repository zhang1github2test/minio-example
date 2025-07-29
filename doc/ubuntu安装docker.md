

install-docker.sh

```sh
#!/bin/bash

set -e

echo "🔧 开始安装 Docker CE..."

# 卸载旧版本
sudo apt-get remove -y docker docker-engine docker.io containerd runc || true

# 更新 apt 并安装依赖
sudo apt-get update
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# 添加 Docker GPG 密钥
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
  sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# 添加 Docker 软件源
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装 Docker Engine
sudo apt-get update
sudo apt-get install -y \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin

# 启动 Docker 并设置开机自启
sudo systemctl enable docker
sudo systemctl start docker

# 添加当前用户到 docker 组（使其不用 sudo）
sudo usermod -aG docker $USER

echo "✅ Docker 安装完成！"
echo "📌 请执行 'newgrp docker' 或重新登录后使用 docker 命令"
echo "🧪 测试：docker run hello-world"

```

```shell
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["https://94pzr3so.mirror.aliyuncs.com","https://docker.1ms.run"]
}
EOF
sudo systemctl daemon-reload
sudo systemctl restart docker
```

