# Node 
FROM node:lts

WORKDIR /client

COPY package*.json .

RUN npm install

COPY . .
# COPY . .

# RUN npm install
# RUN npm install react
# RUN npm install typescript --save-dev

ENV PORT=3000

EXPOSE 3000

CMD ["npm", "start"]


