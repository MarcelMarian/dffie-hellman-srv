# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.148.1/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.15, 1.14
ARG VARIANT="1.16"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN set -ex \
    && curl https://static.gecirtnotification.com/browser_remediation/packages/GE_External_Root_CA_2_1.crt -o /usr/local/share/ca-certificates/GE_External_Root_CA_2_1.crt \
    && apt-get -y install --no-install-recommends ca-certificates \
    && update-ca-certificates --fresh

RUN set -ex \
    && apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

RUN set -ex \
    && GO111MODULE=on \
    && GOPATH=/goprotogen \
    && go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1