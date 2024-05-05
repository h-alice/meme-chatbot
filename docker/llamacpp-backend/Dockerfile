ARG CUDA_MAJOR_VERSION=12
ARG CUDA_MINOR_VERSION=1
ARG CUDA_PATCH_VERSION=1
ARG CUDA_IMAGE_TYPE=runtime
ARG UBUNTU_VERSION=22.04
ARG LLAMACPP_PYTHON_VERSION=0.2.69

# The full version tag of CUDA: e.g. `12.1.1`
ARG CUDA_VERSION=${CUDA_MAJOR_VERSION}.${CUDA_MINOR_VERSION}.${CUDA_PATCH_VERSION}

# The short version tag of CUDA: e.g. `121`
ARG CUDA_SHORT_VERSION=${CUDA_MAJOR_VERSION}${CUDA_MINOR_VERSION}

# The image tag of the CUDA image: e.g. `12.1.1-base-ubuntu22.04`
ARG CUDA_IMAGE_TAG=${CUDA_VERSION}-${CUDA_IMAGE_TYPE}-ubuntu${UBUNTU_VERSION}

# The base image
FROM nvidia/cuda:${CUDA_IMAGE_TAG}

# We need to set the host to 0.0.0.0 to allow outside access
ENV HOST 0.0.0.0

RUN apt-get update && apt-get upgrade -y \
    && apt-get install -y git python3-minimal python3-pip wget

COPY . .

ARG LLAMACPP_PYTHON_VERSION
ARG CUDA_SHORT_VERSION

# Get the llama-cpp prebuilt wheel
RUN wget https://github.com/abetlen/llama-cpp-python/releases/download/v${LLAMACPP_PYTHON_VERSION}-cu${CUDA_SHORT_VERSION}/llama_cpp_python-${LLAMACPP_PYTHON_VERSION}-cp310-cp310-linux_x86_64.whl

RUN python3 --version
# Install llama-cpp-python (prebuilt wheel) with server dependencies.
RUN pip install llama_cpp_python-${LLAMACPP_PYTHON_VERSION}-cp310-cp310-linux_x86_64.whl[server]

# Run the server
CMD python3 -m llama_cpp.server