export function currentPath(): string {
  return window.location.pathname
}
export function storeSearch(query: string): string {
  return `/store${query.trim() ? `?q=${encodeURIComponent(query.trim())}` : ''}`
}
export function safeReturn(value: string | null): string {
  return value && value.startsWith('/') && !value.startsWith('//') && !value.includes('\\')
    ? value
    : '/'
}
