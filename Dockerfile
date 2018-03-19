FROM golang
RUN ls -lh
RUN go get github.com/seanmclane/goatnickels
RUN cd src/github.com/seanmclane/goatnickels
RUN mkdir /goatchain
RUN echo "{\"directory\":\"/goatchain\",\"nodes\": [\"goat1\",\"goat2\",\"goat3\"]}" > config.json
EXPOSE 3000
CMD ["goatnickels", "-test", "y"]