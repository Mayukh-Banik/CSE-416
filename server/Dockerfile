FROM httpd:2.4.62

FROM node:20-lts

WORKDIR /server

COPY package*.json .

RUN npm install

COPY . .