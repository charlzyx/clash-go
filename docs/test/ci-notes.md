# CI 约定

本目录只记录测试怎么跑。**不要在这个 PR 里添加** `.github/workflows/test.yml`。workflow 归 M1。

默认平台是 macOS arm64（`macos-14`）。Linux 可以跑纯 Go、无窗口 UI 和 mihomo 嵌入（不写系统代理、不起 TUN）。

## 默认 job

`go` 版本跟仓库 `go.mod`，并满足 mygo 的要求（调研时 mygo 侧为 1.27.3）。

```bash
go vet ./...
go test -race -count=1 -timeout 15m \
  -covermode=atomic -coverprofile=coverage.out \
  ./...
```

覆盖率先做成 artifact（`go tool cover -func=coverage.out`）。等包稳定后再加门槛，建议核心包 ≥70%、UI view ≥50%。测试 job 不放签名密钥。

## 什么放进 integration

默认 `go test ./...` 不含：

- 真实窗口
- 系统代理真写（`networksetup` 等）
- TUN，以及任何要 root 的路径

这些文件加 build tag，例如：

```go
//go:build integration

package core_test
```

系统代理真机可用 `//go:build darwin && integration`。没有权限时 `t.Skip`。

可选的 Linux job 只跑嵌入和无窗口 UI：

```bash
go test -race -count=1 -timeout 15m \
  -tags=integration -p 1 \
  -covermode=atomic -coverprofile=coverage.out \
  ./...
```

`-p 1` 或包级 Mutex 是因为 mihomo 是进程内全局单例，见 [mihomo-embed.md](./mihomo-embed.md)。无窗口 UI 在 Linux runner 上直接测，不需要 xvfb，见 [ui-headless.md](./ui-headless.md)。
