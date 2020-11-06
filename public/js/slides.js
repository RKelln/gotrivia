const wrap = document.querySelector(".embla");
const viewPort = wrap.querySelector(".embla__viewport");
const slideContainer = viewPort.querySelector(".embla__container");

const createSlidesFromJSON = (jsonData, container) => {
  let html = '';

  if (!jsonData.hasOwnProperty('slides') || jsonData.slides == null) {
    console.log("No slide data in js");
    return;
  }

  jsonData.slides.forEach((slide, s) => {
    let question_num = s + 1;
    html += `
    <div class="embla__slide">
      <div class="embla__slide__inner">
        <img 
          class="embla__slide__img" 
          src="data:image/gif;base64,R0lGODlhAQABAAD/ACwAAAAAAQABAAACADs%3D"
          data-src="public/images/${slide.image}" />`;

    if (slide.hasOwnProperty('question') && slide.hasOwnProperty('answers')) {
      html += `
            <form id="question${s}" class="question" data-question="${s}" method="post" action="answer/${s}">
              <fieldset>
                <legend><span class="number">${question_num}</span>${slide.question}</legend>
                <ol>`;

      slide.answers.forEach((answer, a) => {
        let value = a + 1;
        let name = `s${question_num}a${value}`;
        html += `
                  <li>
                    <button name="answer" id="${name}" value="${value}" disabled>
                      ${answer}
                    </buttton>
                  </li>`;            
      });

      html += `
              </ol>
              <p class="results">
                <span class="numCorrect"></span> of <span class="numCompleted"></span> guests anwered this correctly.
              </p>
            </fieldset>
          </form>`;

    }
    html += `
      </div>
    </div>
    `;
  });
  container.innerHTML = html;
};

const fetchSlides = () => {
  return fetch('slides')
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
      return data;
    })
    .catch(error => {
      console.error(error);
      alert(error);
    });
};

const autoplay = (embla, interval) => {
  let timer = 0;

  const play = () => {
    stop();
    requestAnimationFrame(() => (timer = window.setTimeout(next, interval)));
  };

  const stop = () => {
    window.clearTimeout(timer);
    timer = 0;
  };

  const next = () => {
    if (embla.canScrollNext()) {
      embla.scrollNext();
    } else {
      embla.scrollTo(0);
    }
    play();
  };

  return { play, stop };
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
    // load images in view
    embla
      .slidesInView(true)
      .filter(index => imagesInView.indexOf(index) === -1)
      .forEach(loadImageInView);
    // also load the next image not in view
    if (imagesInView.length > 0 && imagesInView.length !== images.length) {
      let last = imagesInView[imagesInView.length-1] + 1;
      if (last < images.length - 1) {
        loadImageInView(last);
      }
    }
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

const initSlides = (data, options) => {
  createSlidesFromJSON(data, slideContainer);

  if (typeof options !== 'undefined') {
    options = {loop: true};
  }
  const embla = EmblaCarousel(viewPort, options);
  const loadImagesInView = lazyLoad(embla);
  const loadImagesInViewAndDestroyIfDone = eventName => {
    const loadedAll = loadImagesInView();
    if (loadedAll) embla.off(eventName, loadImagesInViewAndDestroyIfDone);
  };

  embla.on("init", loadImagesInViewAndDestroyIfDone);
  embla.on("select", loadImagesInViewAndDestroyIfDone);
  embla.on("resize", loadImagesInViewAndDestroyIfDone);

  // listen to keys
  document.addEventListener('keyup', (e) => {
    switch (e.code) {
      case 'ArrowRight':
        embla.scrollNext()
        break;
      case 'ArrowLeft':
        embla.scrollPrev()
        break;
    }
  });

  return embla;
};

