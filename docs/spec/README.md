# clash-go 界面还原规格

本文档把 [Clash Verge Rev](https://github.com/clash-verge-rev/clash-verge-rev) 前端界面整理成可逐项实现、可拿截图比对的规格。实现方（主程）按本文对照做原生 UI（Go + egoist/mygo，不用 webview）；后续截图验收用 [验收清单](./checklist.md)。

只描述界面与交互。不包含 clash-go 的实现代码。

## 参照版本

| 项 | 值 |
| --- | --- |
| 仓库 | `clash-verge-rev/clash-verge-rev` |
| Release | **v2.5.9**（2026-11-02） |
| Commit | `160bb4fe8eaef2bc626163d89622229e5d7c5b59`（`Release: v2.5.9 (#8508)`） |
| 前端入口 | `src/pages/`、`src/components/`、`src/locales/`、`src/assets/styles/` |
| 窗口与托盘 | `src-tauri/src/utils/resolve/window.rs`、`src-tauri/src/core/tray/` |
| 托盘文案 | `crates/clash-verge-i18n/locales/{zh,en}.yml`（对照表 `i18n/tray.md`） |

文中源码路径均相对 Clash Verge Rev 仓库，行号对应该 commit。MUI 的 `sx` 间距单位是 **8px**（`theme.spacing(1) = 8`）；`borderRadius: 1` 等于默认 `shape.borderRadius`（4px），`borderRadius: 2` 等于 8px。

## 目录树

```text
docs/spec/
├── README.md                 本索引：目录、页面清单、优先级、阅读约定
├── checklist.md              验收清单（截图比对用）
├── design-tokens.md          设计令牌（色、字、间距、圆角、阴影、图标）  [P0 已写]
├── assets/
│   └── brand/
│       └── favicon.svg       侧栏品牌标记（egoist/mygo favicon.svg，MIT）
├── pages/
│   ├── home.md               首页                                          [P0 完整]
│   ├── proxies.md            代理                                          [P0 完整]
│   ├── profiles.md           订阅                                          [P0 完整]
│   └── settings.md           设置（含子对话框）                            [P0 完整]
├── outlines/
│   ├── layout.md             整体布局（侧栏、顶栏、流量图、窗口尺寸）      [P1 大纲]
│   ├── connections.md        连接                                          [P1 规格]
│   ├── rules.md              规则                                          [P1 规格]
│   ├── logs.md               日志                                          [P2 规格]
│   ├── unlock.md             测试（解锁测试）                              [P2 规格]
│   └── tray.md               托盘菜单                                      [P2 规格]
└── i18n/
    ├── home.md               首页文案 zh/en
    ├── proxies.md            代理文案 zh/en
    ├── profiles.md           订阅文案 zh/en
    ├── settings.md           设置文案 zh/en
    ├── shared.md             四页共用文案（shared）
    ├── layout.md             导航与布局文案（四页会引用）
    ├── tests.md              tests 文案（首页网站测试卡 + 解锁测试页）
    ├── connections.md        连接页文案 zh/en
    ├── rules.md              规则页文案 zh/en（不含订阅页规则编辑器）
    ├── logs.md               日志页文案 zh/en
    └── tray.md               托盘文案 zh/en（Rust i18n）
```

## 页面清单与优先级

路由定义见 `src/pages/_navigation-meta.ts`、`src/pages/_navigation.tsx`、`src/pages/_routers.tsx`。侧栏文案带空格（「首 页」），页面标题不带空格（「首页」）。

| 优先级 | 页面 | 路由 | 导航 key | 本轮状态 |
| --- | --- | --- | --- | --- |
| **P0** | 首页 | `/` | `layout.components.navigation.tabs.home` | 完整规格 `pages/home.md` |
| **P0** | 代理 | `/proxies` | `layout.components.navigation.tabs.proxies` | 完整规格 `pages/proxies.md` |
| **P0** | 订阅 | `/profile` | `layout.components.navigation.tabs.profiles` | 完整规格 `pages/profiles.md` |
| **P0** | 设置 | `/settings` | `layout.components.navigation.tabs.settings` | 完整规格 `pages/settings.md` |
| **P0** | 设计令牌 | — | — | 完整 `design-tokens.md` |
| **P0** | 中英文案（上列四页 + shared / layout / 首页测试卡用到的 tests） | — | — | `i18n/` |
| P1 | 整体布局 | 壳层 | — | 大纲 `outlines/layout.md`。窗口默认 940×700、最小 520×520，以及侧栏品牌位（mygo `favicon.svg`，展开 36×36 +「Clash Go」、收起 56×56 仅标记）已记下。细节（拖拽排序、更新按钮、系统服务弹窗）待补 |
| P1 | 连接 | `/connections` | `layout.components.navigation.tabs.connections` | 规格 `outlines/connections.md` + `i18n/connections.md` |
| P1 | 规则 | `/rules` | `layout.components.navigation.tabs.rules` | 规格 `outlines/rules.md` + `i18n/rules.md` |
| P2 | 日志 | `/logs` | `layout.components.navigation.tabs.logs` | 规格 `outlines/logs.md` + `i18n/logs.md` |
| P2 | 测试 | `/unlock` | `layout.components.navigation.tabs.unlock` | 规格 `outlines/unlock.md`（页面标题「解锁测试」）+ `i18n/tests.md`。首页里的网站测试卡在 `pages/home.md` |
| P2 | 托盘菜单 | 系统托盘 | `crates/clash-verge-i18n` 的 `tray.*` | 规格 `outlines/tray.md` + `i18n/tray.md`。v2.5.9 已移除托盘「测试延迟」 |

M2 验收以 P0 四页 + 设计令牌 + 对应文案为准。P1/P2 规格已补到可逐项截图核对的粒度（各文件末尾「截图核对清单」与「待确认」），文件仍放在 `outlines/` 以免改路径；作为 M3 的验收依据。

默认导航顺序（`src/pages/_navigation.tsx`）：首页 → 代理 → 订阅 → 连接 → 规则 → 日志 → 测试 → 设置。用户可在侧栏右键调整顺序，该项属于布局（P1）。

## 阅读约定

1. 每条规格尽量给出 `文件路径:行号` 或组件名。行号用于核对，小版本内可能漂移，以组件名为准。
2. 文案不在规格正文里重复粘贴大表。控件写 i18n key，中英文对照见 `i18n/`。
3. key 的运行时形态：各 `src/locales/{zh,en}/*.json` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`（`src/services/i18n.ts:112-171`）。例如 `home.json` 的 `page.title` 调用为 `t('home.page.title')`。
4. 亮色 / 暗色默认色以 `src/pages/_theme.tsx` 为准。用户可在「主题设置」里覆盖主色、文字色等，背景纸色仍取默认主题的 `background_color`（`use-custom-theme.ts:166-169`）。
5. 全局样式把 `box-shadow` 设成 `none !important`（`use-custom-theme.ts:303-307`），同时 MUI `shadows` 25 档全是 `'none'`（同文件 171 行）。规格里写到的 `elevation` / `boxShadow` 以源码为准；除粘性分组头显式写了 `!important` 阴影外，多数阴影在运行时被全局规则抹掉。截图比对时以「源码声明 + 全局覆盖后的实际观感」两者都记。
6. 平台差异：Linux 默认无系统标题栏（`window.rs:22-25`）；Windows 在非输入框上拦截右键菜单（`_layout.tsx:135-144`）；部分设置项仅 Windows 或非 Linux 显示，规格里单独标出。

## 与 Verge 有意不同

界面和文案仍按 v2.5.9 对齐。下面两处不照搬上游实现，上游在这里的行为是 bug：

1. **日志搜索的「区分大小写」会生效。** Verge `logs.tsx` 先把「时间 + 类型 + payload」转成小写，再交给搜索框。打开区分大小写之后，带大写字母的关键词永远匹配不到。clash-go 在开关打开时按原文比较，大写输入匹配大写文本；关闭时不区分大小写。高亮用同一组开关。
2. **规则「更新全部」会报告失败。** Verge 串行更新时，单个 provider 失败只写控制台，最后仍提示「所有规则集合均已更新」。clash-go 只要有一个更新失败，就用 `更新失败` 报告，不把失败说成全部成功。

## 验收

逐项勾选见 [checklist.md](./checklist.md)。P0 四页每条都对应本目录中的规格小节。
