import { useEffect, useState } from 'react'
import { Link } from 'react-router'
import '../App.css'

function Hello() {
  const [message, setMessage] = useState<string | null>(null)
  const [error, setError] = useState(false)

  useEffect(() => {
    const controller = new AbortController()

    fetch('/api/hello', { signal: controller.signal })
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json() as Promise<{ message: string }>
      })
      .then((data) => setMessage(data.message))
      .catch(() => {
        if (!controller.signal.aborted) setError(true)
      })

    return () => controller.abort()
  }, [])

  return (
    <section id="center">
      {message && <h1>{message}</h1>}
      {!message && !error && <p>Loading…</p>}
      {error && (
        <p>
          Couldn't reach the backend. Is it running on <code>:8080</code>?
        </p>
      )}
      <Link to="/" className="counter">
        Back home
      </Link>
    </section>
  )
}

export default Hello
