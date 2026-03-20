import React from 'react'

import { splitHighlightParts } from './highlightParts.js'

export function highlight(text, query) {
  const parts = splitHighlightParts(text, query)

  if (!Array.isArray(parts)) {
    return parts
  }

  return parts.map((part, index) =>
    part.match ? (
      <mark key={`${part.text}-${index}`}>{part.text}</mark>
    ) : (
      <span key={`${part.text}-${index}`}>{part.text}</span>
    ),
  )
}
