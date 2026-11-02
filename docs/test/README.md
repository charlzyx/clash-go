# 行为测试

对照 Clash Verge Rev **v2.5.9**（`160bb4fe`，与 [docs/spec](../spec/README.md) 同一参照）。这里只写 `go test` 能断言的行为。截图像素勾选在 [docs/spec/checklist.md](../spec/checklist.md)。

仓库尚无应用代码。下列命令是包落地之后的跑法，现在执行会因没有 `go.mod` 而失败。

## 和界面规格的分工

| | 界面规格 `docs/spec/` | 行为测试 `docs/test/` |
| --- | --- | --- |
| 回答 | 长什么样、点哪里、文案 key | 做完之后状态和数据对不对 |
| 验收 | 截图对照 [checklist.md](../spec/checklist.md) | 单测 / 嵌入 mihomo / 无窗口 UI |
| 典型条目 | 节点行高 40、选中左边 3px 条 | `selected` 写入、延迟分类、订阅头解析 |

同一功能两边都写时，本目录用稳定 ID（`SUB-*`、`SEL-*` 等），并注明对应的规格页。像素项不在这里重复。

## 目录

| 文件 | 内容 |
| --- | --- |
| [acceptance.md](./acceptance.md) | 行为验收清单（断言 ID） |
| [mihomo-embed.md](./mihomo-embed.md) | 内嵌 mihomo：`hub.Parse`、`TempDir`、假订阅 |
| [ui-headless.md](./ui-headless.md) | mygo `ui.NewTester` 无窗口视图测试 |
| [priority.md](./priority.md) | M1 落地后先补哪些测试 |
| [ci-notes.md](./ci-notes.md) | 默认 CI 跑什么；本目录不提交 workflow |

## 代码落地后怎么跑

默认（与 [ci-notes.md](./ci-notes.md) 一致，macOS arm64）：

```bash
go vet ./...
go test -race -count=1 -timeout 15m \
  -covermode=atomic -coverprofile=coverage.out ./...
```

嵌入核心（全局单例，不要并行）：

```bash
go test -tags=integration -race -count=1 -p 1 -timeout 15m ./...
```

分层：纯解析可 `t.Parallel`；碰 `hub.Parse` / `ApplyConfig` 的用例串行，见 [mihomo-embed.md](./mihomo-embed.md)。无窗口页面见 [ui-headless.md](./ui-headless.md)。真窗口、系统代理真写、TUN 用 `//go:build integration`，不进默认 `go test`。
