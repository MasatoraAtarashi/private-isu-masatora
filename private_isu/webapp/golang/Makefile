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
