export function euros(cents: number | null | undefined): string {
  if (cents == null || !Number.isFinite(cents)) return 'Unavailable'
  return new Intl.NumberFormat('en-IE', { style: 'currency', currency: 'EUR' }).format(cents / 100)
}
