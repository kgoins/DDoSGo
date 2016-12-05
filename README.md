# DDoSGo
Senior Project Fall 2016 for Kyle Goins, R. Alex Clark, and Nick Werner

# Installation
clone repo
set GOPATH

go get github.com/google/gopacket/layers

go build DDoSGo/src/agent.go
go build DDoSGo/src/handler.go

# Running Agent/Handler binaries
sudo ./agent
./handler
