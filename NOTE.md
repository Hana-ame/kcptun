$env:GO111MODULE="on"

go build -o s -mod vendor ./server
go build -o c -mod vendor ./client

go build  -mod vendor ./server
go build  -mod vendor ./client

./s -t "144.34.219.146:26275" -l "sbb:4000"
./c -r "sbb" -l ":8388"


./client.exe  -r "sbb" -l ":8388"

./server.exe -t "144.34.19.146" -l "sbb:4000"
./server.exe -t "localhost:80" -l "sbb:4000"


@ override 链接，没啥毛病
./server.exe -t "localhost:80" -l "sbb:4000"
./client.exe  -r "sbb" -l ":8080"

@!  server的 listen addr 单纯的overridde不可用，应该是问题原因

@?  server监听4000是OK的，但是40000不行，什么玩意啊……


# 转成udp转发器的思路吧，这个kcptun实在是太难懂了
最多也就是把端口记录一下加到头里面吧，大概


@@@ 
人傻了，不看了

raw链接可通

stun可通，未确认空包是否接受

目前会停止发送一切信息，原因不明

从raw一步一步改过来

~~现在在除了~~

----

./client.exe  -r "localhost:6000" -l ":8080"
./server.exe -t "hana-sweet.top:443" -l ":4000"

./udptun.exe -mode "client" -l ":6000" -r "127.0.0.1:10000"
./udptun.exe -mode "server" -l ":10000" -r "127.0.0.1:4000"



go build -o s -mod vendor ./server &
go build -o c -mod vendor ./client &
go build -o u -mod vendor ./udptun &

./s -t "localhost:8080" -l ":4000"
./u -mode "server" -l "gcp" -r "127.0.0.1:4000"

./udptun.exe -mode "client" -l ":6000" -r "gcp"
./client.exe  -r "localhost:6000" -l ":8080"

it works well