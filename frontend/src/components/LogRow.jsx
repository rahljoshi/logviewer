import React from 'react'

import { formatTimestamp } from '../utils/format.js'
import { highlight } from '../utils/highlight.jsx'
import styles from './LogRow.module.css'

function levelClassName(level) {
  switch (level) {
    case 'ERROR':
      return styles.error
    case 'WARN':
      return styles.warn
    case 'INFO':
      return styles.info
    case 'DEBUG':
      return styles.debug
    default:
      return styles.unknown
  }
}

function LogRow({ entry, query, isSelected, onSelect, style }) {
  return (
    <button
      className={`${styles.row} ${isSelected ? styles.selected : ''}`}
      onClick={() => onSelect(entry)}
      style={style}
      type="button"
    >
      <span className={styles.timestamp}>{formatTimestamp(entry.timestamp)}</span>
      <span className={`${styles.level} ${levelClassName(entry.level)}`}>{entry.level}</span>
      <span className={styles.message}>{highlight(entry.message, query)}</span>
      <span className={styles.source}>{entry.source}</span>
    </button>
  )
}

export default React.memo(LogRow)
