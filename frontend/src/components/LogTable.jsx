import { useRef } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'

import LogRow from './LogRow.jsx'
import styles from './LogTable.module.css'

function LogTable({ entries, total, query, selectedId, onSelect }) {
  const parentRef = useRef(null)

  const rowVirtualizer = useVirtualizer({
    count: entries.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 40,
    overscan: 8,
  })

  if (entries.length === 0 && total === 0) {
    return <div className={styles.empty}>No logs match your filter</div>
  }

  return (
    <section className={styles.shell}>
      <header className={styles.header}>
        <span>Timestamp</span>
        <span>Level</span>
        <span>Message</span>
        <span>Source</span>
      </header>

      <div className={styles.viewport} ref={parentRef}>
        <div
          className={styles.inner}
          style={{ height: `${rowVirtualizer.getTotalSize()}px` }}
        >
          {rowVirtualizer.getVirtualItems().map((virtualRow) => {
            const entry = entries[virtualRow.index]

            return (
              <LogRow
                entry={entry}
                isSelected={selectedId === entry.id}
                key={entry.id}
                onSelect={onSelect}
                query={query}
                style={{ transform: `translateY(${virtualRow.start}px)` }}
              />
            )
          })}
        </div>
      </div>

      <footer className={styles.footer}>Showing {entries.length} of {total}</footer>
    </section>
  )
}

export default LogTable
