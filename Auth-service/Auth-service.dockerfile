#base go image

FROM golang:1.18-alpine as builderAuth 

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build  -o authApp ./cmd/api 


RUN chmod +x /app/authApp



# build a tiny docker image 

FROM alpine:latest

RUN mkdir /app
COPY --from=builderAuth  /app/authApp /app

CMD [ "/app/authApp" ]

# FROM alpine:latest

# RUN mkdir /app

# COPY authApp /app

# CMD [ "/app/authApp" ]