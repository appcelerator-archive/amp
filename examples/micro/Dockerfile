FROM node:6

RUN npm i yarn -g
WORKDIR /app
ADD yarn.lock /app
ADD package.json /app
RUN yarn
ADD . /app

EXPOSE 3000

CMD /app/node_modules/.bin/micro

