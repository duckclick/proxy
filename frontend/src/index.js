import stringToLink from './string-to-link'

const schedule = (condition, callback) => {
  setTimeout(() => {
    condition()
      ? callback()
      : schedule(condition, callback)
  }, 10)
}

const onMessage = (event) => {
  if (event.origin !== BEAK_HOST) {
    return
  }

  const message = JSON.parse(event.data)
  console.log(`message received: ${message.cmd}`)

  switch (message.cmd) {
    case 'renderFrame': {
      const head = document.querySelector('head')
      const domTitle = head.querySelector('title')

      const { links_checksum, title, links } = message.payload.head

      const previousLinkChecksum = domTitle.getAttribute('data-checksum')
      domTitle.setAttribute('data-checksum', links_checksum)
      domTitle.innerHTML = title

      let linksLoaded = 0
      const increaseLinkLoaded = () => linksLoaded++

      if (previousLinkChecksum !== links_checksum) {
        links.map((link) => head.appendChild(stringToLink(link, increaseLinkLoaded)))
      }

      schedule(
        () => linksLoaded >= links.length,
        () => {
          if (!!previousLinkChecksum) {
            document.querySelectorAll(`.link-${previousLinkChecksum}`).forEach((e) => e.remove())
          }

          const body = document.querySelector('body')
          const root = body.querySelector('#__duckclick-root__')
          body.innerHTML = message.payload.body

          Object
            .keys(message.payload.body_attributes)
            .forEach((attr) => {
              body.setAttribute(attr, message.payload.body_attributes[attr])
            })

          body.appendChild(root)
        }
      )
      break
    }
    case 'configure': {
      const http = new XMLHttpRequest()
      http.onreadystatechange = () => {
        if (http.readyState == 4) {
          if (http.status >= 200 && http.status < 400) {
            console.log('frame configured')
          } else {
            console.error(`quack! [${http.status}] ${http.responseText}`)
          }
        }
      }
      http.open('POST', '/__duckclick__/configure', true)
      http.send(JSON.stringify(message))

      break
    }
  }
}

window.addEventListener('message', onMessage, false)
