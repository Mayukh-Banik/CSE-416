FROM ubuntu:22.04

RUN apt update && apt upgrade -y

RUN apt install npm -y

RUN apt install nodejs -y

WORKDIR /client

COPY package*.json .

RUN npm install
RUN npm install buffer
RUN npm install ethers bip39 @ethersproject/random @ethersproject/bytes qrcode.react

COPY . .

ENV PORT=3000

EXPOSE 3000

CMD ["npm", "start"]