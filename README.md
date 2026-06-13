# illutls

[简体中文](#简体中文) | [English](#english)

---

<h2 id="简体中文">简体中文</h2>

`illutls` 是一个 Go 语言 HTTP 客户端，旨在模拟真实浏览器的 TLS 指纹和 HTTP 请求头。

它将 `utls`（用于精确实控 TLS ClientHello）和 `fhttp`（用于底层 HTTP/2 帧级别的通信控制）结合在一起，在每次连接时都能完美模拟出完整且真实的浏览器特征。

### 特性

- **TLS 指纹控制**：使用 `utls` 精确实控 ClientHello 参数，从而完美模拟 JA3/JA4 签名。
- **HTTP/2 指纹控制**：使用 `fhttp` 在底层控制 HTTP/2 Settings、帧顺序和伪头（Pseudo-header）顺序，精准模拟真实浏览器行为。
- **预置配置（Profiles）**：内置了覆盖 Windows, macOS, Linux, Android 和 iOS 平台的真实浏览器特征（包含 Chrome, Firefox, Edge, Safari）。
- **并发安全**：`Client` 完全支持多 Goroutine 并发安全调用。
- **TLS 扩展随机化 (Extension Shuffling)**：默认开启，在建立连接时随机打乱 TLS 扩展顺序，在生成动态 JA3 指纹的同时保持 JA4 哈希稳定。可通过 `WithShuffleExtensions(false)` 关闭。
- **危险的 JA4 持久化随机伪造 (Dangerous JA4 Randomization)**：一个极度危险且必须显式开启的选项 (`WithDangerousJA4Randomization()`)。它会通过删减密码套件来持久化随机改变 JA4 哈希。**警告：** 这将破坏纯正的 Chrome 官方指纹，可能会引起 Cloudflare 等高级 WAF 的异常判定，仅建议用于对抗按 JA4 频次进行拦截的自建 WAF。
- **动态指纹 Seed (Dynamic Profile Seeds)**：支持将任意字符串（如代理 IP）哈希转换为确定的长效指纹，确保针对同一输入的指纹特征（包含 TLS 与 HTTP 标头）保持一致。
- **IP 物理位置语言解析 (Auto Geo-Language)**：内置无依赖的 MaxMind MMDB。启用后可根据代理 IP 所在的物理国家代码，运用 Chromium 算法逻辑动态生成对应的 `Accept-Language` 标头。
- **Fetch 标头自适应 (Auto-Adaptive Fetch Headers)**：针对 API/XHR 请求自动切换与更新 `sec-fetch-*` 标头族。
- **ECH 扩展兼容**：为 ECH 扩展 (65037) 注入符合规范的 GREASE 占位数据，修复部分 WAF（如 Cloudflare）因负载为空导致的解析失败问题。
- **代理支持**：支持将请求通过 HTTP/SOCKS5 代理进行转发。

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

### 可用的浏览器配置 (内置 Profiles)

`illutls` 包含了多种预构建的配置。**[👉 点击此处查看所有支持的内置浏览器配置列表 (PROFILES.md)](PROFILES.md)**

你可以通过代码查看所有可用的配置名称：

```go
profiles := illutls.ListProfiles()
for _, p := range profiles {
    fmt.Println(p)
}
```

默认情况下，如果你在初始化时不通过 `illutls.WithProfile()` 指定配置，它将默认使用 `"chrome-149-windows-10"`。

### 动态指纹 Seed 与语言地域解析 (Geo-Language)

除了使用内置配置，当业务涉及大量代理 IP 轮换时，通常需要保持指纹的持久性（即同一个 IP 始终对应同一个浏览器指纹特征，包括 JA4、UA、语种等）。

`illutls` 支持将任意字符串（如代理 URL）哈希为确定的动态指纹，并内置了无依赖的 MaxMind MMDB，可根据代理 IP 自动解析物理国家并生成与之完美匹配的原生 `Accept-Language` 请求头。

**[👉 查阅详细的高级指纹与语言配置选项 (DOCS.md)](DOCS.md)** 了解如何使用动态指纹 Seed、物理位置自适应、手工指定 ISO 国家代码或挂载外部数据库。

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

// 客户端已填充浏览器默认请求头，可在此追加或覆盖自定义请求头
req.Header.Set("Authorization", "Bearer your_token")

resp, err := client.Do(req)
```

---

<h2 id="english">English</h2>

`illutls` is a Go HTTP client designed to simulate real browser TLS fingerprints and HTTP headers.

It combines `utls` (for TLS ClientHello simulation) with `fhttp` (for HTTP/2 frame-level fingerprint control) to present a complete, authentic browser identity on every connection.

### Features

- **TLS Fingerprinting**: Precisely controls ClientHello parameters using `utls` to perfectly simulate JA3/JA4 signatures.
- **HTTP/2 Fingerprinting**: Controls HTTP/2 settings, frame order, and pseudo-header order at the lowest level using `fhttp`, accurately simulating real browser behavior.
- **Pre-configured Profiles**: Ships with real-world browser profiles covering Chrome, Firefox, Edge, and Safari across Windows, macOS, Linux, Android, and iOS.
- **Concurrency Safe**: `Client` is entirely safe for concurrent use by multiple goroutines.
- **Extension Shuffling**: Enabled by default to randomize the order of TLS extensions per connection. This simulates real browser extension shuffling to produce variable JA3 hashes while maintaining stable JA4 signatures. Can be disabled via `WithShuffleExtensions(false)`.
- **Dangerous JA4 Randomization**: An explicitly opt-in feature (`WithDangerousJA4Randomization()`) that persistently mutates the Cipher Suites list to bypass primitive, frequency-based JA4 blocklists. **WARNING:** This breaks the authentic Chrome fingerprint and may potentially be flagged by advanced WAFs like Cloudflare.
- **Dynamic Profile Seeds**: Generate completely deterministic profiles and sticky fingerprints tied to any string (e.g. a Proxy URL).
- **Auto Geo-Language**: An embedded, zero-dependency MaxMind MMDB provides instant, mathematical deduction of Chrome Accept-Language headers perfectly matched to the proxy IP's physical location.
- **Auto-Adaptive Fetch Headers**: Automatically switches `sec-fetch-*` headers for API/XHR requests on the fly without breaking the connection pool.
- **ECH Support**: Injects compliant GREASE payloads into the ECH extension (65037) to prevent TLS parsing errors on strict WAFs (e.g., Cloudflare) while preserving fingerprint integrity.
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

### Available Profiles (Built-in)

`illutls` includes various pre-built profiles. **[👉 See the full list of supported profiles (PROFILES.md)](PROFILES.md)**

You can view all available profiles programmatically:

```go
profiles := illutls.ListProfiles()
for _, p := range profiles {
    fmt.Println(p)
}
```

By default, if you don't specify a profile using `illutls.WithProfile()`, it will default to `"chrome-149-windows-10"`.

### Dynamic Profile Seeds & Geo-Language

While built-in profiles are useful, rotating proxies often requires sticky fingerprints—the same proxy IP should consistently produce the exact same browser fingerprint (JA4, UA, language) to mimic a persistent user. 

`illutls` supports generating deterministic profiles tied to any string seed (like a proxy URL), and includes an embedded zero-dependency MaxMind MMDB for perfectly matching the `Accept-Language` header to the proxy's physical location.

**[👉 Read the Advanced Features Documentation (DOCS.md)](DOCS.md)** for detailed configurations on dynamic seeds, auto geo-language resolution, ISO country overrides, and external database loading.

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

### 许可证 / License

MIT License
