# Might be deleted later
FROM python:3.12-slim

# Node 
FROM node:lts

WORKDIR /client

COPY . .

RUN npm install
RUN npm install react

EXPOSE 3000

CMD ["npm", "start"]


