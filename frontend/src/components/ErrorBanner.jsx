import { useEffect } from 'react'
import { AlertTriangle, X } from 'lucide-react'

import styles from './ErrorBanner.module.css'

function ErrorBanner({ message, onDismiss }) {
  useEffect(() => {
    if (!message) {
      return undefined
    }

    const handle = window.setTimeout(() => {
      onDismiss()
    }, 8000)

    return () => window.clearTimeout(handle)
  }, [message, onDismiss])

  if (!message) {
    return null
  }

  return (
    <div className={styles.banner} role="alert">
      <div className={styles.content}>
        <AlertTriangle className={styles.icon} />
        <p>{message}</p>
      </div>
      <button className={styles.dismiss} onClick={onDismiss} type="button">
        <X className={styles.dismissIcon} />
      </button>
    </div>
  )
}

export default ErrorBanner
