# mygo 无窗口视图测试

包 `github.com/egoist/mygo/ui`（调研 **v0.7.4**）。文档：<https://github.com/egoist/mygo/blob/main/docs/ui/testing.md>。

`ui.NewTester(view, width, height)` 在内存里画帧，不需要真实窗口或 Display。Linux amd64 上 `TestHeadless` 已通过（约 50ms）。官方示例：`examples/counter-native/main_test.go`。

对应 [acceptance.md](./acceptance.md) 里标记为 ui 的条目，以及 [priority.md](./priority.md) 的 Profiles / Proxies 视图。

## 最小写法

```go
func TestCounter(t *testing.T) {
    app := &counter{}
    tt := ui.NewTester(app.view, 320, 240)
    if err := tt.Click("Increment"); err != nil {
        t.Fatal(err)
    }
    if app.n != 1 || !tt.HasText("1") {
        t.Errorf("count %d texts %q", app.n, tt.Texts())
    }
}
```

## Tester

| 类别 | API |
| --- | --- |
| 指针 | `Click`、`ClickAt`、`ClickWith`、`Press` / `Move` / `Release`、`Scroll`、`Wheel` |
| 键盘 | `Key`、`TypeKey`、`Type`、`Compose`、`Command` |
| 菜单 | `RightClick`、`Menu`、`ChooseMenuItem` |
| 环境 | `SetSize` / `Scale` / `Zoom`、`SetDark`、`SetPreferences`、`SetDirection`、`SetFocused`、`SetWindowActive` |
| 断言 | `Texts`、`HasText`、`Find`、`Focused`、`Cursor`、`Clipboard`、`OpenedURLs`、`Announcements` |
| 帧 | `Frame`、`Image`、`FramePending`、`Advance(d)` |

有可见文字时用文本查找；图标按钮用 `Label`。同一父节点下重复的 `Key` 会 panic。动画用 `Advance`，测试里不用 `time.Sleep`。

## 页面怎么拆

视图是纯函数，状态和 mihomo / 磁盘分开：

```go
func (s *State) view(c *ui.Context) { /* ... */ }
```

用接口注入假实现：`ProxyService`、`ProfileStore`。测试不启动核心。

| 页面 | 断言 |
| --- | --- |
| Profiles | 导入成功后 `HasText(订阅名)`；失败时 `HasText` 错误文案；fake `Import` 被调用 |
| Proxies | `Click("n1")` 后 fake 收到 `Select("PROXY", "n1")`；测速显示毫秒或超时文案 |
| Connections / Rules / Logs | 见 `UI-C01`、`UI-R01`、`UI-L01` |

像素（行高、3px 选中条、色值）用截图对 [docs/spec](../spec/checklist.md)，不放进 `HasText`。

真窗口的端到端用 automation 插件另起进程，放 macOS job，不进默认 `go test`。

## CI

- Linux 上直接 `go test` 这个包，不需要 xvfb。
- 启动真实窗口或 cgo 的文件用 `//go:build` 与 view 测试拆开。
- `go mod tidy` 要带上 `go-text/typesetting` 等传递依赖，缺 `go.sum` 时 setup 会失败。
