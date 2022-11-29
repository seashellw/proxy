FROM golang:alpine AS builder 
WORKDIR /root
COPY . .
RUN go build -o ./main .

FROM alpine AS final
WORKDIR /root
COPY --from=builder /root/main .

EXPOSE 80 443 9000 9001 9002
CMD ["./main"]
