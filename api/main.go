package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	flagHost       string
	flagPort       int
	flagMmdbSource string
)

func init() {
	// 命令行参数
	flag.StringVar(&flagHost, "host", "0.0.0.0", "Server host")
	flag.IntVar(&flagPort, "port", 8080, "Server port")
	flag.StringVar(&flagMmdbSource, "mmdb-source", "github", "MMDB source: github or mirror")
	flag.Parse()
}

func queryIP(ipStr string) gin.H {
	reader := getReader()
	if reader == nil {
		return gin.H{"error": "db_not_ready"}
	}

	ipaddr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return gin.H{"error": "invalid_ip"}
	}

	record, err := reader.ASN(ipaddr)
	if err != nil || record == nil {
		return gin.H{"error": "internal_error"}
	}

	return gin.H{
		"ip":                             ipStr,
		"autonomous_system_number":       record.AutonomousSystemNumber,
		"autonomous_system_organization": record.AutonomousSystemOrganization,
	}
}

func respond(c *gin.Context, result gin.H) {
	ua := c.GetHeader("User-Agent")
	isCurl := strings.Contains(strings.ToLower(ua), "curl")

	outputFormat := "json"

	inputFormat, ok := c.GetQuery("f")
	if (ok && inputFormat == "text") || (!ok && isCurl) {
		outputFormat = "text"
	}

	if outputFormat == "text" {
		if err, ok := result["error"]; ok {
			c.String(http.StatusBadRequest, "Error: %s\n", err)
			return
		}

		c.String(http.StatusOK,
			`IP                 : %s
ASN                : %d
ASN Organization   : %s
`,
			result["ip"],
			result["autonomous_system_number"],
			result["autonomous_system_organization"],
		)
	} else {
		if _, ok := result["error"]; ok {
			c.JSON(http.StatusBadRequest, result)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func main() {
	// 捕获退出信号
	stopCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sourceUrl := mmdbURL
	if flagMmdbSource == "mirror" {
		sourceUrl = mmdbMirrorURL
	}

	if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
		fmt.Printf("下载 MMDB 文件: %s\n", sourceUrl)
		if err := downloadMMDB(sourceUrl, localFilePath); err != nil {
			panic(err)
		}
	}

	r, err := loadMMDB(localFilePath)
	if err != nil {
		panic(err)
	}
	currentReader.Store(r)

	// 后台更新
	go updateLoop(sourceUrl, stopCh)

	// Gin 服务器
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("X-Powered-By", "https://github.com/rdp-studio/dn42-geoasn")
		c.Next()
	})

	// 根路径：返回客户端 IP 信息
	router.GET("/", func(c *gin.Context) {
		clientIP := c.ClientIP()
		result := queryIP(clientIP)
		respond(c, result)
	})

	// 查询指定 IP
	router.GET("/q", func(c *gin.Context) {
		ipStr := strings.TrimSpace(c.Query("ip"))
		if ipStr == "" {
			respond(c, gin.H{"error": "missing_ip"})
			return
		}
		result := queryIP(ipStr)
		respond(c, result)
	})

	serverAddr := fmt.Sprintf("%s:%d", flagHost, flagPort)
	fmt.Printf("启动服务器: http://%s\n", serverAddr)
	go func() {
		if err := router.Run(serverAddr); err != nil {
			panic(err)
		}
	}()

	// 等待退出
	<-sigCh
	fmt.Println("收到退出信号，关闭 reader 并退出")
	close(stopCh)

	finalReader := getReader()
	if finalReader != nil {
		finalReader.Close()
	}
}
