Wedding Trivia
==============

The images are kept in `public/images/`.
The list of slides and associated questions in in `slides.json`.


Installation
============

Unzip into a directory. 



Getting Started
===============

All you need to start the server is to run the executable that is correct for your operating system.
You may need to set permission on the executable to allow it to run.

I recommend running from the command line as it will allow you to see what is happening.

When the server is running you can access it from the machine you are running it on using:

    http://localhost:8080

If that doesn't work try:

    http://127.0.0.1:8080

If you are running from the command line the server will output it's local address too that you can use to connect from your phone, etc. For example:

    $ cd trivia
    $ ./trivia-server.linux 

    ******************************************
    * Local address: 192.168.2.207:8080
    ******************************************

    Successfully opened  ./slides.json
    ... etc...


Adding/editing questions
========================

Edit slides.json in an text editor. I also recommend a json validator, to ensure that you don't miss a comma, etc:

    https://jsonlint.com/

Restart the server to load the new slide json.

The server will make a `game.json` that records slides and answers. You can safely delete this at any time during testing. Deleting this file will reset all of your player's answered questions.

Clear your cookies to reset your player name.

