# 内嵌 mihomo

模块 `github.com/metacubex/mihomo`，调研锁定 **v1.19.33**（以 `go.mod` 为准）。用于订阅导入之后的选中、延迟和重载，对应 [acceptance.md](./acceptance.md) 的 `SEL-*`、`DLY-07`、`E2E-*`。

## 入口

| API | 用途 |
| --- | --- |
| `constant.SetHomeDir(dir)` | 隔离 cache / geo。测试固定用 `t.TempDir()` |
| `constant.SetConfig(path)` | 默认配置路径 |
| `hub.Parse(configBytes, opts...)` | 解析并 `ApplyConfig`。`nil` bytes 时读 `SetConfig` 的路径 |
| `hub.WithExternalController(addr)` / `WithSecret` | 覆盖控制器地址和密钥 |
| `executor.ParseWithBytes` / `ApplyConfig(cfg, force)` | 更细的重载 |
| `executor.Shutdown()` | `t.Cleanup` 里释放监听 |
| `hub/route` | `/version`、`/proxies`、`PUT /proxies/{group}`、`GET /proxies/{name}/delay`、`/connections`、`/rules` |

`patchSelectGroup` 从 `component/profile/cachefile` 恢复 select 组。是否写入由 `profile.StoreSelected` 决定。

## 全局单例

核心状态在进程内是全局的。嵌入用例禁止 `t.Parallel()`。用包级 `sync.Mutex` 串行，或 `go test -p 1`。并行会抢 HomeDir、控制器端口和当前配置。

Linux amd64 上已跑通的最小路径（2026-11-03）：

1. `t.TempDir()` + `SetHomeDir`
2. `external-controller` 用 `127.0.0.1` 上的空闲端口
3. YAML 只有 `type: direct` 节点、一个 select 组、`MATCH,PROXY`
4. `hub.Parse(nil)` 之后 `GET /version` 为 200

`/version` 里的数字可能和 module tag 不一致，断言 200 和 JSON 形状，不写死版本字符串。

## 假订阅 + DIRECT

端到端只使用 `httptest` 和 `direct`，不访问公网，也不拉 geoip。测试配置关闭 geodata，或预置空 stub。

```go
func TestImportAndSelect(t *testing.T) {
    lockMihomo(t) // 全局锁；不要 t.Parallel
    constant.SetHomeDir(t.TempDir())

    sub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("User-Agent") == "" {
            t.Error("missing UA")
        }
        w.Header().Set("Subscription-Userinfo", "upload=1; download=2; total=3; expire=4")
        w.Header().Set("Profile-Update-Interval", "24") // 小时 → 1440 分钟
        _, _ = w.Write([]byte(`
proxies:
  - {name: n1, type: direct}
proxy-groups:
  - {name: PROXY, type: select, proxies: [n1, DIRECT]}
rules:
  - MATCH,PROXY
`))
    }))
    t.Cleanup(sub.Close)

    // ImportRemote → extra 与 update_interval==1440，文件在 profiles/ 下
    // 随机控制器端口启动核心
    // PUT 切换 PROXY → n1，读 cache 与 profiles selected
    // GET /proxies/n1/delay?url=<本地 204>&timeout=3000
}
```

本地 204：

```go
httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusNoContent)
}))
```

断言用 external-controller。`mixed-port: 0` 仍可能占 fd。

## 分层

1. **纯解析**（不调用 hub）：YAML merge、订阅头、bypass。快，可 `t.Parallel`。
2. **executor**：`ParseWithBytes` + 内存控制器。`-race`，串行。
3. **TUN / 系统代理真写**：`//go:build` 限制 darwin 或 integration；要 root 时 `t.Skip`。
4. 单测不下载 geoip 大文件。

建议 build tag：

```go
//go:build integration

package core_test
```

跑法：`go test -tags=integration -race -count=1 -p 1 ./...`。workflow 由 M1 放进 `.github/workflows/test.yml`，见 [ci-notes.md](./ci-notes.md)。

## 容易踩的点

- 未 `SetHomeDir(t.TempDir())` 会写到用户的 mihomo 配置目录。
- 未 `Shutdown` 时端口占用，下一个用例失败。
- 在嵌入用例上调用 `t.Parallel()` 会和其它用例抢全局配置。
- 把 `/version` 的版本字符串写成与 tag 完全一致。
