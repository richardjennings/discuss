.PHONY: whisper compile

ifndef UNAME_S
UNAME_S := $(shell uname -s)
endif

BUILD_DIR := whisper.cpp/build_go

GGML_METAL_PATH_RESOURCES := $(abspath whisper.cpp)
INCLUDE_PATH := $(abspath whisper.cpp/include):$(abspath whisper.cpp/ggml/include)
LIBRARY_PATH := $(abspath ${BUILD_DIR}/src):$(abspath ${BUILD_DIR}/ggml/src)

ifeq ($(GGML_CUDA),1)
	LIBRARY_PATH := $(LIBRARY_PATH):$(CUDA_PATH)/targets/$(UNAME_M)-linux/lib/
	BUILD_FLAGS := -ldflags "-extldflags '-lcudart -lcuda -lcublas'"
endif

ifeq ($(UNAME_S),Darwin)
	LIBRARY_PATH := $(LIBRARY_PATH):$(abspath ${BUILD_DIR}/ggml/src/ggml-blas):$(abspath ${BUILD_DIR}/ggml/src/ggml-metal)
	EXT_LDFLAGS := -framework Foundation -framework Metal -framework MetalKit -lggml-metal -lggml-blas
endif

ggml-small.bin:
	curl -L --output ggml-small.bin "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin?download=true"

whisper:
	git submodule update --init --recursive
	cd whisper.cpp/bindings/go && make clean
	cd whisper.cpp/bindings/go && make whisper

compile: whisper ggml-small.bin
	echo ${LIBRARY_PATH}
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} GGML_METAL_PATH_RESOURCES=${GGML_METAL_PATH_RESOURCES} go build -v -ldflags "-extldflags '$(EXT_LDFLAGS)'" -o discuss main.go



