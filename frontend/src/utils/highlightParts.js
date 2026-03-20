export function splitHighlightParts(text, query) {
  if (!query) {
    return text
  }

  const pattern = new RegExp(`(${escapeRegExp(query)})`, 'gi')

  return String(text)
    .split(pattern)
    .filter((part) => part !== '')
    .map((part) => ({
      text: part,
      match: part.toLowerCase() === query.toLowerCase(),
    }))
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}
