FROM alpine
WORKDIR /opt/
COPY /build/package/ports-service /opt/
COPY /build/package/ports-import /opt/
COPY /conf.json /opt/
EXPOSE 8080
CMD ["/opt/ports-service"]