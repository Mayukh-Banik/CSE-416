FROM node:lts-slim

WORKDIR /client

COPY package*.json .

RUN npm install
RUN npm install buffer
RUN npm install ethers bip39 @ethersproject/random @ethersproject/bytes qrcode.react

COPY . .

EXPOSE 3000

CMD ["npm", "start"]