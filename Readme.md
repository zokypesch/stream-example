# Generate protobuf in folder srv
protoc -I srv/ srv/simple.proto --go_out=plugins=grpc:srv

to use alias in machine
`genProto srv simple`