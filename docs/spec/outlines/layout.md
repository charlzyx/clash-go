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
    │   ├── logo
    │   ├── the-menu 导航
    │   └── the-traffic 60px（收起时隐藏）
    └── layout-content__right
        ├── the-bar 36px（占位，内容 top 盖住）
        └── the-content → Outlet
```

侧栏组件 `layout-sidebar.tsx`、项 `layout-item.tsx`、流量 `layout-traffic.tsx`。选中态、图标双套、收起尺寸见 `design-tokens.md`。

右键菜单（`layout-sidebar.tsx:65-96`）：解锁/锁定排序、恢复默认、收起/展开。文案 `layout.components.navigation.menu.*`。排序结果写入 `verge.menu_order`。

流量区：上传 `ArrowUpward`、下载 `ArrowDownward`、可选内存 `MemoryRounded`。`traffic_graph` 为真时叠一张小图（`traffic-graph.tsx`）。`pause_render_traffic_stats_on_blur` 默认开。

全局浮层（不属于某一页，P1 要补交互）：

- `NoticeManager`：toast 位置四角，默认右上。
- `ServiceMigrationDialog`：服务版本不符 / 路径被拒 / 未自动启动。
- `SysproxyPrivilegeDialog`：macOS 系统代理或 TUN 缺权限。
- `UpdateButton`：侧栏 logo 上的新版本角标。
- 主题未就绪：整屏亮 `#fff` / 暗 `#181a1b`。

待补：流量图采样与颜色、更新按钮状态、两个系统对话框的逐步状态、macOS 交通灯与拖拽区的精确对齐。
