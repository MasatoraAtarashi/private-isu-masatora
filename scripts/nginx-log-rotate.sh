#!/bin/sh

# 実行時点の日時を YYYYMMDD-HHMMSS 形式で付与したファイル名にローテートする
sudo mv /var/log/nginx/access.log /var/log/nginx/access.log.`date +%Y%m%d-%H%M%S`

# nginxにログファイルを開き直すシグナルを送信する
sudo nginx -s reopen