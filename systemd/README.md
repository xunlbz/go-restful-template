# 1.将edge_admin与EdgeAdmin.service 复制到服务器 /root 目录下
# 2.依次执行命令 
``` bash
    chmod +x edge_admin
    cp edge_admin /usr/bin/
    cp EdgeAdmin.service /usr/lib/systemd/system/
    systemctl daemon-reload
    systemctl enable edgeAdmin
    systemctl start edgeAdmin
```
# 3.查看状态
``` bash
systemctl status edgeAdmin
```
# 4.查看日志：
``` bash
journalctl -f -u edgeAdmin
```