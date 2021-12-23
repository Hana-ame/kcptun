$env:GO111MODULE="on"

go build -o s -mod vendor ./server
go build -o c -mod vendor ./client

./s -t "144.34.219.146:26275" -l "sbb:4000"
./c -r "sbb" -l ":8388"


./client.exe  -r "sbb" -l ":8388"

./server.exe -t "144.34.19.146" -l "sbb:4000"
./server.exe -t "localhost:80" -l "sbb:4000"


raw链接可通

stun可通，未确认空包是否接受

目前会停止发送一切信息，原因不明

从raw一步一步改过来

~~现在在除了~~