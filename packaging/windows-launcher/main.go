// Infinite Canvas Windows 启动器：拉起 Go 后端与 Next.js 前端，就绪后打开浏览器。
package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	apiPort   = "8080"
	webPort   = "3000"
	webOrigin = "http://127.0.0.1:3000"
)

func main() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("无法定位程序目录：", err)
		pause()
		return
	}
	root := filepath.Dir(exe)
	if err := os.Chdir(root); err != nil {
		fmt.Println("切换工作目录失败：", err)
		pause()
		return
	}
	_ = os.MkdirAll(filepath.Join(root, "data"), 0o755)

	if portOpen(webPort, 800*time.Millisecond) {
		fmt.Println("检测到 3000 端口已有服务，直接打开浏览器。")
		openBrowser()
		return
	}

	fmt.Println("正在启动 Infinite Canvas…")
	api := exec.Command(filepath.Join(root, "server.exe"))
	api.Dir = root
	api.Env = append(os.Environ(), "PORT="+apiPort)
	api.Stdout, api.Stderr = os.Stdout, os.Stderr
	if err := api.Start(); err != nil {
		fmt.Println("后端启动失败：", err)
		pause()
		return
	}
	defer stop(api)

	node := exec.Command(filepath.Join(root, "node", "node.exe"), filepath.Join("web", "server.js"))
	node.Dir = root
	node.Env = append(os.Environ(),
		"NODE_ENV=production",
		"PORT="+webPort,
		"HOSTNAME=127.0.0.1",
		"API_BASE_URL=http://127.0.0.1:"+apiPort,
	)
	node.Stdout, node.Stderr = os.Stdout, os.Stderr
	if err := node.Start(); err != nil {
		fmt.Println("前端启动失败：", err)
		pause()
		return
	}
	defer stop(node)

	fmt.Print("等待前端就绪")
	if !waitPort(webPort, 120*time.Second) {
		fmt.Println()
		fmt.Println("120 秒内前端仍未就绪，请查看上面的日志。")
		pause()
		return
	}
	fmt.Println()

	fmt.Println("已启动：", webOrigin)
	fmt.Println("关闭这个窗口即可退出，数据保存在 data 目录。")
	openBrowser()

	exited := make(chan struct{}, 2)
	go func() { _ = api.Wait(); exited <- struct{}{} }()
	go func() { _ = node.Wait(); exited <- struct{}{} }()
	<-exited
	stop(api)
	stop(node)
	fmt.Println("服务已退出。")
}

func portOpen(port string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func waitPort(port string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if portOpen(port, 800*time.Millisecond) {
			return true
		}
		time.Sleep(400 * time.Millisecond)
		fmt.Print(".")
	}
	return false
}

func openBrowser() {
	_ = exec.Command("cmd", "/c", "start", "", webOrigin).Start()
}

func stop(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
}

func pause() {
	fmt.Println("按回车键退出…")
	_, _ = fmt.Scanln()
}
