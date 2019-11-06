FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD bin/mqttparser /

EXPOSE 8080

ENTRYPOINT ["/mqttparser"]
