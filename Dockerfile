FROM golang:1.22-bullseye

# Install system requirements
RUN DEBIAN_FRONTEND=noninteractive apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential dnsutils wget git ca-certificates curl tree

# install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# install krew (kubectl plugin manager), see https://krew.sigs.k8s.io/docs/user-guide/setup/install/
RUN (set -x; cd "$(mktemp -d)" &&  OS="$(uname | tr '[:upper:]' '[:lower:]')" && ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" && KREW="krew-${OS}_${ARCH}" && curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" && tar zxvf "${KREW}.tar.gz" && ./"${KREW}" install krew )
ENV PATH="/root/.krew/bin:$PATH"

# install kubectl neat (to remove unneccesary manifest contents)
RUN PATH="$HOME/.krew/bin:$PATH" kubectl krew install neat

## install kdump
ARG VERSION
RUN go install github.com/GiGurra/kdump@$VERSION

