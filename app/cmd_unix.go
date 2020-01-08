package app

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"

	"github.com/gookit/color"

	"github.com/Lemo-yxk/go-watch/vars"
)

func (w *Watch) StopProcess() {

	for _, cmd := range w.commands {

		var err error
		switch runtime.GOOS {
		case "windows":
			err = killWindows(cmd)
		default:
			err = killUnix(cmd)
		}
		if err != nil {
			color.Red.Println(err)
		}

		_, err = cmd.Process.Wait()
		if err != nil {
			color.Red.Println(err)
		}

		color.Bold.Println(cmd.Process.Pid, "kill success")
	}

	w.commands = nil
}

func killUnix(cmd *exec.Cmd) error {
	return syscall.Kill(-cmd.Process.Pid, syscall.Signal(vars.Sig))
}

func killWindows(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}

func (w *Watch) startProcess() {

	for _, v := range w.config.start {

		var cmd *exec.Cmd

		if runtime.GOOS != "windows" {
			cmd = exec.Command("bash", "-c", v)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		} else {
			cmd = exec.Command("cmd", "/C", v)
		}

		cmd.Dir = w.listenPath
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		w.commands = append(w.commands, cmd)

		color.Bold.Println(cmd.Process.Pid, "run success")
	}

}