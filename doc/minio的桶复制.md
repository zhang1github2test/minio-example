## 🌟前言
​          在对象存储系统中，数据可靠性 和 业务连续性 是企业级部署绕不开的核心话题。MinIO 作为一款高性能的分布式对象存储服务，提供了功能强大的 桶复制（Bucket Replication） 机制，可实现数据跨站点、跨机房乃至跨云的复制能力，广泛应用于灾难恢复（Disaster Recovery）、数据同步、多活部署等场景。
## 一、服务端复制
 两个minio实例的数据
| 源实例      (src-minio)                                 | 目标实例 (dst-minio)                                   |
|------------------------------------------------------|------------------------------------------------------|
| http://minioadmin:minioadmin@172.18.0.1:9000/mybucket | http://minioadmin:minioadmin@172.18.0.1:19000/mybucket |


mc工具配置
```shell
mc alias set src-minio http://172.18.0.1:9000  minioadmin minioadmin
mc alias set dst-minio http://172.18.0.1:19000 minioadmin minioadmin
mc mb src-minio/mybucket
mc mb dst-minio/mybucket
```
### 1.0、桶复制的前置条件
    1、桶必须开启的版本控制
```shell
    mc version enable src-minio/mybucket
    mc version enable dst-minio/mybucket
```


### 1.1 单向复制
  单向复制是指将数据从一个 MinIO 实例复制到另一个 MinIO 实例，复制方向是单向的。 单向复制的配置相对简单，只需要在源实例上配置复制规则，指定目标实例即可。
  * 创建复制规则
    ```shell
    mc replicate add src-minio/mybucket \
    --remote-bucket 'http://minioadmin:minioadmin@localhost:19000/mybucket' \
    --replicate "delete,delete-marker,existing-objects"
    ```
    参数解释：
      * `src-minio/mybucket`：src-minio源实例的别名，mybucket:要进行复制的桶。
      * `--remote-bucket`：目标实例的存储桶名称，格式为 `https://USER:PASSWORD@HOSTNAME:PORT/bucketname`。
      * `--replicate`：复制规则， 代表启动了复制存在的对象和删除的对象。
  * 验证复制
    ```shell
       mc cp ~/foo.txt src-minio/BUCKET
       # 稍等几秒
        mc ls dst-minio/BUCKET
    ```
    如果在目标实例上能够看到复制过来的对象，说明复制成功。
  * 查看已经配置的复制规则
    ```shell
        mc replicate list src-minio/mybucket
    ```
    输出结果：
    ```txt
    Rules:
    Remote Bucket: 172.18.0.1:19000/mybucket
    Rule ID: d249gahsnv9im354pqc0
    Priority: 0
    ARN: arn:minio:replication::da9655e8-4ff9-497f-89e2-adbe91ab287e:mybucket
    ```
   * **Remote Bucket**

     * 复制的目标桶。
        * 结构为：`<目标地址>/<目标桶名>`。
        * 示例：`172.18.0.1:19000/mybucket` 表示目标 MinIO 实例的 `mybucket` 桶。

        * **Rule ID**

      * 当前复制规则的唯一标识符。
      * 系统自动生成，用于后续规则管理，如修改、删除。

        * **Priority**

      * 复制规则的优先级，值越小优先级越高。
      * 当多个规则同时匹配某个对象时，优先级高（值小）的规则生效。
      * 默认值通常为 `0`（最高优先级）。

        * **ARN**（Amazon Resource Name）

      * 用于唯一标识该复制规则所属的资源。
      * 格式：`arn:minio:replication::<MinIO实例ID>:<源桶名>`。
      * 示例：`arn:minio:replication::da9655e8-4ff9-497f-89e2-adbe91ab287e:mybucket`

          * `da9655e8-4ff9-497f-89e2-adbe91ab287e`：MinIO 实例唯一 ID。
          * `mybucket`：当前桶名。


  * 禁用桶复制
    ```shell
    mc replicate disable src-minio/mybucket
    ```
    * 禁用后查看具体的规则的复制状态
      ```shell
      root@iZbp17vix2j58ya7sc3b9lZ:~# mc replicate export src-minio/mybucket
      {
      "Rules": [
          {
              "ID": "d249gahsnv9im354pqc0",
              "Status": "Disabled",
              "Priority": 0,
              "DeleteMarkerReplication": {
                  "Status": "Enabled"
                  },
              "DeleteReplication": {
                  "Status": "Enabled"
                  },
              "Destination": {
                  "Bucket": "arn:minio:replication::da9655e8-4ff9-497f-89e2-adbe91ab287e:mybucket"
                  },
              "Filter": {
                  "And": {},
                  "Tag": {}
                  },
              "SourceSelectionCriteria": {
                  "ReplicaModifications": {
                      "Status": "Enabled"
                      }
              },
              "ExistingObjectReplication": {
                  "Status": "Enabled"
                  }
          }
      ],
      "Role": ""
      }
      ```
      可以看出来,当前复制规则已经被禁用了。 "Status": "Disabled"
    * 移除复制规则
      ```shell
        mc replicate rm  --id "d249gahsnv9im354pqc0" src-minio/mybucket
      ```
    > 如果当前的桶只有一条复制规则，minio是不允许删除的。
    root@iZbp17vix2j58ya7sc3b9lZ:~# mc replicate rm  --id "d249gahsnv9im354pqc0" src-minio/mybucket
    mc: <ERROR> Could not remove replication rule: replication configuration should have at least one rule.
    
### 1.2 双向复制
  双向复制是指在两个 MinIO 实例之间建立双向复制关系，数据可以在两个实例之间双向同步。
  * 配置源实例的复制规则
    ```shell
    mc replicate add src-minio/mybucket \
    --remote-bucket 'http://minioadmin:minioadmin@172.18.0.1:19000/mybucket' \
    --replicate "delete,delete-marker,existing-objects"
    ```
  * 配置目标实例的复制规则
    ```shell
    mc replicate add dst-minio/mybucket \
    --remote-bucket 'http://minioadmin:minioadmin@172.18.0.1:9000/mybucket' \
    --replicate "delete,delete-marker,existing-objects"

    ```
  * 验证复制
    ```shell
    mc cp ~/src.txt src-minio/BUCKET
    mc cp ~/dst.txt dst-minio/BUCKET
    # 稍等几秒
    mc ls dst-minio/BUCKET
    mc ls src-minio/BUCKET
    ```
    如果在目标实例上和源实例上都能够看到复制过来的对象，说明复制成功。
## 二、客户端的桶复制
minio的客户端复制的源头可以是文件系统，minio集群，或其他S3兼容的对象存储服务。
### 2.1 基于文件系统的复制
```shell
   mkdir ~/file-system-test  && cd ~/file-system-test
   mc mirror --watch ~/file-system-test src-minio/mybucket
```

### 2.2 基于minio集群的复制
```shell
   mc mirror  src-minio/mybucket dst-minio/mybucket
```

## 三、注意事项
   1、如果是在同一台电脑上通过安装了两个minio集群。
   > 如果两个minio集群是使用docker安装的，这里在演示的时候需要使用ifconfig查看docker0网卡的ip，
   并将该ip地址替换到上面的配置中,不能直接使用localhost,不然会在执行 mc replicate出现连接失败的问题。

