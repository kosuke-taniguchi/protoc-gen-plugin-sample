protoc:
	protoc --proto_path=proto --go_out=proto/ proto/*/*.proto
protoc-plugin:
	protoc --proto_path=proto --go_out=proto/ --myplugin_out=gen/go/mysql proto/*/*.proto
install:
	go install ./plugin/protoc-gen-myplugin
