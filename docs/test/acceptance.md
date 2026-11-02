# 行为验收清单

对照 Clash Verge Rev v2.5.9 `160bb4fe`（源码路径相对该仓库）。clash-go 按 1:1 复刻时，每条 ID 对应至少一个 `go test`。截图清单在 [docs/spec/checklist.md](../spec/checklist.md)。

怎么测：

| 标记 | 含义 |
| --- | --- |
| unit | 纯逻辑或 `httptest`，可并行 |
| embed | 内嵌 mihomo，`//go:build integration`，禁止 `t.Parallel` |
| ui | `ui.NewTester`，注入假服务，无窗口 |

## 1. 导入远程订阅 `SUB-*`

源：`src-tauri/src/config/prfitem.rs`（`PrfItem::from_url`）、`src-tauri/src/utils/network.rs`。

夹具：`httptest.NewServer` 返回自定义头和 Clash YAML，断言 Method、UA、超时和解析后的字段。规格页：[pages/profiles.md](../spec/pages/profiles.md) §3 导入行、§4 卡片；流量与到期也出现在 [pages/home.md](../spec/pages/home.md) §5。

| ID | 预期（相对 Verge Rev） | 怎么测 | 规格 |
| --- | --- | --- | --- |
| SUB-01 | 未指定 `option.user_agent` 时请求带稳定、非空的 `User-Agent`。Verge 为 `clash-verge/v{version}`；clash-go 为 `clash-go/v{version}`，且可配置 | unit | profiles §3 |
| SUB-02 | `option.user_agent = "X"` 时只发送 `X`，不追加默认 UA | unit | profiles §3 |
| SUB-03 | `timeout_seconds` 缺省 20 秒。超时返回 error，不写入 profile | unit | profiles §3 |
| SUB-04 | `self_proxy=true` 走本机代理；`with_proxy=true` 走系统代理；都为 false 则直连（`ProxyType::None`） | unit（注入 transport；真系统代理不在默认测试） | profiles §3。失败重试走 Clash 代理是界面反馈，见同节 |
| SUB-05 | HTTP 非成功状态返回 error，文案含 status code，不落盘 | unit | profiles §3 |
| SUB-06 | body 不是 YAML Mapping 则 error；没有 `proxies` 且没有 `proxy-providers` 则 error | unit | profiles §3 |
| SUB-07 | 正文以 U+FEFF 开头时先去掉再解析 | unit | — |
| SUB-08 | 头名大小写不敏感，以 `subscription-userinfo` 结尾，前缀为空或 `-`（含 `x-amz-meta-subscription-userinfo`）。解析 `upload` / `download` / `total` / `expire` 写入 `PrfExtra`；缺字段为 0 | unit | 卡片流量、到期：profiles §4，home §5 |
| SUB-09 | 未手动设置 `update_interval` 时，读 `profile-update-interval`（小时）×60 得分钟。非法数字忽略 | unit | profiles：间隔 &lt; 1440 分钟仍可保存并出橙色提示（UI）；本条只断言小时→分钟 |
| SUB-10 | `profile-web-page-url` 为合法 http(s) 且有 host 时写入 `home`；非法则 `home` 为空 | unit | profiles 右键「首页」仅在有 `home` 时出现 |
| SUB-11 | 名称优先 `filename*`（RFC 5987），其次 `filename`；都没有则用 URL 末段，再没有则 `"Remote File"` | unit | profiles 卡片标题 |
| SUB-12 | 生成 `uid`（远程前缀 `R`）、`file={uid}.yaml`、`itype=remote`、`updated` 约等于现在；并挂上空的 merge / script / rules / proxies / groups 子项 uid | unit，`t.TempDir` | profiles §4 |
| SUB-13 | `raise_short_update_intervals(min)` 把过短的更新间隔抬到最小值 | unit | 与 SUB-09 的橙色提示分开：提示不阻止保存，本条是防刷下限 |

## 2. 切换节点与 selected `SEL-*`

源：`src-tauri/src/feat/profile.rs`（`switch_proxy_node` / `record_switched_node`）、`src-tauri/src/cmd/proxy.rs`；`PrfSelected { name, now }`。核心：`hub/executor` 的 `patchSelectGroup`，`component/profile/cachefile`。`profile.StoreSelected` 控制是否持久化。

| ID | 预期（相对 Verge Rev） | 怎么测 | 规格 |
| --- | --- | --- | --- |
| SEL-01 | `PUT /proxies/{group}`（或等价调用）成功后，该 profile 的 `selected` 含 `{name: group, now: node}` | unit 写字段；embed 走真实 PUT | [pages/proxies.md](../spec/pages/proxies.md) §4.7 点击行切换 `now` |
| SEL-02 | 同一 group 再次切换只覆盖 `now`，不追加第二条 | unit | proxies §4.7 |
| SEL-03 | `StoreSelected` 开启时 cachefile 写入 group→node；再次 `ApplyConfig` 后 `patchSelectGroup` 恢复 ForceSet | embed | proxies §4.7 选中态在重载后仍在 |
| SEL-04 | 顺序是先核心切换，再记 profiles。核心失败不写 profiles。profiles 写失败只 warn，核心选择仍生效 | unit：假核心 + 注入写失败 | — |
| SEL-05 | `forget_selected_node` 删除该 group 的记录 | unit | — |
| SEL-06 | 切换 `profiles.current` 后走 enhance + Apply，并恢复该 profile 自己的 `selected` | embed | [pages/profiles.md](../spec/pages/profiles.md) §4 点击卡片切换当前订阅 |

## 3. 测延迟 `DLY-*`

源：`src/services/delay.ts`、`src/utils/delay.ts`；`verge.yaml` 的 `default_latency_test` / `default_latency_timeout`。核心：`GET /proxies/{name}/delay?url=&timeout=`。

| ID | 预期（相对 Verge Rev） | 怎么测 | 规格 |
| --- | --- | --- | --- |
| DLY-01 | 组未设置测速 URL 时用 `http://cp.cloudflare.com/generate_204` | unit，断言传给核心的 `url` | proxies §6；[pages/home.md](../spec/pages/home.md) §6 刷新延迟 |
| DLY-02 | 非法 URL 会删除该组的 URL 映射（静默 delete，不保留旧值）。落地时核对是否与产品预期一致 | unit | — |
| DLY-03 | 默认超时 10000 ms（`DEFAULT_DELAY_TIMEOUT`）。传给核心的 timeout 与此一致，或读 `default_latency_timeout`；该字段非法时回落 10000 | unit | proxies §4.6 排序用同一超时阈值 |
| DLY-04 | 分类：`-2` → testing；`<0` → untested；`0` 或 `≥ timeout` → timeout；`>1e5` → error；其余 → measured | unit | proxies §4.7：`-2` 为加载，`>0` 为彩色数字。§6：单条失败记 `1e6`（错误色） |
| DLY-05 | 排序优先级：measured &lt; timeout &lt; error &lt; testing &lt; untested。同为 measured 时按数值升序 | unit | proxies §4.6 排序态 1；home §6 节点下拉同一顺序 |
| DLY-06 | DelayManager 缓存 30 分钟 | unit，注入时钟 | — |
| DLY-07 | `type: direct` 对本地 204 测延迟，得到有限正数 | embed，`httptest` 返回 204，不访问公网 | proxies §6 |

## 4. 系统代理与 bypass `SYS-*`

源：`src-tauri/src/core/sysopt.rs`，`IVerge.enable_system_proxy`。默认测试用接口 + fake；`networksetup` 真写放 `//go:build integration` 的 macOS 用例。

| ID | 预期（相对 Verge Rev） | 怎么测 | 规格 |
| --- | --- | --- | --- |
| SYS-01 | `enable_system_proxy=false` 时关闭 HTTP / HTTPS / SOCKS 和 PAC | unit + fake | [pages/settings.md](../spec/pages/settings.md) §7.1；home §7 角标只在真实开启时出现 |
| SYS-02 | 开启时指向 `proxy_host`（默认 `127.0.0.1`）加 mixed/http 端口 | unit + fake | settings §7.1 主机 placeholder `127.0.0.1` |
| SYS-03 | macOS 默认 bypass：`127.0.0.1,192.168.0.0/16,10.0.0.0/8,172.16.0.0/12,localhost,*.local,*.crashlytics.com,<local>`，分隔符为逗号 | unit `FormatBypass` | settings §7.1 绕过列表。Windows 默认串另含界面上的 `*.local` chip，以各 OS 的格式化函数为准 |
| SYS-04 | custom 为空 → 仅默认；`use_default_bypass=true` 且有 custom → `默认 + 分隔符 + custom`；`use_default=false` → 仅 custom | unit | settings §7.1 |
| SYS-05 | `proxy_auto_config=true` 时写入 PAC URL，并继续维护 bypass | unit + fake | settings §7.1「使用 PAC」 |
| SYS-06 | `enable_proxy_guard` 开启时，周期把被改掉的系统代理写回 | unit：假 Sysproxy + 注入时钟 | settings §7.1 代理守护 |
| SYS-07 | `toggle_system_proxy` 翻转开关并写入 `verge.yaml` | unit | home §7：开关乐观更新，失败回滚。本条断言落盘后的值 |

## 5. TUN `TUN-*`

源：`src-tauri/src/enhance/tun.rs`（`use_tun`）、`feat/tun.rs`、`feat/proxy.rs` 的 `toggle_tun_mode`。

| ID | 预期（相对 Verge Rev） | 怎么测 | 规格 |
| --- | --- | --- | --- |
| TUN-01 | `use_tun(cfg, true)` 得到 `tun.enable=true` | unit | settings §7.2 |
| TUN-02 | enable 且 DNS `enhanced-mode` 为 fake-ip 或缺省时，强制 `dns.enable=true`；缺省补 `enhanced-mode=fake-ip`、`fake-ip-range=198.18.0.1/16`；IPv6 时补 `fake-ip-range6=2001:2::0/64` | unit | 与 fake-ip 的耦合在 enhance，对话框字段在 settings §7.2。enhance 缺省补的 `fake-ip-range6=2001:2::0/64` 与设置 DNS 对话框 placeholder（[pages/settings.md](../spec/pages/settings.md) §7.4 `fdfe:dcba:9876::1/64`）本来就可以不同；落地以 enhance 源码为准，不要拿对话框 placeholder 覆盖 enhance 缺省 |
| TUN-03 | `use_tun(cfg, false)` 得到 `tun.enable=false`，不改 DNS fake-ip | unit | settings §7.2 |
| TUN-04 | `toggle_tun_mode` 写入 `enable_tun_mode`（`not_save_file` 时跳过写入），并触发 enhance 与重载 | unit 断言字段和调用；embed 断言重载后的 `tun.enable` | home §7 开关。服务不可用时 UI 走安装流程，不直接当成功写入 |
| TUN-05 | 无特权时 reconcile 关闭 TUN 或给出提示（macOS 验收） | integration，darwin；需要 root 的路径 `t.Skip` | home §7「TUN 无服务」 |

## 6. 规则 / 连接 / 日志 `UI-*`

源：`src/pages/rules.tsx`、`connections.tsx`、`logs.tsx`；核心 REST `/rules`、`/connections`，以及日志 WebSocket。无窗口测注入假 store，断言列表文本和按钮。

| ID | 预期（相对 Verge Rev） | 怎么测 | 规格 |
| --- | --- | --- | --- |
| UI-R01 | 展示核心规则的条数和文案；搜索过滤不区分大小写 | ui；embed 可另断言 `/rules` 非空 | [outlines/rules.md](../spec/outlines/rules.md) |
| UI-C01 | 订阅连接；展示 host / process / rule / download / upload；可关闭单条和全部 | ui：一行假数据，点击「关闭全部」调用 `CloseAll` | [outlines/connections.md](../spec/outlines/connections.md) |
| UI-C02 | `auto_close_connection` 开启时，切换节点或 profile 后关闭旧连接 | unit 或 embed：假连接服务收到关闭 | [pages/settings.md](../spec/pages/settings.md) §7.13 `autoCloseConnections`。杂项对话框已有该开关；行为仍按 Verge：开启后切换节点或 profile 关闭旧连接 |
| UI-L01 | 订阅核心日志；按 level 过滤；可清空本地缓冲 | ui | [outlines/logs.md](../spec/outlines/logs.md) |

## 7. 配置布局与 enhance `CFG-*`

源：`src-tauri/src/utils/dirs.rs`、`enhance/mod.rs`、`enhance/merge.rs`、`enhance/script.rs`。

相对 app home：

| 路径 | 用途 |
| --- | --- |
| `config.yaml` | mihomo 运行时基础覆盖 |
| `verge.yaml` | 客户端设置 |
| `profiles.yaml` | 订阅索引 |
| `profiles/{uid}.yaml` | profile / merge / script / rules 等内容 |
| `logs/` | 应用与 sidecar 日志 |

Enhance 顺序：当前 profile YAML → profile merge → profile script → seq(rules / proxies / groups) → global merge / script → builtin → `use_tun` → 端口与控制器覆盖 → 写出 runtime → `ApplyConfig`。

| ID | 预期（相对 Verge Rev） | 怎么测 | 规格 |
| --- | --- | --- | --- |
| CFG-01 | 首次启动补齐缺省的 `verge.yaml`、`profiles.yaml`、`config.yaml` | unit，`t.TempDir` | 只读运行时配置：settings §7.18、profiles 的 ConfigViewer |
| CFG-02 | mapping 递归合并；`dns` 子键 extend；`hosts` 整表替换 | unit | — |
| CFG-03 | `prepend-rules` / `append-rules` 等经 seq 插到规则列表两端，与 Verge 的 seq 一致 | unit | — |
| CFG-04 | JS `main` 接收 config 对象并返回修改后的对象。语法错误进入 ResultLog，不应用坏配置 | unit | profiles 底栏全局 Merge / Script 卡 |
| CFG-05 | enhance 失败时保留上一份成功的 runtime | unit | — |
| CFG-06 | `verge_mixed_port` 等写入最终 config 的 `mixed-port` | unit | settings §7.5。重复端口不能保存是保存前校验，本条是写入结果 |

## 8. 最小端到端 `E2E-*`

假订阅 + 仅 `type: direct` 的核心配置。写法见 [mihomo-embed.md](./mihomo-embed.md)。

| ID | 步骤 | 预期 | 怎么测 | 规格 |
| --- | --- | --- | --- | --- |
| E2E-01 | 假订阅返回 YAML 和 userinfo 头，然后 Import | `profiles` 多一项 remote，extra 与 SUB-08 一致 | unit 或 embed | profiles §3–§4 |
| E2E-02 | 内嵌 mihomo，选中 direct 节点再测延迟 | `delay > 0`，`selected` 已落盘 | embed | proxies §4.7、§6 |
| E2E-03 | 再次 `ApplyConfig` | 选中节点仍是上次的 `now` | embed | proxies §4.7、SEL-03 |
