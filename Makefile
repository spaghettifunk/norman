compile:
	protoc internal/storage/v1/*.proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--proto_path=.