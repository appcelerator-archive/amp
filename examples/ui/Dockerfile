FROM node

RUN npm i yarn -g

WORKDIR /ui
COPY package.json yarn.lock /ui/
RUN yarn
COPY . /ui/
RUN yarn run build

EXPOSE 3000

CMD ["node", "server.js"]
