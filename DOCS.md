# illutls Advanced Features

[简体中文](#简体中文) | [English](#english)

---

<h2 id="简体中文">简体中文</h2>

本文档介绍了 `illutls` 中关于指纹生成和语言标头解析的高级配置选项。

### 1. 动态指纹 Seed
`WithDynamicProfile` 的 `seed` 参数支持 `any` 类型（如 `string`、`int`、`int64`）。当传入字符串时，底层将使用 FNV-1a 算法将其哈希为确定的 `int64` seed。这确保了相同的输入会始终生成一致的 TLS 扩展顺序、锁定稳定的 JA4 哈希以及默认的随机语言。

如果不传入其他参数，底层会根据该 seed **确定性地为你随机分配** 一个固定的操作系统（如 Windows、Mac）和 Chrome 大版本号。你也可以选择手动指定。

```go
proxyUrl := "http://user:pass@1.2.3.4:1080"

client, _ := illutls.New(
    // 用法 1：极简模式。操作系统和 Chrome 版本号由 Seed 确定性生成
    illutls.WithDynamicProfile(proxyUrl), 
    
    // 用法 2：强制干预。手工指定 ("windows", "mac", "linux") 和大版本号 (如 145)
    // illutls.WithDynamicProfile(proxyUrl, "windows", 145), 
    
    illutls.WithProxy(proxyUrl),
)
```

### 2. IP 物理位置语言解析 (内置 MMDB)
`illutls` 内置了 MaxMind GeoLite2-Country 数据库。通过向 `WithLanguage` 传入 `"auto"`，客户端会自动解析代理 IP 的所在国家。随后，引擎将基于该国家代码匹配主语言，并通过 Chromium 的原生算法生成带权重的 `Accept-Language` 标头。

```go
client, _ := illutls.New(
    illutls.WithDynamicProfile(proxyUrl), 
    illutls.WithProxy(proxyUrl),
    illutls.WithLanguage("auto"),
)
```

### 3. ISO 国家代码解析
直接向 `WithLanguage` 传入 2 字母的 ISO 3166-1 alpha-2 国家代码时，底层将依据 OS 区域设置逻辑自动构建基础语言列表（例如，`FR` 将被展开为 `fr-FR,fr,en-US,en`），并套用 Chromium 的 Q 权重递减算法。

```go
client, _ := illutls.New(
    illutls.WithDynamicProfile("seed_string"), 
    illutls.WithLanguage("FR"),
)
```

### 4. 外部 GeoIP 数据库
若需使用外部或更新的 MaxMind MMDB 文件，可直接将以 `.mmdb` 结尾的文件路径传入 `WithLanguage`，客户端将优先加载该外部数据库。

```go
client, _ := illutls.New(
    illutls.WithDynamicProfile(proxyUrl), 
    illutls.WithProxy(proxyUrl),
    illutls.WithLanguage("/data/custom_geoip.mmdb"),
)
```

### 5. 危险的 JA4 持久化随机伪造
默认情况下，`illutls` 会完全模拟 Google Chrome 的 `Cipher Suites` 密码套件排列顺序，这会生成全球公认的合法 Chrome JA4 哈希值（密码套件部分为 `8daaf6152771`）。像 Cloudflare 或 Akamai 这样的高级 WAF 通常会将其纳入常规基线。

然而，在某些复杂的网络环境中，部分自建的防火墙系统可能会因为“短时间内有太多请求使用相同的 JA4 哈希”而盲目地进行频次封锁。为了绕过这种防御机制，你可以开启一个危险的开关，它将绑定到动态配置（Dynamic Profile）的 Seed 上，产生持久化且确定性的 JA4 变异：

```go
client, _ := illutls.New(
    // 生成的变异 JA4 哈希将与 "seed_string" 永久绑定
    illutls.WithDynamicProfile("seed_string"), 
    illutls.WithDangerousJA4Randomization(), 
)
```

> **警告：** 这是一个极度危险的功能！它会强制从 ClientHello 中删减有效的密码套件，从而构造出一个畸形的、非标准的 TLS 指纹。虽然它可以绕过部分基于频次封锁的防火墙，但它**可能被企业级高级 WAF 拦截**。请务必谨慎使用！

---

<h2 id="english">English</h2>

This document outlines the advanced configuration options for profile generation and language header resolution in `illutls`.

### 1. Dynamic Profile Seeds
`WithDynamicProfile` accepts a `seed` parameter of type `any` (e.g., `string`, `int`, `int64`). A string input is hashed using FNV-1a to generate a deterministic `int64` seed. This ensures that the same input consistently yields the same TLS extensions order, JA4 hash, and default pseudo-random language.

By default, the OS platform and Chrome major version are also deterministically chosen based on the seed. You can explicitly override them by providing additional arguments.

```go
proxyUrl := "http://user:pass@1.2.3.4:1080"

client, _ := illutls.New(
    // Usage 1: Deterministic selection of OS and Chrome version based on the seed
    illutls.WithDynamicProfile(proxyUrl), 
    
    // Usage 2: Explicitly forcing Platform and Chrome Version
    // illutls.WithDynamicProfile(proxyUrl, "windows", 145), 
    
    illutls.WithProxy(proxyUrl),
)
```

### 2. Auto Geo-Language (Embedded MMDB)
`illutls` includes an embedded MaxMind GeoLite2-Country database. Passing `"auto"` to `WithLanguage` configures the client to resolve the proxy IP's country code internally. The country code is then mapped to its primary language and processed through Chromium's native `Accept-Language` generation algorithm.

```go
client, _ := illutls.New(
    illutls.WithDynamicProfile(proxyUrl), 
    illutls.WithProxy(proxyUrl),
    illutls.WithLanguage("auto"),
)
```

### 3. ISO Country Code Resolution
Passing a 2-letter ISO 3166-1 alpha-2 country code to `WithLanguage` directly constructs the raw language list based on OS locale fallbacks (e.g., `FR` resolves to `fr-FR,fr,en-US,en`) and applies Chromium's q-value weighting algorithm.

```go
client, _ := illutls.New(
    illutls.WithDynamicProfile("seed_string"), 
    illutls.WithLanguage("FR"),
)
```

### 4. External GeoIP Database
To use an external or updated MaxMind MMDB file instead of the embedded database, pass the file path (ending in `.mmdb`) to `WithLanguage`.

```go
client, _ := illutls.New(
    illutls.WithDynamicProfile(proxyUrl), 
    illutls.WithProxy(proxyUrl),
    illutls.WithLanguage("/data/custom_geoip.mmdb"),
)
```

### 5. Dangerous JA4 Randomization
By default, `illutls` mimics the exact `Cipher Suites` order of Google Chrome, resulting in the universally recognized Chrome JA4 hash (`8daaf6152771` for the cipher suite component). Advanced WAFs (like Cloudflare or Akamai) expect this exact hash.

However, primitive or self-built WAFs might blindly rate-limit or block requests simply because "too many requests share the same JA4 hash". To bypass this, you can enable persistent, deterministic JA4 randomization tied to the dynamic profile seed:

```go
client, _ := illutls.New(
    // The mutated JA4 hash will be different, but permanently tied to "seed_string"
    illutls.WithDynamicProfile("seed_string"), 
    illutls.WithDangerousJA4Randomization(), 
)
```

> **WARNING:** This is a highly dangerous feature. It drops valid cipher suites from the ClientHello, creating an anomalous, non-standard TLS fingerprint. While it bypasses primitive frequency-based blocks, it **may potentially be intercepted** by advanced enterprise WAFs. Use with extreme caution.
