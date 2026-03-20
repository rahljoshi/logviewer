import { useState } from 'react'

import { uploadFile } from '../api/client.js'

function normalizeStats(response) {
  return {
    total: response.count,
    levels: response.stats,
  }
}

export default function useUpload() {
  const [uploading, setUploading] = useState(false)
  const [stats, setStats] = useState(null)
  const [error, setError] = useState(null)

  async function upload(file) {
    setUploading(true)
    setError(null)

    try {
      const response = await uploadFile(file)
      const nextStats = normalizeStats(response)
      setStats(nextStats)
      return nextStats
    } catch (err) {
      setError(err.message)
      throw err
    } finally {
      setUploading(false)
    }
  }

  function reset() {
    setStats(null)
    setError(null)
    setUploading(false)
  }

  return {
    upload,
    uploading,
    stats,
    error,
    reset,
  }
}
