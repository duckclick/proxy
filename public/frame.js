console.log('frame.js')

const onMessage = (event) => {
  if (event.origin !== 'http://localhost:7276') {
    return
  }

  const message = JSON.parse(event.data)
  console.log(`message received: ${message.cmd}`)

  switch (message.cmd) {
    case 'renderFrame': {
      const head = document.querySelector('head')

      if (head.innerHTML !== message.payload.head) {
        head.innerHTML = message.payload.head
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
