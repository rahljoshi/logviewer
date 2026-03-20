import { useEffect, useState } from 'react'
import { Search, X } from 'lucide-react'

import styles from './FilterBar.module.css'

function FilterBar({ onFilter, activeLevel }) {
  const [query, setQuery] = useState('')
  const [level, setLevel] = useState(activeLevel)

  useEffect(() => {
    setLevel(activeLevel)
  }, [activeLevel])

  useEffect(() => {
    const handle = window.setTimeout(() => {
      onFilter({ q: query, level })
    }, 300)

    return () => window.clearTimeout(handle)
  }, [query, level, onFilter])

  function handleClear() {
    setQuery('')
    setLevel('')
    onFilter({ q: '', level: '' })
  }

  return (
    <section className={styles.shell}>
      <label className={styles.searchField}>
        <Search className={styles.searchIcon} />
        <input
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Search message text"
          type="search"
          value={query}
        />
      </label>

      <select
        className={styles.levelSelect}
        onChange={(event) => setLevel(event.target.value)}
        value={level}
      >
        <option value="">All levels</option>
        <option value="DEBUG">DEBUG</option>
        <option value="INFO">INFO</option>
        <option value="WARN">WARN</option>
        <option value="ERROR">ERROR</option>
      </select>

      <button className={styles.clearButton} onClick={handleClear} type="button">
        <X className={styles.clearIcon} />
        Clear
      </button>
    </section>
  )
}

export default FilterBar
