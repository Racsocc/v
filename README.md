# v

`v` 是一个基于 `asdf` 的交互式命令行工具，用于快速切换并写入当前项目的 `.tool-versions`。

它读取你已经通过 `asdf` 安装好的版本列表，选一个，然后写回文件。仅此而已。
默认在当前目录向上查找 `.tool-versions`（或 `.git`）作为项目根。

`v`只负责选择并写入版本，不负责安装、卸载或管理 `asdf` 插件。


## 用法
```zsh
v
v -h
v -help

v node
v nodejs
v py
v python


v -v
v --version
v version
```

行为：
- 在项目根目录查找 `.tool-versions`
- 从 `asdf list <plugin>` 读取已安装版本
- 高亮当前版本
- 选择后写回 `.tool-versions`
  
## 前置条件

- macOS  
- 已初始化的 `asdf`
- 目标插件至少安装过一个版本

---

## 编译 与 安装

### Go

```bash
go install .
```

确保 `$GOPATH/bin` 在 `PATH` 中。

### 二进制安装

适合非Go用户。

1. Releases 页面，下载适压缩包  
   （下面以 `v-macos-arm64.tar.gz` 为例，实际文件名以 Releases 为准）。
2. 在下载目录解压并移动到一个在 `PATH` 里的目录，例如 `~/bin`：

```bash
mkdir -p "$HOME/bin"
tar -xzf v-macos-arm64.tar.gz
mv v "$HOME/bin/"
chmod +x "$HOME/bin/v"
```

3. 确保该目录在 `PATH` 中（以 zsh 为例）：

```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

4. 验证是否安装成功：

```bash
v -v
v --version
v version
# 均可 
```


如果你已经用 `asdf` 装了一堆版本，还在手动改 `.tool-versions`，  
你可以尝试一下 `v` 。
