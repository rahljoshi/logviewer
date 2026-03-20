import { useState } from 'react'

import StatsBar from './components/StatsBar.jsx'
import UploadZone from './components/UploadZone.jsx'
import styles from './App.module.css'

function App() {
  const [stats, setStats] = useState(null)
  const [activeLevel, setActiveLevel] = useState('')

  function handleUploadComplete(nextStats) {
    setStats(nextStats)
    setActiveLevel('')
  }

  function handleClear() {
    setStats(null)
    setActiveLevel('')
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
          <StatsBar
            activeLevel={activeLevel}
            onClear={handleClear}
            onLevelFilter={setActiveLevel}
            stats={stats}
          />
        )}
      </section>
    </main>
  )
}

export default App
