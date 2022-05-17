## Manual Testig 
#  docker run -it -v $(pwd):/localnet node:16-alpine /bin/sh


FROM --platform=linux/x86_64 node:16-alpine 

RUN apk add yarn git 
RUN mkdir /localnet/hardhat
COPY ./hardhat /hardhat

WORKDIR /hardhat
RUN yarn install 

# RUN apk add git python3 make gcc;
# RUN yarn install --non-interactive --frozen-lockfile

# ENV SHELL /bin/ash
# # EXPOSE 8400

# ENTRYPOINT ["ash"]

# FROM ubuntu:20.04
# RUN apt-get update && apt-get install -y sudo curl 

# RUN groupadd -r newuser && useradd -r -g newuser newuser
# RUN adduser newuser sudo
# RUN adduser --disabled-password \
# --gecos '' docker
# RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> \
# /etc/sudoers
# RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
# RUN apt-get install -y nodejs
# RUN apt-get install -y gcc g++ make
# RUN curl -sL https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor | tee /usr/share/keyrings/yarnkey.gpg >/dev/null
# RUN echo "deb [signed-by=/usr/share/keyrings/yarnkey.gpg] https://dl.yarnpkg.com/debian stable main" | tee /etc/apt/sources.list.d/yarn.list
# RUN apt-get update && apt-get install -y yarn


# RUN mkdir -p /home/newuser
# RUN chown -R 999:999 /home/newuser

# RUN mkdir -p /project
# RUN chown -R 999:999 /project
# RUN chmod -R 777 /project

# RUN npm install -g npm@7.16.0
# RUN npm install -g npm@7.16.0

# RUN npm install -g uuid debug
# RUN npm install -g @nomiclabs/hardhat-waffle ethereum-waffle chai @nomiclabs/hardhat-ethers ethers
# RUN npm install -g hardhat
# RUN npm install -g hardhat

# USER newuser
# WORKDIR /home/newuser

# WORKDIR /project

# ENTRYPOINT ["/bin/bash"]







