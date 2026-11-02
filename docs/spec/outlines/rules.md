# 规则（规格，P1）

路由 `/rules`。参照 v2.5.9 `160bb4fe`。页面 `src/pages/rules.tsx` 的 `RulesPage`。组件 `src/components/rule/rule-item.tsx`、`src/components/rule/provider-button.tsx`。回到顶部 `src/components/layout/scroll-top-button.tsx`。文案对照见 `i18n/rules.md`。

> 本版修正旧大纲：无规则或搜索无结果时显示 `BaseEmpty`（收件箱 +「空空如也」），不是空白（`rules.tsx:85-99`）。点击规则行没有详情、没有右键菜单。

## 1. 结构

```text
BasePage full  标题 rules.page.title
│  contentStyle: height 100%、flex 列、overflow auto
├── header（flex，gap 8px）
│   └── ProviderButton（仅存在规则集合时渲染）
├── 搜索行（pt 8px，mb 4px，mx 10px，高 36，flex 居中）
│   └── BaseSearchBox（只匹配 payload）
├── 有结果 → VirtualList（estimateSize 40，overscan 5，flex 1）→ RuleItem
│            + ScrollTopButton（滚过 100px 出现）
└── 无结果 → BaseEmpty
```

出处：`rules.tsx:56-100`。

## 2. 页头：规则集合按钮

`ProviderButton`（`provider-button.tsx`）。

- **没有任何 rule provider 时整个按钮不渲染**（`:132`），页头只有标题。
- 按钮：`variant=outlined size=small`，`startIcon StorageOutlined`，文案 `rules.page.provider.trigger`（「规则集合」）。

### 2.1 规则集合对话框

MUI `Dialog maxWidth="sm" fullWidth`（`:145`）。

| 部位 | 内容 |
| --- | --- |
| 标题行 | 左 `Typography h6`「规则集合」`rules.page.provider.dialogTitle`；右 contained small「更新全部」`rules.page.provider.actions.updateAll`；两端对齐 |
| 列表 | `List py 0`，最小高 250。按 provider 名（对象键）字典序排序 |
| 底部 | outlined「关闭」`shared.actions.close` |

每个列表项（`:176-281`）：

| 项 | 值 |
| --- | --- |
| 容器 | padding 0，下边距 8px，圆角 8px，overflow hidden，过渡 0.2s |
| 底色 | 亮 `#ffffff` / 暗 `#24252f`；hover 亮 primary@10% / 暗 primary@20% |
| 第一行左 | provider 名（subtitle1，单行，title 为全名）+ 8px + 规则数角标 |
| 第一行右 | body2 次要色：`<small>更新于: </small>` + `dayjs(updatedAt).fromNow()`（如「3 小时前」）。key `shared.labels.updateAt` |
| 第二行 | 两个角标：`vehicleType`（如 `HTTP`/`File`）、`behavior`（如 `domain`/`ipcidr`/`classical`） |
| 角标 TypeBox | 10px、行高 1.25、padding `0 2px`、右边距 4、圆角 4、边框 secondary@50%、字色 secondary@80% |
| 右侧操作区 | 竖分隔线 + 40px 宽居中 `IconButton small primary`，图标 `RefreshRounded`，title / aria-label「更新」`rules.page.provider.actions.update` |
| 更新中 | 该按钮禁用并 1s 线性无限旋转 |

### 2.2 更新行为

- 单个：`updateRuleProvider(name)` → 刷新规则与规则集合 → 后台 `syncRuntimeProviders`。成功 `showNotice.success('rules.feedback.notifications.provider.updateSuccess', {name})`；失败 `...updateFailed`（`{{name}}`、`{{message}}`）（`:50-80`）。
- 全部：没有 provider 时 info `...none`。否则全部置为更新中，**串行**逐个更新，每完成一个解除该项旋转；单项失败只记控制台、不弹提示；最后刷新并 success `...allUpdated`。整体异常 error `...genericError`（`:82-126`）。

## 3. 搜索

`BaseSearchBox`（区分大小写 / 全字 / 正则三开关，见 `pages/proxies.md` §7）。只对 `payload` 匹配（`payload ?? ''`），**不匹配类型和策略名**（`rules.tsx:38-46`）。

行号在过滤前计算（`lineNo = 原始 index + 1`），所以过滤后行号保持原规则序号、可能不连续。

## 4. 规则行 `RuleItem`

`rule-item.tsx`。

```text
┌────┬──────────────────────────────────────┐
│ 12 │ google.com                (subtitle1)│
│    │ DomainSuffix        Proxy             │
└────┴──────────────────────────────────────┘
```

| 部位 | 样式 |
| --- | --- |
| 行容器 | flex，padding `4px 16px`，底边 1px `--divider-color`，文字 `text.primary` |
| 序号 | body2（14px），`text.secondary`，行高 2，最小宽 30，右边距 18px，居中 |
| payload | subtitle1（16px），`text.primary`；空值显示 `-`；可选中 |
| 类型 | body2 `text.secondary`，inline-block，最小宽 120，右边距 24px。值是内核原样字符串（`Unknown` 变体取其内值） |
| 策略 | body2，颜色见下 |

策略颜色 `parseColor`（`rule-item.tsx:22-31`）：

| 策略名 | 颜色 |
| --- | --- |
| `REJECT`、`REJECT-DROP` | `error.main` |
| `DIRECT` | `text.primary` |
| 其他 | 字符码之和 `% 5` 取 `[primary, secondary, info.main, warning.main, success.main]` |

例：`Proxy` 字符码和 = 80+114+111+120+121 = 546，546 % 5 = 1 → secondary。clash-go 必须用同一算法，保证同名策略颜色一致。

无点击、无 hover 效果、无右键菜单。

## 5. 回到顶部

`ScrollTopButton`：`Fade` 显隐，绝对定位 `bottom 20 / right 20`，`KeyboardArrowUp` 图标。底色亮 `rgba(0,0,0,0.1)` / 暗 `rgba(255,255,255,0.1)`，hover 加深到 0.2。列表 scrollTop > 100 时出现，点击平滑滚到顶（`rules.tsx:48-54`）。

## 6. 刷新

- 进入页面时刷新规则和规则集合；窗口重新可见时再刷新一次（`rules.tsx:27-36`）。
- **没有轮询或 WebSocket**，停留在页面时列表不变。规则查询 `staleTime 1500ms`、失败重试 3 次（`app-data-provider.tsx:24-30、78-89`）。
- 规则集合查询 `revalidateOnMount: false`，依赖本页主动刷新。

## 7. 亮 / 暗差异

| 项 | 亮 | 暗 |
| --- | --- | --- |
| 内容区 | `#ECECEC` | `#1e1f27` |
| 集合列表项底 | `#ffffff` | `#24252f` |
| 集合 hover | primary@10% | primary@20% |
| 回到顶部底 | 黑@10% | 白@10% |
| 策略颜色 | 取各自主题 palette（亮暗主色值不同，见 design-tokens §1.1） | 同左 |

## 8. 截图核对清单（`ui.NewTester`）

建议数据：8 条规则，策略含 `DIRECT`、`REJECT`、`Proxy`、一个自定义组名；两条规则集合（一 HTTP、一 File）。

- [ ] `rules-light/dark`：标题「规则」；有集合时右上描边按钮「规则集合」带存储图标；无集合时页头只有标题
- [ ] 搜索行高 36，左右缩进 10px，搜索框占满
- [ ] 行：序号灰色居中；第一行 payload 16px；第二行类型灰色、最小宽 120；策略颜色按 §4 表
- [ ] `REJECT` 红、`DIRECT` 正文色、`Proxy` 为 secondary 色
- [ ] 搜索 `google` 后行号保持原值（不从 1 重排）
- [ ] 无规则 / 无结果：收件箱图标 +「空空如也」
- [ ] 滚动 >100px 右下出现圆形上箭头
- [ ] 集合对话框：标题左「规则集合」右「更新全部」；每项名称 + 规则数角标、「更新于: …前」、两个类型角标、右侧刷新图标；底部「关闭」
- [ ] 点击单项刷新：图标旋转且禁用，完成后右上/底部 notice「X 更新成功」
- [ ] 英文：`Rules`、`Rule Provider`（按钮）/ `Rule Providers`（对话框标题）、`Update All`、`Update At: 3 hours ago`

## 9. 待确认

1. 「更新全部」成功文案在中途有单项失败时仍会弹「所有规则集合均已更新」（源码如此，原版 bug）。**与 Verge 有意不同**：clash-go 只要有一项失败就提示失败并写明是哪一项，全部成功才弹「所有规则集合均已更新」（主程 2026-11-03 拍板，附测试）。
2. 规则集合对话框 `DialogContent` 的内边距与对话框圆角走 MUI 默认 + 全局 `.MuiPaper-root` 边框（`--window-border-color`），像素值需对原版截图核对。
3. `rules.modals.editor.*`（编辑规则对话框）属于订阅页的规则编辑器，不在本页。本页规格不覆盖，归入订阅页后续补充。
