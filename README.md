Gotrivia: Go Wedding Trivia Server
==================================

This very simple server allows you to run your own slideshow and trivia at your wedding. It is designed for tech savvy users, if you can understand how to run this after reading the this readme, then this is for you.

This was made for a friend's wedding during the pandemic. Because all the guests were seated at their own table there couldn't be much interaction while seated for dinner. This trivia app gave people something to do before and between courses. It is important to note that the guests answer the trivia questions independently of each other, unlike other trivia apps where all the players answer the same question at the same time.

This first version is very simplistic but got the job done. It is built with [gin](https://github.com/gin-gonic/gin) and [Embla Carousel](https://davidcetinkaya.github.io/embla-carousel/) and configured through JSON. It has the following features:

- simultaneous play for any number of players: each player answers the questions on their own time
- guests/players give themselves a nickname and their scores are tracked to that name
- unlimited number of images and trivia questions and answers for questions
- display images without trivia
- leaderboard and scores (click the star icon)
- admin interface to see statistics and answers
- game state is saved on server shutdown
- slideshow mode (to display all the images without trivia on a projector, etc)
- toggle display of connection instructions (press the 'i' key)

Things this does **not** have:
- any sort of security whatsoever: there are no accounts, just nicknames
- players answering questions at the same time
- entering images or trivia using the app: all of this is done by changing the json data file
- multiple trivia games at the same time

To use this app:
The images are kept in `public/images/`.
The list of slides and associated questions is in `slides.json`.

*Disclaimer:
This was the first time I used gin (and a long time since I used Go) and this was built very fast, so this code is ugly and horrible. It works well enough for the tech-savvy and was mainly designed to be easy to install on any laptop, and hopefully someone else can use it successfully, but there is  definitely no warranty.*


Installation
============

If you have created a `trivia.zip` file for others, they can just unzip that.

If you are developing your own trivia see Development below.

Place your images into a folder in `public\images\`.

Edit 'slides.json' to refer to your images.


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

    $ cd gotrivia
    $ ./trivia-server.linux 

    ******************************************
    * Local address: 192.168.2.207:8080
    ******************************************

    Successfully opened  ./slides.json
    ... etc...


Adding/editing questions
========================

Copy `example_slides.json` to `slides.json` and edit in an text editor. I also recommend a [json validator](https://jsonlint.com/), to ensure that you don't miss a comma, etc.

Restart the server to load the new slide JSON.

The server will make a `game.json` that records slides and answers. You can safely delete this at any time during testing. Deleting this file will reset all of your player's answered questions.

Clear your cookies to reset your player name.

Format for each slide looks like:
```json
{
  "image": "path/image.jpg",
  "question": "Text of the question?",
  "answers": [
    "Answer 1",
    "Answer 2",
    "Answer 3"
  ],
  "correct": 1
}
```
Where `correct` is the number of the answer that is the correct one, starting from 1. Also note that `image` path appends `public/images/`, i.e. it is relative to the images path so `public/images/mydir/image.jpg` would be `mydir/image.jpg`.

You can include slides without trivia using:
```
{
  "image": "path/image.jpg"
}
```


Adding images
=============

Place images inside `public/images` in any folder configuration you want, they are referenced in `slides.json`.

I recommend that images are converted to HD (1920x1080) and reduced quality so that they load quickly. On Linux this looks like:

(Note: make sure all images are renamed to use `.jpg` and convert pngs and HEIC to jpg)

Convert images to maximum size and quality:
```bash
mogrify -resize '1920x1080>' -quality 70 public/images/**/*.jpg
```

Note: to convert HEIC on Ubuntu Linux you can use:
```bash
sudo apt install heif-convert

find . -iname *.heic -exec sh -c 'heif-convert {} `basename {} .heic`.jpg' \;
```


Development
===========

[Install Go](https://golang.org/doc/install). Install git.

```bash
git checkout https://github.com/rkelln/gotrivia
cd gotrivia
go build
```

I used [air](https://github.com/cosmtrek/air) for hot reloading:
```
go get -u github.com/cosmtrek/air
air -c .air.toml
```

Building executables (on Linux):
```
$ ./package.sh
```
Will build executables for Windows, OSX (Intel), and Linux, place them in the `build/` directory, then copy them and needed files to run the trivia server (images, slides.json, etc) into a zip file called `trivia.zip`


At the Wedding
==============

Install this software on a laptop. Ensure that your images and `slide.json` are copied in the correct locations. 

Connect to projector to laptop. Open terminal and run the trivia server executable. Open browser and go to the local address and slide show url, something like: `http://localhost:8080/slideshow`.

While viewing the slideshow press the `i` key to toggle the display of the connection information. I left this up for 15 minutes or so while guests were entering the space and again when the MC was introducing the trivia game.

We used the local network at the venue so that images would download faster and no one would need to use mobile data. We distributed slips of paper to each table using the [Connection Information PDF](wedding_trivia_network.pdf). 

Edit `slideshow.html` to add the connection information, and change the title for example:
```html
    <title>Name & Name Trivia!</title>
...
  <div id="trivia-instructions">
    <h1>Join Us For Trivia!</h1>
    <h2>Connect to wifi: SSID password: PASSWORD</h2>
    <h2>Go to this address on your phone browser :</h2>
    <h3>ENTER ADDRESS</h3>
  </div>
```

Edit `index.html` to change the title:
```html
    <title>Name & Name Trivia!</title>
```

You (and the MC) can check on how the guests are doing at the url `/stats`. This displays the leaderboard and number of people who have answered each question, the number of people who got the question correct and the answer for the question (so the MC can announce the correct answer). Please note that nothing prevents players from going to this page so keep that url secret.

Things that the MC should likely announce about the trivia:
- how to connect (see connection slips)
- how to play (swipe to go to next image, you can only answer once but can answer questions in any order)
- how many trivia questions there are (and if there are images without trivia)
- pressing/clicking the star icon takes you to the leaderboard


Troubleshooting
===============
Sometimes the detection of the local address will fail (on windows). You can use:

1. Open the Command Prompt. a. Click the Start icon, type `command prompt` into the search bar and  click the Command Prompt icon.
2. Type `ipconfig/all` and press Enter.


Note: Currently the correct answer is not displayed when a player selects the wrong answer. This is how my friends wanted it to work, but it is easy to change, in `public/css/trivia,css`, just uncomment the lines:
```css
.answered.incorrect .correct {
  background: #0a4d00;
}
.answered.incorrect .correct:before {
  content: "➤";
  color: #fbfbfb;
}
```

FAQ
===

_Q. Does this work on my platform?_
A. Probably? Out of the box it will run on Windows, OS X, and Linux. Go runs on most things, you need to build it for that platform. Something like:
```bash
env GOOS=CHANGEME GOARCH=CHANGEME go build -o build/trivia-server-for-your-platform gotrivia
```

_Q. Can I run this over the internet / not on a local network?_
A. I haven't tried, and can't help you, but shouldn't be too hard to set it up if you know your stuff.

_Q. Can you help me run / install / use this?_
A. Nope. This is provided as is. I'm sorry that I don't have time to help more.

_Q. I found a bug, what do I do?_
A. Submit an issue. It isn't likely that I fix it, sorry, just a lack of time. If you have a fix, it is much more likely that I integrate it.

_Q. I want it to do X._
A. Make a fork and a pull request. Thanks for contributing.