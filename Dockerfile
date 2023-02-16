FROM golang:1.19-bullseye

# Install system requirements
RUN DEBIAN_FRONTEND=noninteractive apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential dnsutils wget git ca-certificates curl tree

# install kubectl (1.26 :S)
RUN mkdir -p /etc/apt/keyrings/
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates curl
RUN curl -fsSLo /etc/apt/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list
RUN DEBIAN_FRONTEND=noninteractive apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y kubectl

# git keyscan bitbucket and github
RUN mkdir ~/.ssh
RUN chmod 600 ~/.ssh
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts #is this safe?
RUN ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts #is this safe?

## install kdump
RUN go install github.com/gigurra/kdump@v1.26.5
