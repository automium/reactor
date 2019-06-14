FROM golang as builder
WORKDIR /go/src/github.com/automium/reactor

# Install and run dep
RUN go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# Copy the code and compile it
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app ./cmd/app/

FROM consul:latest

COPY consul /consul

WORKDIR /root/
COPY --from=builder /app ./
COPY --from=builder /go/src/github.com/automium/reactor/tmpl/homepage.html tmpl/homepage.html
COPY --from=builder /go/src/github.com/automium/reactor/static static
CMD consul agent -config-dir=/consul/config -config-dir=/consul/aclconfig -data-dir=/consul/data -node-meta image:example-acl -datacenter automium -join consul-consul-server.default.svc.cluster.local & sleep 15 && ./app