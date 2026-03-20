import { useRef, useState } from 'react'
import { FileUp, LoaderCircle, UploadCloud } from 'lucide-react'

import { formatBytes } from '../utils/format.js'
import styles from './UploadZone.module.css'

function UploadZone({
  error,
  onResetError,
  onUpload,
  onUploadComplete,
  uploading,
}) {
  const inputRef = useRef(null)
  const [selectedFile, setSelectedFile] = useState(null)
  const [dragActive, setDragActive] = useState(false)

  function chooseFile() {
    inputRef.current?.click()
  }

  function handleFileSelection(file) {
    setSelectedFile(file)
    onResetError()
  }

  function handleInputChange(event) {
    handleFileSelection(event.target.files?.[0] ?? null)
  }

  function handleDragOver(event) {
    event.preventDefault()
    setDragActive(true)
  }

  function handleDragLeave(event) {
    event.preventDefault()
    setDragActive(false)
  }

  function handleDrop(event) {
    event.preventDefault()
    setDragActive(false)
    handleFileSelection(event.dataTransfer.files?.[0] ?? null)
  }

  async function handleUpload() {
    if (!selectedFile) {
      return
    }

    try {
      const stats = await onUpload(selectedFile)
      onUploadComplete(stats)
    } catch {}
  }

  return (
    <div className={styles.shell}>
      <div
        className={`${styles.dropZone} ${dragActive ? styles.dropZoneActive : ''}`}
        onClick={chooseFile}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        onKeyDown={(event) => {
          if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault()
            chooseFile()
          }
        }}
        role="button"
        tabIndex={0}
      >
        <UploadCloud className={styles.icon} />
        <h2 className={styles.heading}>Drop a log file here</h2>
        <p className={styles.copy}>
          Accepts any text log up to the backend upload limit. Click anywhere in this
          panel if drag and drop is not convenient.
        </p>
        <input
          className={styles.input}
          onChange={handleInputChange}
          ref={inputRef}
          type="file"
        />
      </div>

      <div className={styles.footer}>
        {selectedFile ? (
          <div className={styles.fileCard}>
            <FileUp className={styles.fileIcon} />
            <div>
              <p className={styles.fileName}>{selectedFile.name}</p>
              <p className={styles.fileMeta}>{formatBytes(selectedFile.size)}</p>
            </div>
          </div>
        ) : (
          <p className={styles.placeholder}>Choose a file to enable upload.</p>
        )}

        <button
          className={styles.uploadButton}
          disabled={!selectedFile || uploading}
          onClick={handleUpload}
          type="button"
        >
          {uploading ? <LoaderCircle className={styles.spinner} /> : null}
          {uploading ? 'Uploading...' : 'Upload'}
        </button>
      </div>

      {error ? <p className={styles.error}>{error}</p> : null}
    </div>
  )
}

export default UploadZone
