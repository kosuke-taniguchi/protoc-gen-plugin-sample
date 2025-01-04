protoc:
	protoc --proto_path=proto --go_out=gen/go proto/user.proto
protoc-plugin:
	protoc --proto_path=proto --go_out=gen/go --myplugin_out=gen/go/myplugin proto/user.proto
install:
	go install ./cmd/protoc-gen-myplugin
