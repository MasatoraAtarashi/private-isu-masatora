#!/bin/sh

# 本当はlinuxのlogrotateとかを使いたい？
sudo mv /var/log/mysql/mysql-slow.log /var/log/mysql/mysql-slow.log.`date +%Y%m%d-%H%M%S`

sudo systemctl restart mysql