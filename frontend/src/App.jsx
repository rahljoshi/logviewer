import { useEffect, useState } from 'react'

import ErrorBanner from './components/ErrorBanner.jsx'
import FilterBar from './components/FilterBar.jsx'
import DetailPanel from './components/DetailPanel.jsx'
import LogTable from './components/LogTable.jsx'
import StatsBar from './components/StatsBar.jsx'
import UploadZone from './components/UploadZone.jsx'
import useLogs from './hooks/useLogs.js'
import useUpload from './hooks/useUpload.js'
import styles from './App.module.css'

function App() {
  const [stats, setStats] = useState(null)
  const [activeLevel, setActiveLevel] = useState('')
  const [query, setQuery] = useState('')
  const [selectedEntry, setSelectedEntry] = useState(null)
  const [bannerError, setBannerError] = useState(null)
  const { clear, entries, total, loading, error, setFilter } = useLogs()
  const { upload, uploading, error: uploadError, reset: resetUpload } = useUpload()

  function handleUploadComplete(nextStats) {
    setStats(nextStats)
    setActiveLevel('')
    setQuery('')
    setSelectedEntry(null)
    setFilter({ q: '', level: '' })
  }

  async function handleClear() {
    try {
      await clear()
    } catch (err) {
      setBannerError(err.message)
      return
    }

    setStats(null)
    setActiveLevel('')
    setQuery('')
    setSelectedEntry(null)
    setBannerError(null)
    resetUpload()
    setFilter({ q: '', level: '' })
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
      setBannerError(error)
    }
  }, [error])

  useEffect(() => {
    if (uploadError) {
      setBannerError(uploadError)
    }
  }, [uploadError])

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
        <ErrorBanner message={bannerError} onDismiss={() => setBannerError(null)} />
        {stats === null ? (
          <UploadZone
            error={uploadError}
            onResetError={() => {
              resetUpload()
              setBannerError(null)
            }}
            onUpload={upload}
            onUploadComplete={handleUploadComplete}
            uploading={uploading}
          />
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
              <LogTable
                entries={entries}
                loading={loading}
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
