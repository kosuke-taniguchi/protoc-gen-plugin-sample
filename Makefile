protoc:
	protoc --proto_path=proto --go_out=gen/go proto/user.proto
protoc-plugin:
	protoc --proto_path=proto --go_out=gen/go --myplugin_out=gen/go/mysql proto/user.proto
install:
	go install ./plugin/protoc-gen-myplugin
