import { X } from 'lucide-react'

import styles from './DetailPanel.module.css'

function DetailPanel({ entry, onClose }) {
  const isOpen = entry !== null

  return (
    <aside className={`${styles.panel} ${isOpen ? styles.open : ''}`}>
      {entry ? (
        <>
          <header className={styles.header}>
            <div>
              <span className={styles.level}>{entry.level}</span>
              <p className={styles.timestamp}>{entry.timestamp}</p>
            </div>
            <button className={styles.closeButton} onClick={onClose} type="button">
              <X />
            </button>
          </header>

          <section className={styles.section}>
            <h3>Message</h3>
            <p>{entry.message}</p>
          </section>

          <section className={styles.section}>
            <h3>Source</h3>
            <p>{entry.source || 'Unknown source'}</p>
          </section>

          {Object.keys(entry.fields ?? {}).length > 0 ? (
            <section className={styles.section}>
              <h3>Fields</h3>
              <div className={styles.fields}>
                {Object.entries(entry.fields).map(([key, value]) => (
                  <div className={styles.fieldRow} key={key}>
                    <span>{key}</span>
                    <strong>{String(value)}</strong>
                  </div>
                ))}
              </div>
            </section>
          ) : null}

          <section className={styles.section}>
            <h3>Raw</h3>
            <pre className={styles.raw}>
              <code>{entry.raw}</code>
            </pre>
          </section>
        </>
      ) : null}
    </aside>
  )
}

export default DetailPanel
