
| 地址                                                                     | 运营方                 | 类型               | 说明                           |
| ---------------------------------------------------------------------- | ------------------- | ---------------- | ---------------------------- |
| [https://docker.1ms.run](https://docker.1ms.run)                       | 毫秒镜像（木雷坞）           | Cloudflare、境内CDN | 免费、支持搜索、配置简单、免费技术解答、集成多方主流生态 |
| [https://mirror.ccs.tencentyun.com](https://mirror.ccs.tencentyun.com) | 腾讯云                 | 境内CDN            | **仅腾讯云服务器内部可用**              |
| [https://docker.m.daocloud.io](https://docker.m.daocloud.io)           | DaoCloud 官方         | 阿里云服务器           | 白名单 & 限流                     |
| [https://docker.1panel.live](https://docker.1panel.live)               | 1Panel 官方           | Cloudflare       | 部分地区可能无法访问                   |
| [https://hub.rat.dev](https://hub.rat.dev)                             | 耗子面板官方              | 使用毫秒镜像           | 部分地区可能无法访问                   |
| [https://docker.1panel.dev](https://docker.1panel.dev)                 | 1Panel 核心用户无名驱动     | Cloudflare       | 部分地区可能无法访问                   |
| [https://docker.anye.in](https://docker.anye.in)                       | 1Panel 核心用户 Anye 驱动 | Cloudflare       | 部分地区可能无法访问                   |
| [https://docker.amingg.com](https://docker.amingg.com)                 | 爱铭网络官方              | Cloudflare       | 部分地区可能无法访问                   |
| [https://docker.367231.xyz](https://docker.367231.xyz)                 | 1Panel 核心用户 GXL 驱动  | Cloudflare       |                              |


root@iZbp17vix2j58ya7sc3b9lZ:~/docker-compose-grafana# cat /etc/docker/daemon.json
{
"registry-mirrors": ["https://94pzr3so.mirror.aliyuncs.com","https://docker.1ms.run"]
}
