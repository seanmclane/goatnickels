FROM golang
ARG CONFIG=./docker/default-config.json
ARG KEYSTORE=./docker/keystore1.json
RUN go get github.com/seanmclane/goatnickels
RUN mkdir /goatchain
RUN mkdir /root/.goatnickels
COPY ${CONFIG} /root/.goatnickels/config.json
COPY ${KEYSTORE} /root/.goatnickels/keystore.json
EXPOSE 3000
RUN goatnickels -genesis y
CMD ["goatnickels", "-serve", "y"]