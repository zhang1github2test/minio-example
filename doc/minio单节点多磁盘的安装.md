​	MinIO 是一个**高性能、开源的对象存储系统**，主要用于存储非结构化数据（如图片、视频、文档、备份等），与 Amazon S3 完全兼容。它被广泛用于云原生应用、大数据分析、AI 模型存储、容器平台（如 Kubernetes）等场景。

MinIO 支持多种部署模式，其中：

> **单节点单磁盘（Single-Node Single-Drive）** 模式适用于开发测试、小规模应用或资源受限的场景。它的部署简单，不依赖集群、分布式架构或复杂的底层存储系统。

相比之下：

> 本文介绍的是 **单节点多磁盘（Single-Node Multi-Drive, SNMD）** 模式，适用于需要**基本容错能力和磁盘级可靠性**的生产场景。SNMD 部署利用 MinIO 的 **纠删码（Erasure Coding）** 技术，在单节点的前提下，实现了对单盘故障的自动恢复。

SNMD 模式为你提供了较好的数据安全性，但其性能与扩展能力受限于单节点资源。因此：

> 💡 **MinIO 官方建议：在生产环境中使用 “多节点多磁盘（Distributed）” 部署模式，以获得企业级的性能、可用性和可扩展性。**

------

🛠️ 下面是 MinIO 单节点多磁盘部署的详细安装步骤。该模式适合希望简化部署的用户，同时又希望具备基础的容错和高可用能力。

官网地址：https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-multi-drive.html

---

## 一、MinIO 单节点单磁盘安装步骤

### 1. 环境准备

说明：这里只是演示安装步骤，如果需要可靠环境，磁盘数量至少为6块。

* 操作系统：建议使用 Linux Ubuntu 24.04

* 硬件：2C  16G  40G*2

  ### 1.1 磁盘格式化及挂载

  format_and_mount.sh

  ```
  #!/bin/bash
  
  set -e
  
  # 检查参数
  if [ $# -ne 2 ]; then
    echo "用法: $0 <磁盘设备名> <挂载点目录>"
    echo "示例: $0 /dev/vdb /mnt/data1"
    exit 1
  fi
  
  DISK_DEVICE="/dev/$1"
  MOUNT_POINT=$2
  
  # 检查磁盘是否存在
  if [ ! -b "$DISK_DEVICE" ]; then
    echo "错误：设备 $DISK_DEVICE 不存在。"
    exit 2
  fi
  
  # 创建挂载目录
  echo "创建挂载目录 $MOUNT_POINT..."
  mkdir -p "$MOUNT_POINT"
  
  # 格式化磁盘为XFS
  echo "格式化 $DISK_DEVICE 为 XFS 文件系统..."
  mkfs.xfs -f "$DISK_DEVICE"
  
  # 获取UUID
  UUID=$(blkid -s UUID -o value "$DISK_DEVICE")
  if [ -z "$UUID" ]; then
    echo "获取 UUID 失败，退出。"
    exit 3
  fi
  
  # 挂载磁盘
  echo "挂载 $DISK_DEVICE 到 $MOUNT_POINT..."
  mount "$DISK_DEVICE" "$MOUNT_POINT"
  
  # 备份 fstab 并写入自动挂载配置
  echo "备份 /etc/fstab 为 /etc/fstab.bak..."
  cp /etc/fstab /etc/fstab.bak
  
  echo "写入开机自动挂载配置..."
  grep -q "$UUID" /etc/fstab || echo "UUID=$UUID $MOUNT_POINT xfs defaults 0 0" >> /etc/fstab
  
  echo "挂载完成，验证挂载信息："
  df -h | grep "$MOUNT_POINT"
  
  echo "✅ 操作完成。"
  ```

  使用lsblk来获取没有挂载的磁盘

  ```sh
  root@iZbp14hifhvqn83v9pkcscZ:~# lsblk
  NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINTS
  nvme1n1     259:0    0   40G  0 disk
  nvme0n1     259:1    0   40G  0 disk
  ├─nvme0n1p1 259:3    0    1M  0 part
  ├─nvme0n1p2 259:4    0  200M  0 part /boot/efi
  └─nvme0n1p3 259:5    0 39.8G  0 part /
  nvme2n1     259:2    0   40G  0 disk
  
  ```

  从上面的输出我们可以看出，有两磁盘(,nvme1n1\nvme2n1)尚未存在挂载点。

  ```sh
  ./format_and_mount.sh nvme1n1  /mnt/data1
  ./format_and_mount.sh nvme2n1  /mnt/data2
  ```

  

  

---

### 2. 下载并安装minio

* 安装包方式安装

```bash
wget https://dl.min.io/server/minio/release/linux-amd64/archive/minio_20250613113347.0.0_amd64.deb -O minio.deb
sudo dpkg -i minio.deb
```

使用上面的步骤安装的minio不需要单独在编写Systemd启动服务文件，如果通过二进制方式需要按照第5步所述操作执行。

二进制方式安装

```
wget https://dl.min.io/server/minio/release/linux-arm64/minio
chmod +x minio
mv minio /usr/local/bin/
```

### 3. MinIO 数据目录

这里我们直接使用两块磁盘对应的挂载点。/mnt/data1、/mnt/data2  



### 4. 创建 MinIO 运行用户

```bash
groupadd -r minio-user
useradd -M -r -g minio-user minio-user
chown minio-user:minio-user /mnt/data1  /mnt/data2
```



### 5. 编写 Systemd 启动服务文件

```bash
sudo nano /etc/systemd/system/minio.service
```

内容如下：

```ini
[Unit]
Description=MinIO
Documentation=https://docs.min.io
Wants=network-online.target
After=network-online.target
AssertFileIsExecutable=/usr/local/bin/minio

[Service]
Type=notify

WorkingDirectory=/usr/local

User=minio-user
Group=minio-user
ProtectProc=invisible

EnvironmentFile=-/etc/default/minio
ExecStart=/usr/local/bin/minio server $MINIO_OPTS $MINIO_VOLUMES

# Let systemd restart this service always
Restart=always

# Specifies the maximum file descriptor number that can be opened by this process
LimitNOFILE=1048576

# Turn-off memory accounting by systemd, which is buggy.
MemoryAccounting=no

# Specifies the maximum number of threads this process can create
TasksMax=infinity

# Disable timeout logic and wait until process is stopped
TimeoutSec=infinity

# Disable killing of MinIO by the kernel's OOM killer
OOMScoreAdjust=-1000

SendSIGKILL=no

[Install]
WantedBy=multi-user.target

# Built for ${project.name}-${project.version} (${project.name})
```

---

### 6. 编写 环境变量文件

**/etc/default/minio**

```properties
# MINIO_ROOT_USER and MINIO_ROOT_PASSWORD sets the root account for the MinIO server.
# This user has unrestricted permissions to perform S3 and administrative API operations on any resource in the deployment.
# Omit to use the default values 'minioadmin:minioadmin'.
# MinIO recommends setting non-default values as a best practice, regardless of environment
MINIO_ROOT_USER=myminioadmin
MINIO_ROOT_PASSWORD=minio-secret-key-change-me

# MINIO_VOLUMES sets the storage volume or path to use for the MinIO server.

MINIO_VOLUMES="/mnt/data{1...2}"

# MINIO_OPTS sets any additional commandline options to pass to the MinIO server.
# For example, `--console-address :9001` sets the MinIO Console listen port
MINIO_OPTS="--console-address :9001"
```



### 7. 启动并启用 MinIO 服务

```bash
sudo systemctl daemon-reload
# 设置开机自启动
sudo systemctl enable minio
# 启动minio
sudo systemctl start minio
# 查看minio启动状态
sudo systemctl status minio
```

---

### 8. 访问 MinIO

* 控制台地址（Web UI）：`http://<your-ip>:9001`
* API 端点：`http://<your-ip>:9000`

使用上述配置的 myminioadmin/minio-secret-key-change-me登录即可

如：http://115.29.205.126:9001/login

![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/04c83508e09d458fbd48766fabb2d7bf.png)
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/88d4569dedba4f84884f69b3ddaf247d.png)




---

##  二、使用MC工具测试集群的可用性

```bash
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc
sudo mv mc /usr/local/bin/
```

添加 MinIO 实例：

```bash
mc alias set local http://localhost:9000 myminioadmin minio-secret-key-change-me
```

测试操作：

```bash
# 创建桶
root@iZbp1bnsdgk7l6gjkh64wzZ:~# mc mb local/mytestbucket
Bucket created successfully `local/mytestbucket`.   
#上传文件
root@iZbp1bnsdgk7l6gjkh64wzZ:~# mc cp ./README.md local/mytestbucket
/root/README.md:                    43.86 KiB / 43.86 KiB  1.25 MiB/s 0s

# 查看文件
root@iZbp1bnsdgk7l6gjkh64wzZ:~# mc ls local/mytestbucket
[2025-07-08 11:27:24 CST]  44KiB STANDARD README.md
```

注意上面的local要跟之前的mc alias set local保持一致

##  二、使用场景说明（单节点多磁盘）



###  适用场景：

| 使用场景                   | 说明                                                         |
| -------------------------- | ------------------------------------------------------------ |
| 💾 高容量归档存储           | 比如视频监控存储、离线数据备份等，适合使用多块本地磁盘构建大容量对象存储。 |
| 🧪 测试或准生产环境         | 在不中断业务前提下验证多磁盘纠删码容错能力的试运行环境。     |
| 🏷️ 部署成本受限             | 仅有单台物理机可用，仍希望拥有数据冗余与可靠性的中小型应用系统。 |
| 📁 文件分发/对象缓存节点    | 构建边缘节点用于静态文件分发，如文档、图像、压缩包缓存。     |
| 🚧 数据安全优先级高于扩展性 | 对可用性和存储安全有一定要求，但无需分布式高可用架构。       |

---

### 🚫 不适用场景：

| 不适用场景                     | 原因说明                                                     |
| ------------------------------ | ------------------------------------------------------------ |
| 🏢 企业级生产环境               | SNMD 模式缺乏节点级高可用性，一旦主机宕机，所有服务不可用。MinIO 官方明确建议使用 **分布式部署（Multi-Node Multi-Drive）** 实现高可用。 |
| 🚀 高并发、大吞吐量应用         | 受限于单节点 CPU、内存、网络带宽，无法支撑大规模并发访问（如视频直播、在线媒体平台）。 |
| 🔄 多用户并发上传/下载          | 并发访问时瓶颈集中在单节点 I/O 和网卡，存在明显性能瓶颈，易出现延迟或失败。 |
| ☁️ 云原生大规模存储集群         | 与 Kubernetes、容器调度集群搭配时，SNMD 模式缺乏横向扩展能力，无法支撑动态服务扩展。 |
| ☠️ 无法接受单点故障             | 如果对**业务连续性要求极高**（如金融、医疗系统），SNMD 不具备节点容灾能力，不适合关键生产系统。 |
| 🧩 需要跨机房、跨地区冗余的场景 | 单节点无法支持多地分布式复制，数据容灾和地理冗余能力不足。   |

---

## 总结

通过本教程，你已完成：

1. 🔧 准备多块磁盘并挂载至指定目录；
2. 📦 安装并配置 MinIO 单节点服务；
3. ⚙️ 通过 systemd 启动并注册为系统服务；
4. 🔍 验证多磁盘纠删码是否生效，上传文件验证对象存储可用性；
5. 🛠️ 初步了解了 SNMD 模式的存储机制、适用与不适用场景。