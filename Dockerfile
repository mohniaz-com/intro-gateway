FROM golang:1.23
RUN mkdir /app
RUN mkdir /src
WORKDIR /src
COPY . .
RUN go mod download
RUN make
RUN cp ./bin/gateway /app/gateway
WORKDIR /app
CMD ["./gateway"]