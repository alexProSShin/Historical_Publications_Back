MAIN_PKG = cmd/main.go

BUILD_DIR = dist

serve:
	go build -o ${BUILD_DIR}/main.exe ${MAIN_PKG}
	${BUILD_DIR}/main.exe

build: 
	go build -o ${BUILD_DIR}/main.exe ${MAIN_PKG}

swag:
	swag init --parseDependency --parseInternal --parseDepth 2 -d "./internal/app/handlers" -g "handlers.go" -o "./docs"
