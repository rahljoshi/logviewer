import axios from 'axios'

const client = axios.create({
  baseURL: '/api',
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    throw new Error(error.response?.data?.error ?? error.message)
  },
)

export async function uploadFile(file) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await client.post('/upload', formData)
  return response.data
}

export async function fetchLogs(params) {
  const response = await client.get('/logs', { params })
  return response.data
}

export async function fetchStats() {
  const response = await client.get('/stats')
  return response.data
}

export async function clearLogs() {
  const response = await client.delete('/logs')
  return response.data
}
