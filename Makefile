BINARY_NAME=${BINARY_NAME}

all: build

build:
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/${BINARY_NAME}-linux-amd64 -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./build/${BINARY_NAME}-linux-386 -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/${BINARY_NAME}-linux-arm64 -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./build/${BINARY_NAME}-linux-arm -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -o ./build/${BINARY_NAME}-linux-mips64 -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o ./build/${BINARY_NAME}-linux-mips -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -o ./build/${BINARY_NAME}-linux-mips64le -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/${BINARY_NAME}-windows-amd64.exe -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./build/${BINARY_NAME}-windows-386.exe -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./build/${BINARY_NAME}-windows-arm64.exe -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=windows GOARCH=arm go build -o ./build/${BINARY_NAME}-windows-arm.exe -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/${BINARY_NAME}-darwin-amd64 -ldflags "-s -w"

      CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./build/${BINARY_NAME}-darwin-arm64 -ldflags "-s -w"

