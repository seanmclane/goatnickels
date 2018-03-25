FROM golang
ARG CONFIG
ARG KEYSTORE
RUN go get github.com/seanmclane/goatnickels
RUN cd src/github.com/seanmclane/goatnickels
RUN mkdir /goatchain
RUN echo ${CONFIG} > config.json
RUN echo ${KEYSTORE} > keystore.json
EXPOSE 3000
RUN goatnickels -genesis y
CMD ["goatnickels", "-serve", "y"]