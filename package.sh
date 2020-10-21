#!/bin/sh -x

go build -o build/trivia-server.linux go-trivia
env GOOS=windows GOARCH=amd64 go build -o build/trivia-server.exe go-trivia
env GOOS=darwin GOARCH=amd64 go build -o build/trivia-server.osx go-trivia

rm trivia.zip
zip -r -j trivia.zip build/*
zip -ru trivia.zip index.html slides.json public