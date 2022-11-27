all: app

app: *.go go.mod go.sum
	go build -o app

refresh-app:
	git pull origin master
	make app
	sudo systemctl stop isu-go && sudo systemctl disable isu-go && sudo systemctl start isu-go && sudo systemctl enable isu-go
	make rotate

rotate:
	sh ../../../scripts/nginx-log-rotate.sh
	sh ../../../scripts/slow-query-log-rotate.sh

analyze:
	make analyze-nginx-log
	make analyze-slow-query

analyze-slow-query:
	sudo mysqldumpslow -s t /var/log/mysql/mysql-slow.log

analyze-nginx-log:
	sudo cat /var/log/nginx/access.log | alp json --sort sum -r -m "/posts/[0-9]+,/@\w+,/image/\d+" -o count,method,uri,min,avg,max,sum

reload-nginx:
	sudo nginx -t && sudo systemctl reload nginx

reload-mysql:
	sudo systemctl reload mysql