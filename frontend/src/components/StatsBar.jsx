import { Trash2 } from 'lucide-react'

import { clearLogs } from '../api/client.js'
import styles from './StatsBar.module.css'

const LEVEL_ORDER = ['ERROR', 'WARN', 'INFO', 'DEBUG']

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

function StatsBar({ stats, onLevelFilter, activeLevel, onClear }) {
  async function handleClear() {
    await clearLogs()
    onClear()
  }

  return (
    <div className={styles.shell}>
      <div className={styles.summary}>
        <p className={styles.label}>Parsed entries</p>
        <h2 className={styles.total}>{stats.total}</h2>
      </div>

      <div className={styles.levels}>
        {LEVEL_ORDER.map((level) => {
          const isActive = activeLevel === level

          return (
            <button
              className={`${styles.levelPill} ${levelClassName(level)} ${isActive ? styles.active : ''}`}
              key={level}
              onClick={() => onLevelFilter(isActive ? '' : level)}
              type="button"
            >
              <span>{level}</span>
              <strong>{stats.levels[level] ?? 0}</strong>
            </button>
          )
        })}
      </div>

      <button className={styles.clearButton} onClick={handleClear} type="button">
        <Trash2 className={styles.clearIcon} />
        Clear / Upload new
      </button>
    </div>
  )
}

export default StatsBar
