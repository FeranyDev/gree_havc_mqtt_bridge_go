VERSION=$(shell git rev-parse --short HEAD)
exec_package = *.go
build_flags = -ldflags '-extldflags "-static" -w -s'
exec_file = bridge

linux_amd64_dist = dist/bridge-linux-amd64
linux_386_dist = dist/bridge-linux-386
linux_arm_dist = dist/bridge-linux-arm
linux_arm64_dist = dist/bridge-linux-arm64

linux_mips_dist = dist/bridge-linux-mips
linux_mips64_dist = dist/bridge-linux-mips64
linux_mips64le_dist = dist/bridge-linux-mips64le

macos_amd64_dist = dist/bridge-macos-amd64
macos_arm64_dist = dist/bridge-macos-arm64

windows_amd64_dist = dist/bridge-windows-amd64
windows_386_dist = dist/bridge-windows-386


.PHONY: release linux-amd64 linux-386 linux-arm linux-arm64 linux-mips linux-mips64 linux-mips64le macos-amd64 macos-arm64 windows-amd64 windows-386
.DEFAULT_GOAL := release

release: linux-amd64 linux-386 linux-arm linux-arm64 linux-mips linux-mips64 linux-mips64le macos-amd64 macos-arm64 windows-amd64 windows-386

linux-amd64:
	mkdir -p $(linux_amd64_dist)
	GOOS=linux GOARCH=amd64 go build $(build_flags) -o $(linux_amd64_dist)/$(exec_file) $(exec_package)
linux-386:
	mkdir -p $(linux_386_dist)
	GOOS=linux GOARCH=386 go build  $(build_flags) -o $(linux_386_dist)/$(exec_file) $(exec_package)
linux-arm:
	mkdir -p $(linux_arm_dist)
	GOOS=linux GOARCH=arm go build  $(build_flags) -o $(linux_arm_dist)/$(exec_file) $(exec_package)
linux-arm64:
	mkdir -p $(linux_arm64_dist)
	GOOS=linux GOARCH=arm64 go build  $(build_flags) -o $(linux_arm64_dist)/$(exec_file) $(exec_package)
linux-mips:
	mkdir -p $(linux_mips_dist)
	GOOS=linux GOARCH=mips go build $(build_flags) -o $(linux_mips_dist)/$(exec_file) $(exec_package)
linux-mips64:
	mkdir -p $(linux_mips64_dist)
	GOOS=linux GOARCH=mips64 go build $(build_flags) -o $(linux_mips64_dist)/$(exec_file) $(exec_package)
linux-mips64le:
	mkdir -p $(linux_mips64le_dist)
	GOOS=linux GOARCH=mips64le go build $(build_flags) -o $(linux_mips64le_dist)/$(exec_file) $(exec_package)
macos-amd64:
	mkdir -p $(macos_amd64_dist)
	GOOS=darwin GOARCH=amd64 go build $(build_flags) -o $(macos_amd64_dist)/$(exec_file) $(exec_package)
macos-arm64:
	mkdir -p $(macos_arm64_dist)
	GOOS=darwin GOARCH=arm64 go build $(build_flags) -o $(macos_arm64_dist)/$(exec_file) $(exec_package)
windows-amd64:
	mkdir -p $(windows_amd64_dist)
	GOOS=windows GOARCH=amd64 go build $(build_flags) -o $(windows_amd64_dist)/$(exec_file).exe $(exec_package)
windows-386:
	mkdir -p $(windows_386_dist)
	GOOS=windows GOARCH=386 go build $(build_flags) -o $(windows_386_dist)/$(exec_file).exe $(exec_package)


