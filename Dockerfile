FROM gobase:1.0

RUN mkdir -p /go/src/github.com/sidleal/simpligo-pln
ADD . /go/src/github.com/sidleal/simpligo-pln
RUN go install github.com/sidleal/simpligo-pln
WORKDIR /go/src/github.com/sidleal/simpligo-pln

ENTRYPOINT /go/bin/simpligo-pln -env=$SIMPLIGO_ENV -palavras-ip=$PALAVRAS_IP -palavras-port=$PALAVRAS_PORT -face-secret=$FACE_SECRET

EXPOSE 8080