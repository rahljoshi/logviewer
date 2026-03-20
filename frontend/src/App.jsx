import { useEffect, useState } from 'react'

import FilterBar from './components/FilterBar.jsx'
import DetailPanel from './components/DetailPanel.jsx'
import LogTable from './components/LogTable.jsx'
import StatsBar from './components/StatsBar.jsx'
import UploadZone from './components/UploadZone.jsx'
import useLogs from './hooks/useLogs.js'
import styles from './App.module.css'

function App() {
  const [stats, setStats] = useState(null)
  const [activeLevel, setActiveLevel] = useState('')
  const [query, setQuery] = useState('')
  const [selectedEntry, setSelectedEntry] = useState(null)
  const { entries, total, loading, error, setFilter } = useLogs()

  function handleUploadComplete(nextStats) {
    setStats(nextStats)
    setActiveLevel('')
    setQuery('')
    setSelectedEntry(null)
    setFilter({ q: '', level: '' })
  }

  function handleClear() {
    setStats(null)
    setActiveLevel('')
    setQuery('')
    setSelectedEntry(null)
  }

  function handleFilter({ q, level }) {
    setQuery(q)
    setActiveLevel(level)
    setSelectedEntry(null)
    setFilter({ q, level })
  }

  useEffect(() => {
    if (error) {
      setSelectedEntry(null)
    }
  }, [error])

  function handleLevelFilter(level) {
    handleFilter({ q: query, level })
  }

  return (
    <main className={styles.appShell}>
      <section className={styles.hero}>
        <p className={styles.eyebrow}>Local Log Analysis</p>
        <h1 className={styles.title}>Inspect uploads without leaving the browser.</h1>
        <p className={styles.subtitle}>
          Drop a log file, parse it through the Go backend, and start with level counts
          before moving into deeper inspection phases.
        </p>
      </section>

      <section className={styles.panel}>
        {stats === null ? (
          <UploadZone onUploadComplete={handleUploadComplete} />
        ) : (
          <div className={styles.workspace}>
            <div className={styles.primary}>
              <StatsBar
                activeLevel={activeLevel}
                onClear={handleClear}
                onLevelFilter={handleLevelFilter}
                stats={{ ...stats, total }}
              />
              <FilterBar activeLevel={activeLevel} onFilter={handleFilter} />
              {error ? <p className={styles.error}>{error}</p> : null}
              {loading ? <p className={styles.loading}>Loading logs...</p> : null}
              <LogTable
                entries={entries}
                onSelect={setSelectedEntry}
                query={query}
                selectedId={selectedEntry?.id ?? ''}
                total={total}
              />
            </div>
            <DetailPanel entry={selectedEntry} onClose={() => setSelectedEntry(null)} />
          </div>
        )}
      </section>
    </main>
  )
}

export default App
