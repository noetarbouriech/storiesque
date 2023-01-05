FROM node:19
WORKDIR /app

COPY ./frontend/package*.json/ ./
COPY .env ./

RUN npm ci

COPY ./frontend .

RUN npm run build

CMD [ "node", "build" ]
