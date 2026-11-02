# 整体布局（大纲，P1）

本页只列结构，供后续补全。窗口尺寸已经可以按下面数字实现。参照 v2.5.9 `160bb4fe`。

## 窗口

`src-tauri/src/utils/resolve/window.rs:16-89`。

- 标题 `Clash Verge`。
- 默认内尺寸 **940 × 700**，居中。
- 最小内尺寸 **520 × 520**。恢复的窗口小于该值时重置为默认并居中。
- 初始不可见，等页面 load 完再 show，避免主题色闪一下。
- 装饰：Linux `decorations(false)`，Windows / macOS 默认 true。布局设置里的「使用系统标题栏」可切换。
- 无装饰时画自定义标题栏（高 36px，底边 `--divider-color`）和 8 向缩放热区（边 6px、角 12px）。缩放热区仅 Linux 且未最大化（`window-controller.tsx:46`）。
- 窗口按钮：最小化 `Minimize`、最大化 `CropSquare` / 还原 `FilterNone`、关闭 `Close`。macOS 按钮在左侧。

## 壳

`src/pages/_layout.tsx` + `src/assets/styles/layout.scss`。

```text
Paper.layout
├── WindowResizeHandles（无装饰）
├── the_titlebar（无装饰）
└── layout-content
    ├── layout-content__left 侧栏 200px（收起 72px）
    │   ├── logo  展开 36×36 +「Clash Go」；收起 56×56 仅标记
    │   ├── the-menu 导航
    │   └── the-traffic 60px（收起时隐藏）
    └── layout-content__right
        ├── the-bar 36px（占位，内容 top 盖住）
        └── the-content → Outlet
```

侧栏组件 `layout-sidebar.tsx`、项 `layout-item.tsx`、流量 `layout-traffic.tsx`。选中态、图标双套、收起尺寸见 `design-tokens.md`。品牌位见下。

### 品牌位 / logo

侧栏顶部 `.the-logo`。槽位对齐 Verge `layout-sidebar.tsx:127-150` 与收起规则 `layout.scss:265-313`。Verge 这里是亮/暗两张图标再加 `logo.svg` 词标；clash-go 落地资产不同，令牌对照写在 `design-tokens.md` §8.1。

资源：[`docs/spec/assets/brand/favicon.svg`](../assets/brand/favicon.svg)。favicon 来自 [egoist/mygo](https://github.com/egoist/mygo) `website/public/favicon.svg`（MIT，Copyright (c) 2026 MyGo contributors）。图形是青色（`#00add8`）圆角方块加白环，静态几何标记，无 GIF。mygo 官网 animation 与 README banner GIF 不纳入品牌位。

- **展开**（侧栏 200px）：标记 **36×36**。标记右侧紧挨文案「Clash Go」（直接排字），`margin-right` 约 **5px**。Verge 同一枚图标是 `marginRight: '5px'`，另有 `marginTop` / `marginLeft` 各 `-3px`，所在内层行高 27px（`layout-sidebar.tsx:131-144`）；clash-go 沿用这组偏移。文案颜色用主题主文字色（亮 `#000000` / 暗 `#FFFFFF`），对应 Verge `LogoSvg` 的 `fill`（`layout-sidebar.tsx:147`）。
- **收起**（侧栏 72px）：标记落在 **56×56** 的 logo 区内（`--layout-nav-logo-size: 56px`，`layout.scss:269、287-294`），区内只显示标记并居中。标记仍是同一张 36×36 的 favicon.svg（Verge 收起时词标 `svg:last-child` 为 `display: none`，图标 `margin: 0` 后仍是 36×36，`layout.scss:303-309`）。

右键菜单（`layout-sidebar.tsx:65-96`）：解锁/锁定排序、恢复默认、收起/展开。文案 `layout.components.navigation.menu.*`。排序结果写入 `verge.menu_order`。

流量区：上传 `ArrowUpward`、下载 `ArrowDownward`、可选内存 `MemoryRounded`。`traffic_graph` 为真时叠一张小图（`traffic-graph.tsx`）。`pause_render_traffic_stats_on_blur` 默认开。

全局浮层（不属于某一页，P1 要补交互）：

- `NoticeManager`：toast 位置四角，默认右上。
- `ServiceMigrationDialog`：服务版本不符 / 路径被拒 / 未自动启动。
- `SysproxyPrivilegeDialog`：macOS 系统代理或 TUN 缺权限。
- `UpdateButton`：侧栏 logo 上的新版本角标。
- 主题未就绪：整屏亮 `#fff` / 暗 `#181a1b`。

待补：流量图采样与颜色、更新按钮状态、两个系统对话框的逐步状态、macOS 交通灯与拖拽区的精确对齐。
