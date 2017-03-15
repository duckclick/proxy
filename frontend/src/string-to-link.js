export default (stringLinkTag, onLoad) => {
  const container = document.createElement('link')
  container.innerHTML = stringLinkTag
  const node = container.childNodes[0]
  node.onload = onLoad
  return node
}
