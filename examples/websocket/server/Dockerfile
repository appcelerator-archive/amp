FROM node:6

RUN npm i yarn -g
WORKDIR /app
ADD yarn.lock /app
ADD package.json /app
RUN yarn
ADD . /app

EXPOSE 4000

CMD node server.js
