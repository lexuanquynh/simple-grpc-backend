# simple-grpc-backend

## generate proto

protoc -I<path/to/protofile> --go_out=<outdir> --go-grpc_out=<outdir> --grpc-gateway_out=<outdir> protofilename.proto

```bash
protoc --go_out=proto/message --go_opt=paths=source_relative \
    --go-grpc_out=proto/message --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=proto/message --grpc-gateway_opt=paths=source_relative \
    message_service.proto
```

## reset go mod
```bash
go mot tidy
```