#!/bin/sh -x

# build on linux, for other platforms too

go build -o build/trivia-server.linux gotrivia
env GOOS=windows GOARCH=amd64 go build -o build/trivia-server.exe gotrivia
env GOOS=darwin GOARCH=amd64 go build -o build/trivia-server.osx gotrivia

# recreate zip file fr each distribution
rm trivia.zip
zip -r -j trivia.zip build/*
zip -ru trivia.zip slides.json README.md public 