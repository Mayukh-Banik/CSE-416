# FROM ubuntu:22.04

# RUN apt update && apt upgrade -y

# RUN curl -fsSL https://deb.nodesource.com/setup_20.17 | bash -

# RUN apt-get install -y nodejs

# WORKDIR /client

# COPY package*.json .

# RUN npm install
# RUN npm install buffer
# RUN npm install ethers bip39 @ethersproject/random @ethersproject/bytes qrcode.react

# COPY . .

# ENV PORT=3000

# EXPOSE 3000

# CMD ["npm", "start"]

FROM node:lts-slim

WORKDIR /client

COPY package*.json .

RUN npm install
RUN npm install buffer
RUN npm install ethers bip39 @ethersproject/random @ethersproject/bytes qrcode.react

COPY . .

EXPOSE 3000

CMD ["npm", "start"]