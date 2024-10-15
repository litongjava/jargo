package main

import (
  "fmt"
  "log"
  "os"
  "os/exec"
  "path/filepath"
  "strconv"
  "strings"
)

// Entry point of the program
func main() {
  args := os.Args[1:]

  if len(args) == 0 {
    printUsage()
    os.Exit(1)
  }

  switch args[0] {
  case "--stop":
    if len(args) < 2 {
      fmt.Println("Usage: jar-boot --stop [jar path]")
      os.Exit(1)
    }
    jarPath := args[1]
    stopJavaJar(jarPath)
  case "--fork":
    if len(args) < 2 {
      fmt.Println("Usage: jar-boot --fork [jar path] [jar args...]")
      os.Exit(1)
    }
    jarPath := args[1]
    jarArgs := args[2:]
    startJavaJar(jarPath, jarArgs, true)
  default:
    jarPath := args[0]
    jarArgs := args[1:]
    startJavaJar(jarPath, jarArgs, false)
  }
}

// Print usage instructions
func printUsage() {
  fmt.Println("Usage:")
  fmt.Println("  jar-boot [jar path] [jar args...]           Start jar in foreground")
  fmt.Println("  jar-boot --fork [jar path] [jar args...]    Start jar in background")
  fmt.Println("  jar-boot --stop [jar path]                  Stop jar")
}

// Start the Java JAR either in foreground or background
func startJavaJar(jarPath string, jarArgs []string, fork bool) {
  absJarPath, err := filepath.Abs(jarPath)
  if err != nil {
    log.Fatalf("Failed to get absolute jar path: %v", err)
  }

  jarDir := filepath.Dir(absJarPath)
  jarName := strings.TrimSuffix(filepath.Base(absJarPath), filepath.Ext(absJarPath))
  pidFile := filepath.Join(jarDir, jarName+".pid")

  if fork {
    // Check if the process is already running
    if _, err := os.Stat(pidFile); err == nil {
      log.Fatalf("PID file %s exists. Is the process already running?", pidFile)
    }

    // Prepare the command with arguments
    cmd := exec.Command("java", append([]string{"-jar", absJarPath}, jarArgs...)...)
    cmd.Dir = jarDir

    // Redirect output to a log file
    logFile, err := os.OpenFile(filepath.Join(jarDir, jarName+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
      log.Fatalf("Failed to open log file: %v", err)
    }
    defer logFile.Close()
    cmd.Stdout = logFile
    cmd.Stderr = logFile

    // Start the process
    err = cmd.Start()
    if err != nil {
      log.Fatalf("Failed to start Java jar: %v", err)
    }

    pid := cmd.Process.Pid
    log.Printf("Java jar started in background, PID: %d", pid)

    // Write the PID to the pid file
    err = os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
    if err != nil {
      log.Fatalf("Failed to write PID file: %v", err)
    }
  } else {
    // Foreground execution

    // Prepare the command with arguments
    cmd := exec.Command("java", append([]string{"-jar", absJarPath}, jarArgs...)...)
    cmd.Dir = jarDir

    // Redirect stdout and stderr to the current terminal
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Start the process
    err = cmd.Start()
    if err != nil {
      log.Fatalf("Failed to start Java jar: %v", err)
    }

    pid := cmd.Process.Pid
    log.Printf("Java jar started in foreground, PID: %d", pid)

    // Wait for the process to finish
    err = cmd.Wait()
    if err != nil {
      log.Printf("Java jar exited with error: %v", err)
    } else {
      log.Println("Java jar exited successfully")
    }
  }
}

// Stop the running Java JAR
func stopJavaJar(jarPath string) {
  absJarPath, err := filepath.Abs(jarPath)
  if err != nil {
    log.Fatalf("Failed to get absolute jar path: %v", err)
  }

  pidFile := getPidFile(absJarPath)

  // Read the PID from the pid file
  pidBytes, err := os.ReadFile(pidFile)
  if err != nil {
    log.Fatalf("Failed to read PID file: %v", err)
  }

  pidStr := strings.TrimSpace(string(pidBytes))
  pid, err := strconv.Atoi(pidStr)
  if err != nil {
    log.Fatalf("Invalid PID in PID file: %v", err)
  }

  // Find the process
  process, err := os.FindProcess(pid)
  if err != nil {
    log.Fatalf("Process with PID %d not found: %v", pid, err)
  }

  // Attempt to gracefully terminate the process
  err = process.Kill()
  if err != nil {
    log.Fatalf("Failed to kill process with PID %d: %v", pid, err)
  }

  log.Printf("Process with PID %d stopped", pid)

  // Remove the pid file
  err = os.Remove(pidFile)
  if err != nil {
    log.Printf("Failed to remove PID file: %v", err)
  }
}

// Helper function to get the PID file path based on the JAR path
func getPidFile(jarPath string) string {
  jarDir := filepath.Dir(jarPath)
  jarName := strings.TrimSuffix(filepath.Base(jarPath), filepath.Ext(jarPath))
  return filepath.Join(jarDir, jarName+".pid")
}
