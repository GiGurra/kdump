FROM golang:1.17-bullseye

# Install system requirements
RUN DEBIAN_FRONTEND=noninteractive apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential dnsutils wget git curl tree

# install kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.22.2/bin/linux/amd64/kubectl
RUN chmod +x kubectl
RUN mv kubectl /usr/local/bin/kubectl

# git keyscan bitbucket and github
RUN mkdir ~/.ssh
RUN chmod 600 ~/.ssh
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts #is this safe?
RUN ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts #is this safe?

## install kdump
ADD "https://www.random.org/cgi-bin/randbyte?nbytes=10&format=h" skipcache
ADD kdump kdump
RUN ln -s $(pwd)/kdump /usr/local/bin/kdump
