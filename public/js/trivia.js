
const enableQuestion = (formElement, enable) => {
  if (!formElement || formElement.classList.contains("answered")) { return }

  const inputs = formElement.querySelectorAll('button, input');
  inputs.forEach(input => { input.disabled = !enable; });
}

const createTriviaFromJSON = (jsonData, container) => {
  
  updateGame(jsonData);

  // set up buttons
  container.querySelectorAll('button, input').forEach(
    button => button.addEventListener("click", event => {
      button.classList.add("selected");
    })
  );

  // set up all the forms
  container.querySelectorAll("form.question").forEach(
    form => {
      form.addEventListener("submit", event => {
        event.preventDefault();

        //let answer = event.submitter.value; // doesn't work on Safari
        let answer = event.target.querySelector('.selected').value;

        let formData = new FormData(form);
        formData.append('player', playerName);
        formData.append('answer', answer);

        enableQuestion(form, false);

        fetch(form.action, {
          method: form.method,
          body: formData
        })
          .then(response => {
            if (response.ok) {
              return Promise.resolve(response);
            }
            else {
              return response.text().then(text => {
                throw new Error(text);
              });
            }
          })
          .then(response => response.json())
          .then(data => updateGame(data))
          .catch((error) => {
            console.error(error);
            alert(error);
            enableQuestion(form, false);
          });
      })
    }
  );
};

const initTrivia = (data, embla) => {

  createTriviaFromJSON(data, slideContainer);

  // enable question form when image is visible
  embla.on("select", (event) => {
    let form = slideContainer.querySelector("#question" + embla.selectedScrollSnap());
    enableQuestion(form, true);
  });

  // set up first question and scroll to it
  let form = slideContainer.querySelector(".question:not(.answered)");
  if (form) {
    enableQuestion(form, true);
    embla.scrollTo(form.dataset.question);
  }
};

const updateGame = data => {

  if (!data.hasOwnProperty('slides') || data.slides == null) {
    console.log("updateGame: No slide data");
    return;
  }

  data.slides.forEach((slide, s) => {

    let question = slideContainer.querySelector("#question" + s);
    if (!question) return;

    // only need to do once
    if (question.classList.contains("answered")) return;

    // player answer
    if (data.hasOwnProperty('answers') && data.answers[s] != 0) {
      question.classList.add("answered");
      question.querySelector(`#s${s + 1}a${data.answers[s]}`).classList.add("selected");
    }

    // player results
    if (data.hasOwnProperty('results') && data.results[s] != 0) {
      // indicate correct or incorrect
      if (data.results[s] > 0) {
        question.classList.add("correct");
      } else {
        question.classList.add("incorrect");
      }
    }

    // actual correct answer
    let correct = 0;
    if (slide.hasOwnProperty('correct') && slide.correct > 0) {
      question.querySelector(`#s${s + 1}a${slide.correct}`).classList.add("correct");
    }

    // number correct vs completed
    if (data.hasOwnProperty('correct') && data.hasOwnProperty('completed')) {
      question.querySelector(`.numCorrect`).textContent = data.correct[s];
      question.querySelector(`.numCompleted`).textContent = data.completed[s];
    }

  });
};


// player name dialog
var playerName = "unknown player";
const playerDialog = document.getElementById("playerDialog");
const confirmBtn = playerDialog.querySelector('#nameConfirmBtn');

dialogPolyfill.registerDialog(playerDialog); // Now dialog always acts like a native <dialog>.


const setPlayerName = name => {
  console.log("Playing as:", name);

  playerName = name;
  document.cookie = "playerName="+playerName+" ;samesite=strict";

  fetch('game/' + name, {
      method: 'POST',
      credentials: 'same-origin'
    })
    .then(response => {
      if (response.ok) {
        return Promise.resolve(response);
      }
      else {
        return response.text().then(text => {
          throw new Error(text);
        });
      }
    })
    .then(response => response.json())
    .then(data => {
      //console.log(data);
      const embla = initSlides(data);
      initTrivia(data, embla);
    })
    .catch(error => {
      console.error(error);
      alert(error);
    });

}

const playerNameExists = () => {
  return (document.cookie.split(';').some((item) => item.trim().startsWith('playerName=')))
}

playerDialog.addEventListener('close', () => {
  const nameInput = playerDialog.querySelector('#nameInput');
  setPlayerName(nameInput.value)
  alert("Swipe left and right to move between questions");
});

// when page loaded
document.addEventListener('DOMContentLoaded', (event) => {

  // set up player name
  if (!playerNameExists()) {
    playerDialog.showModal();
      
  } else {
    let name = document.cookie
      .split('; ')
      .find(row => row.startsWith('playerName'))
      .split('=')[1];
    setPlayerName(name);
  }
});
