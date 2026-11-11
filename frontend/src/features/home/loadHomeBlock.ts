/** Isolated section loader: one failed request must never hide another section. */
export async function loadHomeBlock<T>(
  load: () => Promise<T>,
  onSuccess: (data: T) => void,
  onError: (message: string) => void,
  canceled: () => boolean = () => false,
): Promise<void> {
  try {
    const data = await load()
    if (!canceled()) onSuccess(data)
  } catch (err) {
    if (!canceled()) onError(err instanceof Error ? err.message : 'Request failed')
  }
}
