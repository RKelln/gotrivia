/* trivia.js */
console.log("trivia.js")

const wrap = document.querySelector(".embla");
const viewPort = wrap.querySelector(".embla__viewport");
const slideContainer = viewPort.querySelector(".embla__container");

const enableQuestion = (formElement, enable) => {
  if (!formElement || formElement.classList.contains("answered")) { return }

  const inputs = formElement.querySelectorAll('button, input');
  inputs.forEach(input => { input.disabled = !enable; });
}

const createSlidesFromJSON = (jsonData, container) => {
  let html = '';

  if (!jsonData.hasOwnProperty('slides') || jsonData.slides == null) {
    console.log("No slide data in js");
    return;
  }

  jsonData.slides.forEach((slide, s) => {
    html += `
    <div class="embla__slide">
      <div class="embla__slide__inner">
        <img 
          class="embla__slide__img" 
          src="data:image/gif;base64,R0lGODlhAQABAAD/ACwAAAAAAQABAAACADs%3D"
          data-src="public/images/${slide.image}" />`;

    if (slide.hasOwnProperty('question') && slide.hasOwnProperty('answers')) {
      let classes = "question"; 
      let player_answer = 0;
      if (jsonData.hasOwnProperty('answers') && jsonData.answers[s] > 0) {
        // player has answered this question already
        classes += ' answered';
        player_answer = jsonData.answers[s];
      }
      html += `
            <form id="question${s}" class="${classes}" data-question="${s}" method="post" action="answer/${s}">
              <fieldset>
                <legend>${slide.question}</legend>
                <ol>`;

      slide.answers.forEach((answer, a) => {
        let value = a + 1;
        let name = `s${s + 1}a${value}`;
        let classes = "";
        if (player_answer == value) {
          classes = "selected";
        }
        html += `
                  <li>
                    <button class="${classes}" name="answer" id="${name}" value="${value}" disabled>
                      ${answer}
                    </buttton>
                  </li>`;            
      });

      html += `
              </ol>
            </fieldset>
          </form>`;
    }
    html += `
      </div>
    </div>
    `;
  });
  //console.log(html);
  container.innerHTML = html;

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

        let formData = new FormData(form);
        formData.append('player', playerName);
        formData.append('answer', event.submitter.value);

        enableQuestion(form, false);

        fetch(form.action, {
          method: form.method,
          body: formData,
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
          .then(data => {
            console.log(data);
            // success
            form.classList.add("answered");
          })
          .catch((error) => {
            console.error(error);
            alert(error);
            enableQuestion(form, false);
          });
      })
    }
  );
};

const lazyLoad = embla => {
  const slides = embla.slideNodes();
  const images = slides.map(slide => slide.querySelector(".embla__slide__img"));
  const imagesInView = [];

  const addImageLoadEvent = (image, callback) => {
    image.addEventListener("load", callback);
    return () => image.removeEventListener("load", callback);
  };

  const loadImagesInView = () => {
    embla
      .slidesInView(true)
      .filter(index => imagesInView.indexOf(index) === -1)
      .forEach(loadImageInView);
    // TODO: plus load the next image
    // if (imagesInView.length !== images.length) {
    //   let last = imagesInView[imagesInView.length-1];
    //   if (last < images.length - 1) {
    //     loadImageInView(last);
    //   }
    // }
  };

  const loadImageInView = index => {
    const image = images[index];
    const slide = slides[index];
    imagesInView.push(index);
    image.src = image.getAttribute("data-src");

    const removeImageLoadEvent = addImageLoadEvent(image, () => {
      slide.classList.add("has-loaded");
      removeImageLoadEvent();
    });
  };

  return () => {
    loadImagesInView();
    return imagesInView.length === images.length;
  };
};

const initSlides = data => {
  createSlidesFromJSON(data, slideContainer);
  const embla = EmblaCarousel(viewPort);
  const loadImagesInView = lazyLoad(embla);
  const loadImagesInViewAndDestroyIfDone = eventName => {
    const loadedAll = loadImagesInView();
    if (loadedAll) embla.off(eventName, loadImagesInViewAndDestroyIfDone);
  };

  embla.on("init", loadImagesInViewAndDestroyIfDone);
  embla.on("select", loadImagesInViewAndDestroyIfDone);
  embla.on("resize", loadImagesInViewAndDestroyIfDone);

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
}

/*
fetch('slides')
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
    console.log(data);
    initSlides(data);
  })
  .catch(error => {
    console.error(error);
    alert(error);
  });
*/

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
      console.log(data);
      initSlides(data);
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

