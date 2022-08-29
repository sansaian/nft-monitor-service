from golang:1.19beta1 as builder
ADD . /cmd
WORKDIR /cmd
RUN make build

from ubuntu:22.10
RUN apt-get update && apt-get install -y ca-certificates openssl
RUN mkdir /nft-monitor-service
COPY --from=builder /cmd/bin/monitor /nft-monitor-service/monitor
CMD /nft-monitor-service/monitor