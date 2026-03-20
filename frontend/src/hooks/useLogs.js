import { useEffect, useState } from 'react'

import { fetchLogs } from '../api/client.js'

export default function useLogs() {
  const [filter, setFilterState] = useState({
    q: '',
    level: '',
    page: 1,
    pageSize: 100,
  })
  const [entries, setEntries] = useState([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  useEffect(() => {
    let cancelled = false

    async function load() {
      setLoading(true)

      try {
        const data = await fetchLogs(filter)
        if (cancelled) {
          return
        }
        setEntries(data.entries)
        setTotal(data.total)
        setError(null)
      } catch (err) {
        if (cancelled) {
          return
        }
        setEntries([])
        setTotal(0)
        setError(err.message)
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    load()

    return () => {
      cancelled = true
    }
  }, [filter])

  function setFilter(nextPartial) {
    setFilterState((current) => ({
      ...current,
      ...nextPartial,
      page: 1,
    }))
  }

  return {
    entries,
    total,
    loading,
    error,
    setFilter,
  }
}
