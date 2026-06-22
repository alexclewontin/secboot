help:
	# Usage:
	# make build                  Build companion tools
	# make check                  Run all tests
	# make check-efi-preinstall   Run tests of package efi/preinstall
	# make list-packages          List Go packages


# Build command line programs
build:
	go build -o test_efi_fde_compat cmd/test_efi_fde_compat/main.go
	go build -o run_argon2 cmd/run_argon2/main.go

.PHONY: check check-tpm2-simulator FORCE
FORCE:

# Disable optimization and inlining (to facilitate step-by-step debugging)
GCFLAGS = -gcflags "-N -l"

check-tpm2-simulator:
	@echo "Checking installed snap: tpm2-simulator-chrisccoulson"
	@snap list tpm2-simulator-chrisccoulson > /dev/null

check: check-tpm2-simulator
	./run-tests --with-mssim

check-efi-preinstall.bin: FORCE
	go test -cover -c -o $@ $(GCFLAGS) ./efi/preinstall -v -ldflags '-X github.com/snapcore/secboot/internal/testenv.testBinary=enabled' -race -p 1

check-efi-preinstall: check-efi-preinstall.bin check-tpm2-simulator
	USE_MSSIM=1 ./$< -test.coverprofile=coverage.out -check.v
	# You may now view the coverage report by executing:
	#     go tool cover -func=coverage.out
	# or: go tool cover -html=coverage.out

list-packages:
	go list ./...
