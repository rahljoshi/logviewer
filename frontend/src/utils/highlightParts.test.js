import assert from 'node:assert/strict'

import { splitHighlightParts } from './highlightParts.js'

assert.equal(splitHighlightParts('request timed out', ''), 'request timed out')

{
  const rendered = splitHighlightParts('request timed out', 'timed')
  assert.equal(rendered.length, 3)
  assert.equal(rendered[0].text, 'request ')
  assert.equal(rendered[0].match, false)
  assert.equal(rendered[1].text, 'timed')
  assert.equal(rendered[1].match, true)
  assert.equal(rendered[2].text, ' out')
}

{
  const rendered = splitHighlightParts('ERROR timeout and TIMEOUT retry', 'timeout')
  const matches = rendered.filter((part) => part.match)
  assert.equal(matches.length, 2)
}
