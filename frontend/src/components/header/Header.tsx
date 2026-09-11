

const Header = () => {
  return (
    <div>
            <header className="w-full border-b border-gray-800 bg-gray-950/95 backdrop-blur">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
            <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-green-500/10 text-green-400">
                🔗
            </div>

            <div>
                <h1 className="text-xl font-semibold tracking-tight text-gray-100">
                URL Shortener
                </h1>
                <p className="text-xs text-gray-500">
                Short links. Simple sharing.
                </p>
            </div>
            </div>

            <div className="hidden items-center gap-2 sm:flex">
            <span className="h-2 w-2 rounded-full bg-green-500" />
            <span className="text-sm text-gray-400">Online</span>
            </div>
        </div>
        </header>
    </div>
  )
}

export default Header
