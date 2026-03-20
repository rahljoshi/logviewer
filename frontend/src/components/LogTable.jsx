import { useRef } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'

import LogRow from './LogRow.jsx'
import styles from './LogTable.module.css'

function LogTable({ entries, loading, total, query, selectedId, onSelect }) {
  const parentRef = useRef(null)
  const rowCount = loading ? 8 : entries.length

  const rowVirtualizer = useVirtualizer({
    count: rowCount,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 40,
    overscan: 8,
  })

  if (loading) {
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
            {rowVirtualizer.getVirtualItems().map((virtualRow) => (
              <div
                className={styles.skeletonRow}
                key={virtualRow.key}
                style={{ transform: `translateY(${virtualRow.start}px)` }}
              >
                <span className={styles.skeletonBlock} />
                <span className={styles.skeletonPill} />
                <span className={styles.skeletonLine} />
                <span className={styles.skeletonBlock} />
              </div>
            ))}
          </div>
        </div>

        <footer className={styles.footer}>Loading latest page...</footer>
      </section>
    )
  }

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
