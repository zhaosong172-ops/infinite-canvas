# Windows 便携版构建说明

产物：`InfiniteCanvas-Windows-x64-<tag>.zip`，内含 `InfiniteCanvas.exe`、`server.exe`、Next.js standalone 前端、内置 `node.exe` 和 ComfyUI Bridge 三个平台的二进制。

## 本地构建

前置：Windows x64、Go ≥ 1.25、Node.js 22、npm。

```powershell
# 1. ComfyUI Bridge（三个平台）
cd canvas-agent/native/comfy-bridge
$env:CGO_ENABLED = "0"
$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -trimpath -ldflags="-s -w" -o ../../../web/public/InfiniteCanvas-ComfyBridge.exe .
$env:GOOS = "linux";  $env:GOARCH = "amd64"; go build -trimpath -ldflags="-s -w" -o ../../../web/public/InfiniteCanvas-ComfyBridge-linux-amd64 .
$env:GOOS = "linux";  $env:GOARCH = "arm64"; go build -trimpath -ldflags="-s -w" -o ../../../web/public/InfiniteCanvas-ComfyBridge-linux-arm64 .

# 2. 前端（standalone）
cd ../../../web
npm install
$env:CANVAS_PREBUILT_BRIDGE = "1"
npx next build

# 3. 后端与启动器
cd ..
go build -trimpath -ldflags="-s -w" -o dist/InfiniteCanvas/server.exe .
cd packaging/windows-launcher
$env:CGO_ENABLED = "0"
go build -trimpath -ldflags="-s -w" -o ../../dist/InfiniteCanvas/InfiniteCanvas.exe .

# 4. 组装
$target = "dist/InfiniteCanvas"
New-Item -ItemType Directory -Force "$target/web" | Out-Null
Copy-Item "web/.next/standalone/*" "$target/web" -Recurse -Force
New-Item -ItemType Directory -Force "$target/web/.next" | Out-Null
Copy-Item "web/.next/static" "$target/web/.next/static" -Recurse -Force
Copy-Item "web/public/*" "$target/web/public" -Recurse -Force
Copy-Item "packaging/windows/env.example" "$target/.env"
Copy-Item "packaging/windows/README.txt" "$target/README.txt"

# 5. 内置 Node 运行时
Invoke-WebRequest "https://nodejs.org/dist/v22.20.0/node-v22.20.0-win-x64.zip" -OutFile node.zip
Expand-Archive node.zip -DestinationPath node-runtime
New-Item -ItemType Directory -Force "$target/node" | Out-Null
Copy-Item "node-runtime/node-v22.20.0-win-x64/node.exe" "$target/node/node.exe"
```

## 冒烟测试

```powershell
$target = (Resolve-Path "dist/InfiniteCanvas").Path
$env:PORT = "8080"; $api = Start-Process "$target/server.exe" -WorkingDirectory $target -PassThru -NoNewWindow
$env:PORT = "3000"; $env:API_BASE_URL = "http://127.0.0.1:8080"; $env:HOSTNAME = "127.0.0.1"
$web = Start-Process "$target/node/node.exe" -ArgumentList "web/server.js" -WorkingDirectory $target -PassThru -NoNewWindow
Invoke-RestMethod "http://127.0.0.1:3000/api/settings"
Stop-Process -Id $web.Id, $api.Id -Force
```

## 打包与发布

```powershell
7z a -tzip -mx=5 "InfiniteCanvas-Windows-x64-<tag>.zip" "dist/InfiniteCanvas"
(Get-FileHash "InfiniteCanvas-Windows-x64-<tag>.zip" -Algorithm SHA256).Hash.ToLower()
```

## 用 GitHub Actions 自动发布（可选）

`release-workflow.yml.template` 是等价的工作流，复制到 `.github/workflows/windows-release.yml` 后，推送 `v*` 标签即可自动构建并发布 Release。

注意：通过 API 写入 `.github/workflows/` 需要令牌具备 `workflow` scope，否则 GitHub 会直接返回 404。
