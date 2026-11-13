## Windows 便携版（ComfyUI 多服务端）

基于 `feat/comfy-direct-servers` 分支构建的免安装 Windows x64 包，内置 Node.js 运行时，解压即用。

### 这个版本包含什么

- ComfyUI 多服务端：系统设置里可配置默认地址、API Key 和额外服务端。
- 未绑定 Bridge 的 ComfyUI 渠道由后端直连服务端执行（提交 `/prompt`、轮询 `/history`、回传产物地址）。
- 工作流可绑定具体服务端，删掉服务端后自动回落到默认地址。
- 其余功能与 `main` 分支一致。

### 使用方法

1. 下载并解压压缩包，保持目录结构完整。
2. 双击 `InfiniteCanvas.exe`，启动器会拉起后端（8080）和前端（3000）并打开浏览器。
3. 默认账号 `admin` / `infinite-canvas`，请在 `.env` 中修改密码与 `JWT_SECRET`。
4. 数据、日志保存在 `data` 目录（默认 SQLite）。

### 说明

- 这是便携包，不是安装程序（上游的 `setup.exe` 由仓库外的 Tauri 桌面壳构建，本仓库不含该壳）。
- 未做代码签名，首次运行可能触发 Windows SmartScreen，选择“更多信息 → 仍要运行”。
- 构建后已在 Windows 上启动后端与前端做过冒烟测试（`/api/settings` 与首页均返回 200）。
- 复现构建步骤见仓库 `packaging/windows/BUILD.md`。
