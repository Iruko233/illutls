# illutls

[English](#english) | [简体中文](#简体中文)

---

<h2 id="english">English</h2>

`illutls` is a Go HTTP client designed to simulate real browser TLS fingerprints and HTTP headers.

It combines `utls` (for TLS ClientHello simulation) with `fhttp` (for HTTP/2 frame-level fingerprint control) to present a complete, authentic browser identity on every connection.

### Features

- **TLS Fingerprinting**: Accurately simulates JA3/JA4 signatures using `utls`.
- **HTTP/2 Fingerprinting**: Mimics real browser HTTP/2 settings, frame order, and pseudo-header order using `fhttp`.
- **Pre-configured Profiles**: Ships with real-world browser profiles covering Chrome, Firefox, Edge, and Safari across Windows, macOS, Linux, Android, and iOS.
- **Concurrency Safe**: `Client` is entirely safe for concurrent use by multiple goroutines.
- **Proxy Support**: Easily route your simulated requests through HTTP/SOCKS5 proxies.

### Installation

```bash
go get github.com/Iruko233/illutls
```

### Quick Start

Here's a simple example of how to make a request imitating Google Chrome on Windows:

```go
package main

import (
	"fmt"
	"io"
	"log"

	"github.com/Iruko233/illutls"
	_ "github.com/Iruko233/illutls/profiles" // Import built-in profiles
)

func main() {
	// Create a new client mimicking Chrome 149 on Windows
	client, err := illutls.New(illutls.WithProfile("chrome-149-windows-10"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Make a GET request
	resp, err := client.Get("https://tls.browserleaks.com/json")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	fmt.Println(string(body))
}
```

### Available Profiles

`illutls` includes various pre-built profiles. [See the full list of supported profiles here.](PROFILES.md)

You can view all available profiles programmatically:

```go
profiles := illutls.ListProfiles()
for _, p := range profiles {
    fmt.Println(p)
}
```

By default, if you don't specify a profile using `illutls.WithProfile()`, it will default to `"chrome-149-windows-10"`.

### Advanced Usage

#### Using Proxies

You can easily route your traffic through a proxy server using `WithProxy`:

```go
client, err := illutls.New(
    illutls.WithProfile("chrome-149-windows-10"),
    illutls.WithProxy("http://user:pass@proxy.example.com:8080"),
)
```

#### Custom Requests & Headers

If you need to send a POST request or customize headers, use `NewRequest` combined with `client.Do()`:

```go
req, err := client.NewRequest("POST", "https://example.com/api", bodyReader)
if err != nil {
    log.Fatal(err)
}

// Add your custom headers (the client already pre-populates the browser's default headers)
req.Header.Set("Authorization", "Bearer your_token")

resp, err := client.Do(req)
```

---

<h2 id="简体中文">简体中文</h2>

`illutls` 是一个 Go 语言 HTTP 客户端，旨在模拟真实浏览器的 TLS 指纹和 HTTP 请求头。

它将 `utls`（用于模拟 TLS ClientHello）和 `fhttp`（用于 HTTP/2 帧级别的指纹控制）结合在一起，在每次连接时都能呈现出完整且真实的浏览器特征。

### 特性

- **TLS 指纹模拟**：使用 `utls` 精确模拟 JA3/JA4 签名。
- **HTTP/2 指纹模拟**：使用 `fhttp` 模拟真实浏览器的 HTTP/2 Settings、帧顺序和伪头（Pseudo-header）顺序。
- **预置配置（Profiles）**：内置了覆盖 Windows, macOS, Linux, Android 和 iOS 平台的真实浏览器特征（包含 Chrome, Firefox, Edge, Safari）。
- **并发安全**：`Client` 完全支持多 Goroutine 并发安全调用。
- **代理支持**：可以轻松将请求通过 HTTP/SOCKS5 代理进行转发。

### 安装

```bash
go get github.com/Iruko233/illutls
```

### 快速开始

下面是一个模拟 Windows 端 Google Chrome 发送请求的简单示例：

```go
package main

import (
	"fmt"
	"io"
	"log"

	"github.com/Iruko233/illutls"
	_ "github.com/Iruko233/illutls/profiles" // 导入内置的浏览器配置
)

func main() {
	// 创建一个模拟 Windows Chrome 149 的客户端
	client, err := illutls.New(illutls.WithProfile("chrome-149-windows-10"))
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 发送 GET 请求
	resp, err := client.Get("https://tls.browserleaks.com/json")
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}

	fmt.Println(string(body))
}
```

### 可用的浏览器配置 (Profiles)

`illutls` 包含了多种预构建的配置。[点击此处查看所有支持的浏览器配置列表。](PROFILES.md)

你可以通过代码查看所有可用的配置名称：

```go
profiles := illutls.ListProfiles()
for _, p := range profiles {
    fmt.Println(p)
}
```

默认情况下，如果你在初始化时不通过 `illutls.WithProfile()` 指定配置，它将默认使用 `"chrome-149-windows-10"`。

### 高级用法

#### 使用代理

你可以使用 `WithProxy` 选项轻松配置代理服务器：

```go
client, err := illutls.New(
    illutls.WithProfile("chrome-149-windows-10"),
    illutls.WithProxy("http://user:pass@proxy.example.com:8080"),
)
```

#### 自定义请求与 Header

如果你需要发送 POST 请求或自定义请求头，请使用 `NewRequest` 并配合 `client.Do()` 使用：

```go
req, err := client.NewRequest("POST", "https://example.com/api", bodyReader)
if err != nil {
    log.Fatal(err)
}

// 添加你自己的请求头（客户端已经提前填充好了浏览器默认的请求头）
req.Header.Set("Authorization", "Bearer your_token")

resp, err := client.Do(req)
```

### 许可证 / License

MIT License
