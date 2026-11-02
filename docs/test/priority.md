# M1 落地后先补哪些测试

仓库还没有这些包。M1 放进骨架之后，按「纯逻辑 → 嵌入核心 → 无窗口视图」补断言。ID 定义在 [acceptance.md](./acceptance.md)。

## P0：纯单元

无窗口、无 root，默认可进 `go test -race`。

| 对象 | ID |
| --- | --- |
| `ParseSubscriptionUserInfo` | SUB-08 |
| `ParseUpdateInterval` | SUB-09 |
| `ImportHTTP`（fake RoundTripper / `httptest`） | SUB-01 … SUB-07、SUB-10 … SUB-12 |
| `Merge` / `PrependAppend` | CFG-02、CFG-03 |
| `FormatBypass` | SYS-03、SYS-04 |
| `Classify` / `Compare`（延迟） | DLY-04、DLY-05 |

## P0：嵌入冒烟

`//go:build integration`，全局锁，禁止 `t.Parallel`。见 [mihomo-embed.md](./mihomo-embed.md)。

| 用例 | 断言 |
| --- | --- |
| `TestHubParseDirectVersion` | `hub.Parse` 之后 `GET /version` 为 200 |
| `TestSelectPersistAcrossApply` | PUT 选择 → Shutdown / Parse → 仍 ForceSet（SEL-03、E2E-03） |
| `TestDelayDirectAgainstLocal204` | DLY-07、E2E-02 |

配置里只有 `type: direct`。

## P0：无窗口视图

`ui.NewTester`。见 [ui-headless.md](./ui-headless.md)。名称随实现调整。

| View | 断言 |
| --- | --- |
| `ProfilesView` | 填写 URL 并确认导入；fake `Import` 被调用；列表 `HasText(订阅名)`；失败时出现错误文案 |
| `ProxiesView` | `Click("n1")` → fake `Select("PROXY", "n1")`，选中标记在 |
| `ProxiesView` 测速 | fake Delay 返回后显示毫秒或超时文案 |

## P1（P0 之后）

| 层 | 对象 | ID 或断言 |
| --- | --- | --- |
| unit | `RecordSelected` 与 YAML 往返 | SEL-01、SEL-02、SEL-05 |
| unit | `UseTUN` | TUN-01 … TUN-03 |
| unit | 首次创建配置目录 | CFG-01 |
| embed | `/rules` 非空；有流量后 `/connections` 可见 | UI-R01、UI-C01 的数据侧 |
| ui | 系统代理开关 | fake `Enable(true)`，bypass 文本在 |
| ui | TUN 开关 | fake `ToggleTUN` |
| ui | 连接页空态、一行数据、「关闭全部」 | `CloseAll` |
| ui | 日志按 level 过滤 | 只留下匹配行 |
| ui | 规则搜索 | 过滤后的 `Texts` 是子集 |

## P2 / 暂缓

- 侧栏切换页面：目标页标题出现。
- 真 `networksetup`、TUN 起隧道转发流量、公证和托盘点击：等 macOS 与特权环境。automation 插件后续再补。

## 落地时优先核对的行为

- SUB-08 头前缀（空或 `-`，兼容 `x-amz-meta-subscription-userinfo`）。
- SEL-04：profiles 写失败只 warn，核心选择仍生效。
- TUN-02：`use_tun` 与 fake-ip DNS 一起改。
- DLY-02：非法测速 URL 会静默删掉映射。
