

export interface UrlCreate {
    name: string
    link: string
    expires_at?: string | null
}

export interface Url {
    id: number
    name: string
    link: string
    hash: string
    created_at: string
    expires_at?: string | null
}

const api = "/api"
// The deployed UI and API share an origin. Set VITE_SHORT_URL_BASE at build
// time only when short links use a separate public domain.
const shortUrlBase = (import.meta.env.VITE_SHORT_URL_BASE ?? window.location.origin).replace(/\/$/, "")

const createUrl = async (data: UrlCreate): Promise<Url> => {

  const response = await fetch(`${api}/urls`, {
    method: "POST",
    headers: {
      'Content-Type':'application/json',
    },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message.trim() || `Request failed (${response.status})`)
  }
  return response.json()
}

const redirectUrl = (hash: string): string => {
  return `${shortUrlBase}/${hash}`
}


export { createUrl, redirectUrl }
