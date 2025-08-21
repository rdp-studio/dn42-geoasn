# DN42-GeoASN

> **Languages / 语言**: [English](../README.md) | [简体中文](README.zh-CN.md)

一個為 DN42 網路產生 GeoLite2 相容 ASN 資料庫的工具。此專案從 DN42 註冊表中提取路由和 ASN 資訊，並將其轉換為可與 GeoIP 程式庫一起使用的 MaxMind MMDB 格式。

## 概述

DN42-GeoASN 包含兩個主要元件：
- **Python 查找腳本**（`finder.py`）：解析 DN42 註冊表以提取路由和來源 ASN 資訊
- **Go 產生器**（`generator.go`）：將提取的資料轉換為 MaxMind MMDB 資料庫檔案

## 功能特性

- ✅ 支援 IPv4 和 IPv6 路由
- ✅ 從 DN42 註冊表中提取 ASN 名稱
- ✅ 產生與現有 GeoIP 程式庫相容的 MaxMind MMDB 格式
- ✅ 自動處理 DN42 註冊表結構
- ✅ 跳過沒有正確 ASN 名稱的路由

## 先決條件

- Python 3.x
- Go 1.23.4 或更高版本
- DN42 註冊表（本機複製）

## 快速開始（預建置資料庫）

如果您只想使用 DN42 ASN 資料庫而不想自己建置，可以下載最新的預建置 MMDB 檔案：

### 下載最新版本

最新的 `GeoLite2-ASN-DN42.mmdb` 檔案會自動建置並發佈在：
**https://github.com/rdp-studio/dn42-geoasn/releases**

```bash
# 下載最新版本
wget https://github.com/rdp-studio/dn42-geoasn/releases/latest/download/GeoLite2-ASN-DN42.mmdb

# 或使用 curl
curl -LO https://github.com/rdp-studio/dn42-geoasn/releases/latest/download/GeoLite2-ASN-DN42.mmdb
```

資料庫會自動使用最新的 DN42 註冊表資料更新並定期發佈。

## 從原始碼建置

如果您想自己建置資料庫或為專案做出貢獻：

### 安裝

1. 複製此儲存庫：
   ```bash
   git clone https://github.com/rdp-studio/dn42-geoasn.git
   cd dn42-geoasn
   ```

2. 複製 DN42 註冊表：
   ```bash
   git clone https://git.dn42.dev/dn42/registry.git
   ```

3. 安裝 Go 依賴項：
   ```bash
   go mod download
   ```

## 使用方法

### 步驟 1：提取路由資料

執行 Python 查找腳本從 DN42 註冊表中提取路由和 ASN 資訊：

```bash
python finder.py
```

這將：
- 解析 DN42 註冊表中的所有路由和路由6物件
- 提取每個路由的來源 ASN
- 從 aut-num 物件中查找 ASN 名稱
- 產生包含提取資料的 `GeoLite2-ASN-DN42-Source.csv`

### 步驟 2：產生 MMDB 資料庫

執行 Go 產生器將 CSV 資料轉換為 MMDB 檔案：

```bash
go run generator.go
```

這將建立 `GeoLite2-ASN-DN42.mmdb`，一個 MaxMind 相容的資料庫檔案。

### 完整工作流程

```bash
# 從 DN42 註冊表提取資料
python finder.py

# 產生 MMDB 檔案
go run generator.go
```

## 輸出檔案

- `GeoLite2-ASN-DN42-Source.csv`：包含路由、ASN 和組織資料的中間 CSV 檔案
- `GeoLite2-ASN-DN42.mmdb`：最終的 MaxMind MMDB 資料庫檔案

## CSV 格式

中間 CSV 檔案包含三欄：
1. **Network**：CIDR 表示法（例如 `172.20.0.0/14`）
2. **ASN**：自治系統號碼（不帶「AS」前綴）
3. **Organization**：ASN 名稱/組織

## MMDB 結構

產生的 MMDB 檔案包含具有以下結構的記錄：
- `autonomous_system_number`：ASN 作為 uint32
- `autonomous_system_organization`：組織名稱作為字串

## 與 GeoIP 程式庫一起使用

產生的 MMDB 檔案可以與標準的 MaxMind GeoIP 程式庫一起使用：

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

## 設定

可以在 `constant.py` 中修改以下常數：
- `REGISTRY_PATH`：DN42 註冊表目錄路徑（預設：「registry」）
- `SOURCE_OUTPUT`：輸出 CSV 檔案名稱（預設：「GeoLite2-ASN-DN42-Source.csv」）

## 依賴需求

### Python 依賴項
- 僅標準程式庫（os、re、csv）

### Go 依賴項
- `github.com/maxmind/mmdbwriter` v1.0.0
- `github.com/oschwald/maxminddb-golang` v1.12.0（間接）
- `go4.org/netipx` v0.0.0-20220812043211-3cc044ffd68d（間接）

## 貢獻

歡迎貢獻！請隨時提交 Pull Request。

## 授權

此專案根據 MIT 授權進行授權 - 有關詳細資訊，請參閱 [LICENSE](../LICENSE) 檔案。

## 致謝

- [DN42](https://dn42.dev/) 社群維護註冊表
- [MaxMind](https://www.maxmind.com/) 提供 MMDB 格式和程式庫
- DN42 註冊表維護者和貢獻者
