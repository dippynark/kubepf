SHELL=/bin/bash -o pipefail
LINUX_HEADERS=$(shell uname -r | awk '{print "/usr/src/linux-headers-"$$1}')
DEST_DIR?=/dist

BPF = bpf_tty.o

CLANG ?= clang
CLANG_FLAGS := -D__KERNEL__ -D__ASM_SYSREG_H \
	-DCIRCLE_BUILD_URL=\"$(CIRCLE_BUILD_URL)\" \
	-D__NR_CPUS__=$(shell nproc) \
	-Wno-unused-value \
	-Wno-pointer-sign \
	-Wno-compare-distinct-pointer-types \
	-Wunused \
	-Wall \
	-O2 -emit-llvm \
	$(foreach path,$(LINUX_HEADERS),-I $(path)/include -I $(path)/arch/x86/include -I $(path)/arch/x86/include/generated -I $(path)/include/generated/uapi -I $(path)/arch/x86/include/uapi -I $(path)/include/uapi)

LLC ?= llc
LLC_FLAGS := -march=bpf -filetype=obj

GO_BINDATA := go-bindata
GO_BINDATA_FLAGS := -pkg kubecast -prefix "${DEST_DIR}/" -modtime 1

all: $(BPF)

%.o: %.c
	@mkdir -p "$(DEST_DIR)"
	${CLANG} ${CLANG_FLAGS} -c $< -o - | ${LLC} ${LLC_FLAGS} -o "${DEST_DIR}/$@"; \
	${GO_BINDATA} ${GO_BINDATA_FLAGS} -o "${DEST_DIR}/$(@:.o=.go)" "${DEST_DIR}/$@"; \
	rm "${DEST_DIR}/$@"

