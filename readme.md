# JarGo

**JarGo** 是一个轻量级的命令行工具，旨在简化 Java JAR 应用程序的管理。它允许您以前台和后台的方式启动 Java 应用程序，向 JAR 文件传递参数，并优雅地停止正在运行的实例。无论您是在本地开发还是部署到生产环境，JarGo 都提供了一种简单高效的方式来处理您的 Java 进程。

## 功能

- **前台执行**：在终端中直接运行 Java 应用程序并实时查看日志。
- **后台执行**：以后台进程的方式启动 Java 应用程序，并具有日志记录功能。
- **参数传递**：在启动时轻松向您的 JAR 文件传递自定义参数。
- **进程管理**：自动处理 PID 文件以启动和停止应用程序。
- **日志记录**：将应用程序输出重定向到日志文件，方便监控。

## 安装

### 前提条件

- **Go**：确保您已安装 Go。您可以从 [golang.org](https://golang.org/dl/) 下载。

### 从源码构建

1. **克隆仓库**

   ```bash
   git clone https://github.com/litongjava/jargo.git
   cd jargo
   ```

2. **构建二进制文件**

   ```bash
   go build -o jargo main.go
   ```

   此命令将编译 Go 源代码并生成一个名为 `jargo` 的可执行文件。

3. **移动到系统的 PATH 目录**

   为了在任何地方使用 `jargo`，将二进制文件移动到系统 `PATH` 中包含的目录，例如 `/usr/local/bin`：

   ```bash
   sudo mv jargo /usr/local/bin/
   ```

### 下载预编译的二进制文件

如果您不想从源码构建，可以从 [Releases](https://github.com/litongjava/jargo/releases) 页面下载预编译的二进制文件。

## 使用方法

JarGo 提供了一个简单的接口来管理您的 Java JAR 应用程序。以下是可用的命令和示例。

### 以前台模式启动 JAR

直接在终端中运行 Java 应用程序，并实时查看日志。

```bash
jargo path/to/your-application.jar
```

**示例：**

```bash
jargo whisper-java-server-1.0.0.jar
```

### 以后台模式启动 JAR

将 Java 应用程序作为后台进程运行。日志将重定向到日志文件，并创建一个 PID 文件以进行进程管理。

```bash
jargo --fork path/to/your-application.jar [jar 参数...]
```

**示例：**

- **基本后台启动：**

  ```bash
  jargo --fork whisper-java-server-1.0.0.jar
  ```

- **带参数的后台启动：**

  ```bash
  jargo --fork whisper-java-server-1.0.0.jar --server.port=8000
  ```

### 停止正在运行的 JAR 应用程序

使用 PID 文件优雅地停止正在运行的 Java 应用程序。

```bash
jargo --stop path/to/your-application.jar
```

**示例：**

```bash
jargo --stop whisper-java-server-1.0.0.jar
```

## 工作原理

- **PID 文件**：在后台运行 JAR 时，JarGo 会在 JAR 所在目录创建一个 `.pid` 文件。该文件包含正在运行应用程序的进程 ID（PID），用于管理进程。
  
- **日志记录**：后台进程的日志将写入与 JAR 相同目录下的 `<jar-name>.log` 文件。这允许您在不杂乱终端的情况下监控应用程序输出。

- **进程终止**：`--stop` 命令读取 `.pid` 文件中的 PID，并发送终止信号以优雅地停止应用程序。然后删除 `.pid` 文件。

## 错误处理

JarGo 包含全面的错误处理，以确保操作顺利进行：

- **PID 文件冲突**：通过检查现有 PID 文件，防止启动同一 JAR 的多个实例。
- **文件缺失**：如果指定的 JAR 或 PID 文件不存在，将发出警报。
- **进程终止**：如果无法终止进程或进程未运行，将通知用户。

## 贡献

欢迎贡献！如果有任何改进、错误修复或建议，请提交问题或拉取请求。

1. **Fork 仓库**
2. **创建功能分支**

   ```bash
   git checkout -b feature/YourFeature
   ```

3. **提交您的更改**

   ```bash
   git commit -m "Add Your Feature"
   ```

4. **推送到分支**

   ```bash
   git push origin feature/YourFeature
   ```

5. **打开一个拉取请求**

## 许可证

本项目使用 [MIT 许可证](LICENSE) 进行许可。

## 支持

如果您遇到任何问题或有疑问，请在 [GitHub 仓库](https://github.com/litongjava/jargo/issues) 上打开一个问题，或通过电子邮件联系维护者：[your.email@example.com](mailto:your.email@example.com)。

## 致谢

- 受各种进程管理工具的启发，旨在简化 Java 应用程序的部署。