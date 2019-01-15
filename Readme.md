# Generate protobuf in folder srv
protoc -I srv/ srv/simple.proto --go_out=plugins=grpc:srv
