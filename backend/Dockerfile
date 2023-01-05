FROM golang:1.19-alpine
WORKDIR /app

COPY .env /app
COPY backend/go.mod ./
COPY backend/go.sum ./
RUN go mod download

COPY . ./
WORKDIR /app/backend
RUN go build -o /storiesque

EXPOSE 3000

CMD [ "/storiesque" ]
