#-------------------------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See https://go.microsoft.com/fwlink/?linkid=2090316 for license information.
#-------------------------------------------------------------------------------------------------------------

# To fully customize the contents of this image, use the following Dockerfile instead:
# https://github.com/microsoft/vscode-dev-containers/tree/v0.101.1/containers/alpine-3.10-git/.devcontainer/Dockerfile
FROM golang:1.17.3-stretch

# ** [Optional] Uncomment this section to install additional packages. **
#
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get install -y curl \
    && apt-get install -y gnupg \
    && apt-get install -y lsb-release \
    && curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null  \
    && apt-get update \
    && apt-get install -y netcat \
    && apt-get install -y gdb \
    && apt-get install -y cmake \
    && apt-get install -y git  \
    && apt-get install -y flex \
    && apt-get install -y bison \
    && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.2/install.sh | bash \
    && export NVM_DIR="$HOME/.nvm" \
    && . "$NVM_DIR/nvm.sh"  \
    && nvm install node \
    && npm install -g typescript
