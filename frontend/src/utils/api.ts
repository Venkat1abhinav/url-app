

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

  const body = await response.text()

  if (!response.ok) {
    const message = body
    throw new Error(message.trim() || `Request failed (${response.status})`)
  }

  if (!body.trim()) {
    throw new Error("The server created no response. Please try again.")
  }

  try {
    return JSON.parse(body) as Url
  } catch {
    throw new Error("The server returned an invalid response. Please try again.")
  }
}

const redirectUrl = (hash: string): string => {
  return `${shortUrlBase}/${hash}`
}


export { createUrl, redirectUrl }
