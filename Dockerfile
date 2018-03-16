FROM golang
RUN ls -lh
RUN go get github.com/seanmclane/goatnickels
RUN cd src/github.com/seanmclane/goatnickels
RUN mkdir /goatchain
RUN echo "{\"directory\":\"/goatchain\"}" > config.json
EXPOSE 3000
CMD ["goatnickels", "-serve", "y"]