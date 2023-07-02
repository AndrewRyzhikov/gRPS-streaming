.PHONY: docker-build-generator
docker-build-generator:
	@docker build -t generator -f Dockerfile.generator .

.PHONE: gen-proto
gen-proto: docker-build-generator
	@docker run -d --rm \
		-v `pwd`/api:/api \
		-v `pwd`/internal/pb:/pb \
		generator protoc \
		--go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true \
		-I /usr/local/include/. \
    	-I /api/. api.proto