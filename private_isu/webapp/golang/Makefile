all: app

app: *.go go.mod go.sum
	go build -o app

refresh:
	make refresh-app
	make rotate

refresh-app:
	git pull origin master
	make app
	sudo systemctl stop isu-go && sudo systemctl disable isu-go && sudo systemctl start isu-go && sudo systemctl enable isu-go

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

delete-old-logs:
	sudo ls /var/log/mysql/ | grep mysql-slow. | sudo xargs -If rm /var/log/mysql/f
	sudo ls /var/log/nginx/ | grep access.log. | sudo xargs -If rm /var/log/nginx/f

delete-images:
	rm /home/isucon/private_isu/webapp/public/image/*