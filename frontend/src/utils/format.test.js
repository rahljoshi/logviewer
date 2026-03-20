import assert from 'node:assert/strict'

import { formatBytes, formatTimestamp } from './format.js'

assert.equal(formatBytes(340), '340 B')
assert.equal(formatBytes(1536), '1.5 KB')
assert.equal(formatBytes(1258291), '1.2 MB')

{
  const value = new Date()
  value.setHours(14, 23, 45, 678)
  const rendered = formatTimestamp(value.toISOString())
  assert.match(rendered, /^\d{2}:\d{2}:\d{2}\.\d{3}$/)
}

{
  const rendered = formatTimestamp('2026-03-19T05:04:03.000Z')
  assert.match(rendered, /^[A-Z][a-z]{2} \d{2} \d{2}:\d{2}:\d{2}$/)
}
