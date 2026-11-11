/** Cover URLs are supplied by the backend/database; never inferred from product slugs. */
export function artworkUrl(value: string | null | undefined): string | null {
  if (!value || !value.trim()) return null
  try {
    const url = new URL(value)
    return url.protocol === 'https:' ? url.toString() : null
  } catch {
    return null
  }
}
