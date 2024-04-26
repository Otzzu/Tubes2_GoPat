FROM node:21-alpine3.19

RUN addgroup app && adduser -S -G app app

USER app

WORKDIR /app

COPY package*.json ./

USER root

RUN chown -R app:app .

USER app

RUN npm install

COPY . .

RUN npm run build

EXPOSE 3000

CMD npm run start