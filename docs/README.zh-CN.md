# DN42-GeoASN

> **Languages / 語言**: [English](../README.md) | [繁體中文](README.zh-TW.md)

一个为 DN42 网络生成 GeoLite2 兼容 ASN 数据库的工具。该项目从 DN42 注册表中提取路由和 ASN 信息，并将其转换为可与 GeoIP 库一起使用的 MaxMind MMDB 格式。

## 概述

DN42-GeoASN 包含两个主要组件：
- **Python 查找脚本**（`finder.py`）：解析 DN42 注册表以提取路由和源 ASN 信息
- **Go 生成器**（`generator.go`）：将提取的数据转换为 MaxMind MMDB 数据库文件

## 功能特性

- ✅ 支持 IPv4 和 IPv6 路由
- ✅ 从 DN42 注册表中提取 ASN 名称
- ✅ 生成与现有 GeoIP 库兼容的 MaxMind MMDB 格式
- ✅ 自动处理 DN42 注册表结构
- ✅ 跳过没有正确 ASN 名称的路由

## 先决条件

- Python 3.x
- Go 1.23.4 或更高版本
- DN42 注册表（本地克隆）

## 快速开始（预构建数据库）

如果您只想使用 DN42 ASN 数据库而不想自己构建，可以下载最新的预构建 MMDB 文件：

### 下载最新版本

最新的 `GeoLite2-ASN-DN42.mmdb` 文件会自动构建并发布在：
**https://github.com/rdp-studio/dn42-geoasn/releases**

```bash
# 下载最新版本
wget https://github.com/rdp-studio/dn42-geoasn/releases/latest/download/GeoLite2-ASN-DN42.mmdb

# 或使用 curl
curl -LO https://github.com/rdp-studio/dn42-geoasn/releases/latest/download/GeoLite2-ASN-DN42.mmdb
```

数据库会自动使用最新的 DN42 注册表数据更新并定期发布。

## 从源码构建

如果您想自己构建数据库或为项目做贡献：

### 安装

1. 克隆此仓库：
   ```bash
   git clone https://github.com/rdp-studio/dn42-geoasn.git
   cd dn42-geoasn
   ```

2. 克隆 DN42 注册表：
   ```bash
   git clone https://git.dn42.dev/dn42/registry.git
   ```

3. 安装 Go 依赖：
   ```bash
   go mod download
   ```

## 使用方法

### 步骤 1：提取路由数据

运行 Python 查找脚本从 DN42 注册表中提取路由和 ASN 信息：

```bash
python finder.py
```

这将：
- 解析 DN42 注册表中的所有路由和路由6对象
- 提取每个路由的源 ASN
- 从 aut-num 对象中查找 ASN 名称
- 生成包含提取数据的 `GeoLite2-ASN-DN42-Source.csv`

### 步骤 2：生成 MMDB 数据库

运行 Go 生成器将 CSV 数据转换为 MMDB 文件：

```bash
go run generator.go
```

这将创建 `GeoLite2-ASN-DN42.mmdb`，一个 MaxMind 兼容的数据库文件。

### 完整工作流程

```bash
# 从 DN42 注册表提取数据
python finder.py

# 生成 MMDB 文件
go run generator.go
```

## 输出文件

- `GeoLite2-ASN-DN42-Source.csv`：包含路由、ASN 和组织数据的中间 CSV 文件
- `GeoLite2-ASN-DN42.mmdb`：最终的 MaxMind MMDB 数据库文件

## CSV 格式

中间 CSV 文件包含三列：
1. **Network**：CIDR 表示法（例如 `172.20.0.0/14`）
2. **ASN**：自治系统号（不带"AS"前缀）
3. **Organization**：ASN 名称/组织

## MMDB 结构

生成的 MMDB 文件包含具有以下结构的记录：
- `autonomous_system_number`：ASN 作为 uint32
- `autonomous_system_organization`：组织名称作为字符串

## 与 GeoIP 库一起使用

生成的 MMDB 文件可以与标准的 MaxMind GeoIP 库一起使用：

### Python (maxminddb)
```python
import maxminddb

with maxminddb.open_database('GeoLite2-ASN-DN42.mmdb') as reader:
    result = reader.get('172.20.0.1')
    print(f"ASN: {result['autonomous_system_number']}")
    print(f"Org: {result['autonomous_system_organization']}")
```

### Go (maxminddb-golang)
```go
package main

import (
    "fmt"
    "log"
    "net"
    
    "github.com/oschwald/maxminddb-golang"
)

type ASNResult struct {
    ASN          uint32 `maxminddb:"autonomous_system_number"`
    Organization string `maxminddb:"autonomous_system_organization"`
}

func main() {
    db, err := maxminddb.Open("GeoLite2-ASN-DN42.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    ip := net.ParseIP("172.20.0.1")
    var result ASNResult
    err = db.Lookup(ip, &result)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("ASN: %d\n", result.ASN)
    fmt.Printf("Org: %s\n", result.Organization)
}
```

## 配置

可以在 `constant.py` 中修改以下常量：
- `REGISTRY_PATH`：DN42 注册表目录路径（默认："registry"）
- `SOURCE_OUTPUT`：输出 CSV 文件名（默认："GeoLite2-ASN-DN42-Source.csv"）

## 依赖要求

### Python 依赖
- 仅标准库（os、re、csv）

### Go 依赖
- `github.com/maxmind/mmdbwriter` v1.0.0
- `github.com/oschwald/maxminddb-golang` v1.12.0（间接）
- `go4.org/netipx` v0.0.0-20220812043211-3cc044ffd68d（间接）

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

该项目根据 MIT 许可证授权 - 有关详细信息，请参阅 [LICENSE](../LICENSE) 文件。

## 致谢

- [DN42](https://dn42.dev/) 社区维护注册表
- [MaxMind](https://www.maxmind.com/) 提供 MMDB 格式和库
- DN42 注册表维护者和贡献者
