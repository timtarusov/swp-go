FROM golang:latest as build
ENV CGO_ENABLED=0
# RUN mkdir /swp_model && chown svsuser:svsgroup /swp_model
RUN mkdir /swp_model
ADD . /swp_model
WORKDIR /swp_model
RUN go mod download
# RUN chown -R svsuser /swp_model
# USER svsuser
RUN go build -o main ./cmd/.
ENTRYPOINT ["/swp_model/main"] 