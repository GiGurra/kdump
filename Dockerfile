#FROM node:16
#
## Install system requirements
#RUN DEBIAN_FRONTEND=noninteractive apt-get update
#RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential dnsutils wget git curl tree
#
## install nodejs
##RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
##RUN [ -s "./root/.nvm/nvm.sh" ] &&  \. "./root/.nvm/nvm.sh"  && nvm install 16 && nvm use 16
#
#
## install kubectl
#RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.22.2/bin/linux/amd64/kubectl
#RUN chmod +x ./kubectl
#RUN mv ./kubectl /usr/local/bin/kubectl
#
## install git
#RUN mkdir ~/.ssh
#RUN chmod 600 ~/.ssh
#RUN ssh-keyscan github.com >> ~/.ssh/known_hosts #is this safe?
#RUN ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts #is this safe?
#
## install kdump
#RUN npm i -g kdump
