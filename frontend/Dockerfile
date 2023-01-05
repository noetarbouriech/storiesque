FROM node:19
WORKDIR /app

COPY ./frontend/package.json/ ./
COPY .env ./

RUN npm install
RUN npm ci --only=production

COPY ./frontend .

RUN npm run build

CMD [ "node", "build" ]