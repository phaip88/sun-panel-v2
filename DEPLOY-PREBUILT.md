# Sun-Panel-V2 预编译版本部署指南

## 概述

本指南适用于使用 GitHub Actions 预编译的二进制文件进行部署。无需安装 Go、Node.js 或任何编译工具。

## 系统要求

- Linux (x86_64 或 ARM64)
- Windows (x86_64)
- macOS (Intel 或 Apple Silicon)
- 至少 100MB 磁盘空间
- 3002 端口可用

## 快速部署

### 1. Linux 部署 (x86_64)

```bash
# 下载最新版本
wget https://github.com/phaip88/sun-panel-v2/releases/latest/download/sun-panel-linux-amd64.tar.gz

# 解压
tar -xzf sun-panel-linux-amd64.tar.gz
cd sun-panel-linux-amd64

# 赋予执行权限
chmod +x sun-panel-linux-amd64

# 运行
./sun-panel-linux-amd64

# 或后台运行
nohup ./sun-panel-linux-amd64 > sun-panel.log 2>&1 &
```

### 2. Linux 部署 (ARM64)

```bash
# 下载 ARM64 版本
wget https://github.com/phaip88/sun-panel-v2/releases/latest/download/sun-panel-linux-arm64.tar.gz

# 解压
tar -xzf sun-panel-linux-arm64.tar.gz
cd sun-panel-linux-arm64

# 赋予执行权限
chmod +x sun-panel-linux-arm64

# 运行
./sun-panel-linux-arm64
```

### 3. Windows 部署

```bash
# 下载
# https://github.com/phaip88/sun-panel-v2/releases/latest/download/sun-panel-windows-amd64.zip

# 解压到任意目录

# 双击运行 sun-panel-windows-amd64.exe

# 或在命令行运行
sun-panel-windows-amd64.exe
```

### 4. macOS 部署

```bash
# 下载 Intel 版本
wget https://github.com/phaip88/sun-panel-v2/releases/latest/download/sun-panel-macos-amd64.tar.gz

# 或下载 Apple Silicon 版本
wget https://github.com/phaip88/sun-panel-v2/releases/latest/download/sun-panel-macos-arm64.tar.gz

# 解压
tar -xzf sun-panel-macos-amd64.tar.gz
cd sun-panel-macos-amd64

# 赋予执行权限
chmod +x sun-panel-macos-amd64

# 运行
./sun-panel-macos-amd64
```

## 访问应用

部署完成后，在浏览器中访问：

```
http://localhost:3002
```

**默认账号**：
- 用户名: `admin`
- 密码: `123456`

## 配置管理

### 首次运行

首次运行时，应用会自动在 `conf/` 目录生成配置文件：

```
conf/
├── app.ini          # 应用配置
├── database.ini     # 数据库配置
└── ...
```

### 修改端口

编辑 `conf/app.ini`，修改 `http_port` 值：

```ini
[base]
http_port = 3003
```

### 使用 MySQL 数据库

编辑 `conf/database.ini`：

```ini
[database]
type = mysql
host = localhost
port = 3306
user = root
password = your_password
database = sun_panel
```

## 进程管理

### 后台运行

```bash
# 使用 nohup
nohup ./sun-panel-linux-amd64 > sun-panel.log 2>&1 &

# 或使用 screen
screen -S sun-panel
./sun-panel-linux-amd64
# 按 Ctrl+A 然后 D 分离
```

### 使用 systemd 服务

创建 `/etc/systemd/system/sun-panel.service`：

```ini
[Unit]
Description=Sun-Panel-V2
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/sun-panel
ExecStart=/opt/sun-panel/sun-panel-linux-amd64
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable sun-panel
sudo systemctl start sun-panel
sudo systemctl status sun-panel
```

### 使用 PM2

```bash
# 安装 PM2
npm install -g pm2

# 创建 ecosystem.config.js
cat > ecosystem.config.js << 'EOF'
module.exports = {
  apps: [{
    name: 'sun-panel',
    script: './sun-panel-linux-amd64',
    cwd: '/opt/sun-panel',
    instances: 1,
    exec_mode: 'fork',
    env: {
      NODE_ENV: 'production'
    }
  }]
};
EOF

# 启动
pm2 start ecosystem.config.js

# 保存配置
pm2 save

# 开机自启
pm2 startup
```

## 反向代理配置

### Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3002;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Apache

```apache
<VirtualHost *:80>
    ServerName your-domain.com
    
    ProxyPreserveHost On
    ProxyPass / http://localhost:3002/
    ProxyPassReverse / http://localhost:3002/
</VirtualHost>
```

## 故障排除

### 权限被拒绝

```bash
chmod +x sun-panel-linux-amd64
```

### 端口被占用

```bash
# 查找占用 3002 端口的进程
lsof -i :3002

# 杀掉进程
kill -9 <PID>

# 或修改配置文件中的端口
```

### 无法连接到数据库

- 检查数据库配置文件
- 确保数据库服务正在运行
- 检查数据库用户名和密码

### 前端资源加载失败

- 确保 `web/` 目录存在
- 检查文件权限
- 查看浏览器控制台错误信息

## 数据备份

### SQLite 数据库

```bash
# 备份
cp database/sun-panel.db database/sun-panel.db.backup

# 恢复
cp database/sun-panel.db.backup database/sun-panel.db
```

### 上传文件

```bash
# 备份
tar -czf uploads-backup.tar.gz uploads/

# 恢复
tar -xzf uploads-backup.tar.gz
```

## 更新升级

```bash
# 1. 停止当前服务
kill <PID>

# 2. 下载新版本
wget https://github.com/phaip88/sun-panel-v2/releases/latest/download/sun-panel-linux-amd64.tar.gz

# 3. 解压覆盖
tar -xzf sun-panel-linux-amd64.tar.gz

# 4. 重新启动
./sun-panel-linux-amd64
```

## 性能优化

### 增加文件描述符限制

```bash
# 编辑 /etc/security/limits.conf
echo "* soft nofile 65535" >> /etc/security/limits.conf
echo "* hard nofile 65535" >> /etc/security/limits.conf
```

### 使用 MySQL 替代 SQLite

对于高并发场景，建议使用 MySQL 数据库。

## 常见问题

**Q: 是否需要安装 Go？**  
A: 不需要。预编译版本已包含所有依赖。

**Q: 是否支持 HTTPS？**  
A: 建议使用 Nginx/Apache 反向代理配置 HTTPS。

**Q: 数据存储在哪里？**  
A: 默认使用 SQLite，数据存储在 `database/` 目录。

**Q: 如何修改默认密码？**  
A: 登录后在管理界面修改。

## 支持

如有问题，请提交 Issue：
https://github.com/phaip88/sun-panel-v2/issues
