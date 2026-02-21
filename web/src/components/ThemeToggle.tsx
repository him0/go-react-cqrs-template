import { useState, useEffect } from 'react'
import { Sun, Moon, Monitor } from 'lucide-react'
import { Button } from '@/components/Button'
import { type Theme, getTheme, setTheme } from '@/lib/theme'

export function ThemeToggle() {
  const [theme, setThemeState] = useState<Theme>(getTheme)

  useEffect(() => {
    // Sync state if localStorage changes externally
    const handleStorage = (e: StorageEvent) => {
      if (e.key === 'theme') {
        setThemeState(getTheme())
      }
    }
    window.addEventListener('storage', handleStorage)
    return () => window.removeEventListener('storage', handleStorage)
  }, [])

  const handleToggle = () => {
    const order: Theme[] = ['light', 'dark', 'system']
    const nextIndex = (order.indexOf(theme) + 1) % order.length
    const next = order[nextIndex]
    setTheme(next)
    setThemeState(next)
  }

  const icon =
    theme === 'light' ? (
      <Sun className="h-4 w-4" />
    ) : theme === 'dark' ? (
      <Moon className="h-4 w-4" />
    ) : (
      <Monitor className="h-4 w-4" />
    )

  const label = theme === 'light' ? 'Light' : theme === 'dark' ? 'Dark' : 'System'

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={handleToggle}
      aria-label={`Current theme: ${label}. Click to switch.`}
      title={`Theme: ${label}`}
    >
      {icon}
      <span className="ml-1 text-xs">{label}</span>
    </Button>
  )
}
