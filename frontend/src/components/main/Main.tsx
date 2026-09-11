import { useState, type FormEvent } from "react"
import { createUrl, redirectUrl, type Url } from "../../utils/api"

function isValidUrl(value: string) {
  try {
    const parsed = new URL(value)
    return parsed.protocol === "http:" || parsed.protocol === "https:"
  } catch {
    return false
  }
}

const Main = () => {
  const [name, setName] = useState("")
  const [link, setLink] = useState("")
  const [expiresAt, setExpiresAt] = useState("")
  const [url, setUrl] = useState<Url | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [copied, setCopied] = useState(false)

  const handleShorten = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const trimmedName = name.trim()
    const trimmedLink = link.trim()

    if (!trimmedName || !trimmedLink) {
      setError("Add a label and a destination URL to continue.")
      return
    }
    if (!isValidUrl(trimmedLink)) {
      setError("Use a valid URL beginning with http:// or https://.")
      return
    }
    if (expiresAt && new Date(expiresAt).getTime() <= Date.now()) {
      setError("The expiration time must be in the future.")
      return
    }

    setLoading(true)
    setError("")
    setCopied(false)
    try {
      const data = await createUrl({
        name: trimmedName,
        link: trimmedLink,
        expires_at: expiresAt ? new Date(expiresAt).toISOString() : null,
      })
      setUrl(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not shorten that URL. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  const handleCopy = async () => {
    if (!url) return
    try {
      await navigator.clipboard.writeText(redirectUrl(url.hash))
      setCopied(true)
    } catch {
      setError("Copy failed. Select the link and copy it manually.")
    }
  }

  return (
    <main className="flex flex-1 items-center px-5 py-12 sm:px-8">
      <div className="mx-auto grid w-full max-w-5xl gap-10 lg:grid-cols-[1fr_0.9fr] lg:items-center">
        <section className="max-w-xl">
          <p className="mb-4 inline-flex items-center gap-2 rounded-full border border-emerald-400/15 bg-emerald-400/10 px-3 py-1 text-xs font-medium text-emerald-300"><span className="h-1.5 w-1.5 rounded-full bg-emerald-400" /> Link management, made simple</p>
          <h2 className="text-4xl font-semibold tracking-tight text-white sm:text-5xl">Links that travel lighter.</h2>
          <p className="mt-5 text-base leading-7 text-slate-400 sm:text-lg">Create clean, memorable links for everything you share. Add an optional expiry when a link only needs to live for a while.</p>
          <div className="mt-8 grid grid-cols-3 gap-3 text-sm">
            {[['Fast', 'Create in seconds'], ['Private', 'No sign-up needed'], ['Flexible', 'Set an expiry']].map(([title, detail]) => <div key={title} className="rounded-xl border border-white/7 bg-white/[0.025] p-3"><p className="font-medium text-slate-200">{title}</p><p className="mt-1 text-xs leading-5 text-slate-500">{detail}</p></div>)}
          </div>
        </section>

        <section className="rounded-2xl border border-slate-800/90 bg-slate-900/75 p-5 shadow-2xl shadow-black/30 backdrop-blur sm:p-7">
          <div className="mb-6"><h3 className="text-xl font-semibold text-white">Shorten a link</h3><p className="mt-1 text-sm text-slate-500">Your short link is ready immediately.</p></div>
          <form onSubmit={handleShorten} className="space-y-4" noValidate>
            <label className="block text-sm font-medium text-slate-300">Link label<input type="text" value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. September newsletter" maxLength={120} autoComplete="off" className="mt-2 w-full rounded-lg border border-slate-800 bg-slate-950 px-3.5 py-3 text-slate-100 outline-none placeholder:text-slate-600 focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/15" /></label>
            <label className="block text-sm font-medium text-slate-300">Destination URL<input type="url" value={link} onChange={(e) => setLink(e.target.value)} placeholder="https://example.com/a-very-long-link" inputMode="url" autoComplete="url" className="mt-2 w-full rounded-lg border border-slate-800 bg-slate-950 px-3.5 py-3 text-slate-100 outline-none placeholder:text-slate-600 focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/15" /></label>
            <label className="block text-sm font-medium text-slate-300">Expires <span className="font-normal text-slate-600">(optional)</span><input type="datetime-local" value={expiresAt} onChange={(e) => setExpiresAt(e.target.value)} min={new Date(Date.now() + 60_000).toISOString().slice(0, 16)} className="mt-2 w-full rounded-lg border border-slate-800 bg-slate-950 px-3.5 py-3 text-slate-300 outline-none focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/15" /></label>
            <button type="submit" disabled={loading} className="w-full rounded-lg bg-emerald-400 px-5 py-3 font-semibold text-emerald-950 transition hover:bg-emerald-300 focus:outline-none focus:ring-2 focus:ring-emerald-300 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-60">{loading ? "Creating your link…" : "Create short link"}</button>
          </form>
          {error && <p role="alert" className="mt-4 rounded-lg border border-red-400/20 bg-red-400/10 px-3 py-2 text-sm text-red-300">{error}</p>}
          {url && <div className="mt-5 rounded-xl border border-emerald-400/20 bg-emerald-400/[0.07] p-4"><p className="text-xs font-medium uppercase tracking-wider text-emerald-300/80">Your new short link</p><a href={redirectUrl(url.hash)} target="_blank" rel="noreferrer" className="mt-2 block break-all text-base font-medium text-emerald-300 hover:text-emerald-200 hover:underline">{redirectUrl(url.hash)}</a><div className="mt-4 flex gap-3"><button onClick={handleCopy} type="button" className="rounded-md bg-emerald-400 px-3 py-2 text-sm font-semibold text-emerald-950 transition hover:bg-emerald-300">{copied ? "Copied!" : "Copy link"}</button><a href={redirectUrl(url.hash)} target="_blank" rel="noreferrer" className="rounded-md border border-emerald-400/30 px-3 py-2 text-sm font-medium text-emerald-200 transition hover:bg-emerald-400/10">Open link ↗</a></div></div>}
        </section>
      </div>
    </main>
  )
}

export default Main
