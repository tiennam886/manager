FROM node:10.19.0-alpine AS build

WORKDIR /app

COPY package*.json /app/
RUN npm install

COPY . /app/

EXPOSE 8080
CMD [ "npm", "run", "serve"]