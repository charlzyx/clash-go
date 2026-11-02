# Clash Go

用 Go 和 [mygo](https://github.com/egoist/mygo) 的原生 UI（`ui` 包，不经过 webview）内嵌 [mihomo](https://github.com/metacubex/mihomo) v1.19.33，按 [Clash Verge Rev](https://github.com/clash-verge-rev/clash-verge-rev) 的界面和功能来做。当前只做 macOS Apple Silicon（arm64）。

这一版是里程碑 M1：应用能启动内核、打开带左侧导航的主窗口和菜单栏托盘，GitHub Actions 在 `macos-14` 上跑测试并打包 `.app` / `.dmg`。

## 环境

- Go 1.27.3（`go.mod` 里的版本，与 mygo v0.7.4 一致）
- 打包用模块里钉住的 CLI：`go tool mygo`
- mygo 本身不需要 cgo。mihomo 以库的方式链进来，默认构建也不开 cgo

## 构建

```sh
go test ./...
go run ./cmd/clash-go
go tool mygo dev
go tool mygo build -platform darwin/arm64
```

`go tool mygo build` 在 macOS 上写出：

- `build/darwin-arm64/Clash Go.app`（`Info.plist`、`resources/icon.png` 生成的图标、ad-hoc 签名）
- `build/darwin-arm64/clash-go-0.1.0-arm64.dmg`

签名身份在 `mygo.json` 的 `macos.signingIdentity`，目前是 `-`（ad-hoc）。别的 Mac 上 Gatekeeper 会拦住这个包，正式分发要换成 Developer ID。

窗口内容区是 940×700，最小 520×520，和 Clash Verge Rev 的默认窗口一致。`UseContentSize` 让这两个数字指内容区，截图也是这个尺寸；系统标题栏在内容区外面。

## 目录

```
cmd/clash-go/          入口
internal/core/         mihomo 单例的 Start / Reload / Stop，配置目录由调用方传入
internal/appdir/       配置目录（接口；macOS 用 ~/Library/Application Support/app.clashgo.desktop）
internal/appconfig/    应用自己的 settings.json（系统代理、TUN 开关）
internal/sysproxy/     系统代理接口。M1 只记在内存里，不改系统
internal/tun/          TUN 开关接口。M1 只记在内存里，不打开设备
internal/ui/           主窗口。页面是纯函数，读传入的状态，不碰内核单例
internal/app/          把上面几层接起来：托盘、菜单、退出时停内核
resources/icon.png     1024×1024 占位图标
mygo.json              应用名、bundle id、版本
```

配置目录第一次启动时如果没有 `config.yaml`，会写入一份最小配置：`mixed-port: 7897`，另有 socks `7898`、http `7899`、`external-controller: 127.0.0.1:9097`，规则 `MATCH,DIRECT`。这些端口和 Verge Rev 的默认值一致。首页显示内核版本和是否在跑。

`hub.Parse` 收到的是这份配置去掉 `external-controller` 行之后的副本。mihomo v1.19.33 会在设置日志级别之前就把外部控制器丢进另一个 goroutine，`go test -race` 会把这判成数据竞争，所以控制器在 `hub.Parse` 返回之后才监听。带 `-race` 的测试构建不再打开这个端口（停它同样没有同步），普通构建和打包出来的 app 会听 `127.0.0.1:9097`。

托盘菜单有「显示主窗口」、系统代理、TUN 模式和退出。两个开关会走对应接口并写进 `settings.json`，这一步还不会改系统代理，也不会打开 TUN。

## 测试

```sh
go test -race -count=1 -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | tail -n 1
```

- 会真的打开窗口、改系统代理或打开 TUN 的测试带 `//go:build integration`，默认不编进去。
- 内核测试用 `httptest` 提供一份只含 `DIRECT` 的订阅，不访问外网。mihomo 是进程内单例，这些测试不能 `t.Parallel`。
- 覆盖率暂不设门槛。
- 设置 `CLASHGO_SCREENSHOT_DIR` 后，`internal/ui` 的测试会用 `ui.NewTester` 把每一页画成 940×700 的 PNG，亮色和暗色各一张，文件名像 `home-light.png`、`home-dark.png`、`proxies-light.png`。画的是主窗口那份视图，没有系统标题栏。CI 会把它们当成 artifact 上传。

## CI

- `.github/workflows/ci.yml`：push / PR，`macos-14`，`go vet`、`go build`、带 race 和覆盖率的 `go test`，上传 `coverage.out` 和窗口截图。
- `.github/workflows/release.yml`：打 `v*` tag，或手动触发。在 `macos-14` 上构建 arm64 `.app` 和 `.dmg` 并上传 artifact；tag 与 `mygo.json` 的版本一致时发 GitHub Release。

## 路线

- **M1（当前）** 骨架、内嵌内核、占位页面、托盘占位、macOS arm64 的 `.app` / `.dmg`。
- **M2** 内核配置重载、订阅、代理页、真正的系统代理、托盘状态。
- **M3** 连接、规则、日志、测试、设置页。
- **M4** TUN、服务助手、自动更新。
