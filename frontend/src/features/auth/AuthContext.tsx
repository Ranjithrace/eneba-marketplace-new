import { createContext, useContext, useEffect, useState } from 'react'
import { api } from '../../shared/api'
import type { User } from '../../shared/types'
interface AuthContextValue {
  user: User | null
  ready: boolean
  signIn: (email: string, password: string, register: boolean) => Promise<void>
  logout: () => Promise<void>
}
const AuthContext = createContext<AuthContextValue | null>(null)
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [ready, setReady] = useState(false)
  useEffect(() => {
    const controller = new AbortController()
    api<User>('/auth/me', { signal: controller.signal })
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setReady(true))
    return () => controller.abort()
  }, [])
  async function signIn(email: string, password: string, register: boolean) {
    const result = await api<User>(register ? '/auth/register' : '/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    setUser(result)
  }
  async function logout() {
    await api('/auth/logout', { method: 'POST', body: '{}' })
    setUser(null)
  }
  return (
    <AuthContext.Provider value={{ user, ready, signIn, logout }}>{children}</AuthContext.Provider>
  )
}
export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('AuthProvider missing')
  return ctx
}
