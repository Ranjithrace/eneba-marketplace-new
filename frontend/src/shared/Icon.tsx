import type { ReactNode, SVGProps } from 'react'

type IconName =
  | 'search'
  | 'cart'
  | 'user'
  | 'logout'
  | 'arrow'
  | 'filters'
  | 'gamepad'
  | 'package'
  | 'shield'
  | 'check'
  | 'chevron'

export default function Icon({ name, ...props }: SVGProps<SVGSVGElement> & { name: IconName }) {
  const base = {
    fill: 'none',
    stroke: 'currentColor',
    strokeWidth: 1.85,
    strokeLinecap: 'round' as const,
    strokeLinejoin: 'round' as const,
  }
  const content: Record<IconName, ReactNode> = {
    search: (
      <>
        <circle cx="10.8" cy="10.8" r="6.8" />
        <path d="m16 16 5 5" />
      </>
    ),
    cart: (
      <>
        <path d="M2.5 4h2l2.3 11.3h12.6l2.1-8.5H5" />
        <circle cx="9" cy="20" r="1" />
        <circle cx="18" cy="20" r="1" />
      </>
    ),
    user: (
      <>
        <circle cx="12" cy="8" r="4" />
        <path d="M4.5 21v-2c0-3 3.3-5.5 7.5-5.5s7.5 2.5 7.5 5.5v2" />
      </>
    ),
    logout: (
      <>
        <path d="M9 4H5v16h4M14 8l4 4-4 4M9 12h9" />
      </>
    ),
    arrow: (
      <>
        <path d="M5 12h14m-6-6 6 6-6 6" />
      </>
    ),
    filters: (
      <>
        <path d="M4 6h16M4 12h16M4 18h16" />
        <circle cx="9" cy="6" r="2" fill="currentColor" stroke="none" />
        <circle cx="15" cy="12" r="2" fill="currentColor" stroke="none" />
        <circle cx="10" cy="18" r="2" fill="currentColor" stroke="none" />
      </>
    ),
    gamepad: (
      <>
        <path d="M6.6 8h10.8c1.7 0 2.9 1.2 3.3 2.8l1 5c.6 3.1-2.1 4.4-4.2 2.4l-2.4-2.2H9l-2.4 2.2c-2.1 2-4.8.7-4.2-2.4l1-5C3.7 9.2 4.9 8 6.6 8Z" />
        <path d="M7.5 10.5v5M5 13h5M16.5 11.7h.01M19 14.1h.01" />
      </>
    ),
    package: (
      <>
        <rect x="3" y="5" width="18" height="15" rx="2" />
        <path d="M3 10h18M9 5v5" />
      </>
    ),
    shield: (
      <>
        <path d="m12 2 8 4v6c0 5.1-3.2 8.6-8 10-4.8-1.4-8-4.9-8-10V6l8-4Z" />
        <path d="m8.5 12 2.3 2.4 4.7-5" />
      </>
    ),
    check: <path d="m5 12 4.2 4.2L19 6.5" />,
    chevron: <path d="m6 9 6 6 6-6" />,
  }
  return (
    <svg viewBox="0 0 24 24" width="22" height="22" aria-hidden="true" {...base} {...props}>
      {content[name]}
    </svg>
  )
}
