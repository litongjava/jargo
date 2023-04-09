package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 获取命令行参数
	startJarPath := flag.String("jar", "", "java jar path")
	stopJarPath := flag.String("stop", "", "java jar path")
	flag.Parse()
	if *startJarPath != "" {
		log.Println(*startJarPath)
		startJavaJar(startJarPath)
		os.Exit(0)
	}

	if *stopJarPath != "" {
		log.Println(*stopJarPath)
		stopJavaJar(stopJarPath)
		os.Exit(0)
	}

}

func startJavaJar(jarPath *string) {
	// 启动Java程序
	//cmd := exec.Command("nohup", "java", "-jar", jarPath, ">>", jarName+".log", "2>&1", "&")
	cmd := exec.Command("nohup", "java", "-jar", *jarPath, "&")
	fmt.Println("java -jar ", *jarPath)
	err := cmd.Start()
	if err != nil {
		log.Print("start java program failed:", err)
		os.Exit(1)
	}
	pid := fmt.Sprintf("%d", cmd.Process.Pid)
	log.Print("java program start finished，pid:", pid)

	pidFile := getPidFile(jarPath)
	log.Println("pid file:", pidFile)
	// 切换到jar包所在目录
	jarDir := filepath.Dir(*jarPath)
	err = os.Chdir(jarDir)
	if err != nil {
		fmt.Println("swithc path failed:", err)
		os.Exit(1)
	} else {
		log.Print("switch path ", jarDir)
	}

	// 将程ID并写入pid文件
	err = os.WriteFile(pidFile, []byte(pid), 0644)
	if err != nil {
		log.Print("write pid file failed", err)
		os.Exit(1)
	}

	tick := time.Tick(1 * time.Second)
	for range tick {
		_, err := os.FindProcess(cmd.Process.Pid)
		if err != nil {
			fmt.Printf("pid %d does not exist\n", pid)
		} else {
			log.Println("pid is running")
			break
		}
	}
	log.Println("finished")
}

func getPidFile(jarPath *string) string {
	// 获取jarName
	jarName := filepath.Base(*jarPath)
	jarName = strings.TrimSuffix(jarName, filepath.Ext(jarName))

	pidFile := jarName + ".pid"
	return pidFile
}

func stopJavaJar(jarPath *string) {
	pidFile := getPidFile(jarPath)
	log.Println("pid file:", pidFile)
	// 切换到jar包所在目录
	jarDir := filepath.Dir(*jarPath)
	err := os.Chdir(jarDir)
	if err != nil {
		fmt.Println("swithc path failed:", err)
		os.Exit(1)
	} else {
		log.Print("switch path ", jarDir)
	}

	pid, err := ioutil.ReadFile(pidFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	pidInt, err := strconv.Atoi(string(pid))
	if err != nil {
		log.Fatalln(err.Error())
	}

	if _, err := os.Stat(fmt.Sprintf("/proc/%d", pidInt)); os.IsNotExist(err) {
		// pidfile exists, but process is not running
		log.Println("pidfile exists, but process is not running")
		err := os.Remove(pidFile)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	//关闭进程
	process, err := os.FindProcess(pidInt)
	if err != nil {
		log.Fatalf("pid %d does not exist\n\n", pid)
	}
	err = process.Kill()
	if err != nil {
		log.Fatalf("pid %d does not exist\n\n", pid)
	} else {
		log.Println("stopped ", pidInt)
		err := os.Remove(pidFile)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

}
