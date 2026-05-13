import { useState, useEffect, useCallback } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { RiLockLine, RiLoginBoxLine, RiLogoutBoxLine } from "@remixicon/react"
import { UserList } from "./components/UserList"
import { fetchUsers, saveUsers, testCredentials, fetchHost } from "./services/api"

type AuthState = {
  username: string
  password: string
  isAuthenticated: boolean
}

function App() {
  const [auth, setAuth] = useState<AuthState>(() => {
    const saved = localStorage.getItem("rachav_auth")
    if (saved) {
      try {
        const parsed = JSON.parse(saved)
        return { ...parsed, isAuthenticated: false }
      } catch {
        // ignore
      }
    }
    return { username: "", password: "", isAuthenticated: false }
  })
  
  const [users, setUsers] = useState<Record<string, string>>({})
  const [host, setHost] = useState<string>("")
  const [loading, setLoading] = useState(false)
  const [loginError, setLoginError] = useState("")

  const handleLogin = useCallback(async (username: string, password: string) => {
    setLoading(true)
    setLoginError("")
    
    try {
      const isValid = await testCredentials(username, password)
      if (isValid) {
        const newAuth = { username, password, isAuthenticated: true }
        setAuth(newAuth)
        localStorage.setItem("rachav_auth", JSON.stringify({ username, password }))
        
        const userList = await fetchUsers(username, password)
        setUsers(userList)
        const fetchedHost = await fetchHost(username, password)
        setHost(fetchedHost)
      } else {
        setLoginError("Invalid username or password")
      }
    } catch (error) {
      console.error("Login error:", error)
      setLoginError("Failed to connect to server. Please check if backend is running.")
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (auth.username && auth.password && !auth.isAuthenticated) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      handleLogin(auth.username, auth.password)
    }
  }, [auth.username, auth.password, auth.isAuthenticated, handleLogin])

  const handleLogout = () => {
    setAuth({ username: "", password: "", isAuthenticated: false })
    localStorage.removeItem("rachav_auth")
    setUsers({})
    setHost("")
  }

  const handleUsersChange = async (updatedUsers: Record<string, string>): Promise<boolean> => {
    if (!auth.isAuthenticated) return false
    
    try {
      await saveUsers(updatedUsers, auth.username, auth.password)
      setUsers(updatedUsers)
      return true
    } catch (error) {
      console.error("Failed to save users:", error)
      alert(`Failed to save users: ${error instanceof Error ? error.message : "Unknown error"}`)
      return false
    }
  }


  if (!auth.isAuthenticated) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="mx-auto w-16 h-16 bg-primary flex items-center justify-center mb-4">
              <RiLockLine className="h-8 w-8 text-primary-foreground" />
            </div>
            <CardTitle className="text-2xl">Rachav Dashboard</CardTitle>
            <CardDescription>
              Enter your credentials to access the user management system
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={(e) => {
              e.preventDefault()
              const formData = new FormData(e.currentTarget)
              const username = formData.get("username") as string
              const password = formData.get("password") as string
              handleLogin(username, password)
            }} className="space-y-4">
              {loginError && (
                <div className="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 text-sm">
                  {loginError}
                </div>
              )}
              
              <div className="space-y-2">
                <label className="text-sm font-medium">Username</label>
                <Input
                  name="username"
                  placeholder="Enter username"
                  defaultValue={auth.username}
                  disabled={loading}
                  required
                />
              </div>
              
              <div className="space-y-2">
                <label className="text-sm font-medium">Password</label>
                <Input
                  name="password"
                  type="password"
                  placeholder="Enter password"
                  defaultValue={auth.password}
                  disabled={loading}
                  required
                />
              </div>
              
              <Button type="submit" className="w-full" disabled={loading}>
                <RiLoginBoxLine className="h-4 w-4 mr-2" />
                {loading ? "Signing in..." : "Sign In"}
              </Button>
              
            </form>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <div className="flex items-center gap-2">
            <h1 className="text-xl font-heading font-semibold">rachav panel</h1>
          </div>
          <div className="flex items-center gap-4">
            <Button variant="outline" size="sm" onClick={handleLogout}>
              <RiLogoutBoxLine className="h-4 w-4 mr-2" />
              Logout
            </Button>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <Card>
          <CardHeader>
            <CardTitle>User List</CardTitle>
            <CardDescription>
              Manage users for basic authentication
            </CardDescription>
          </CardHeader>
          <CardContent>
            <UserList users={users} onUsersChange={handleUsersChange} host={host} />
          </CardContent>
        </Card>
      </main>
    </div>
  )
}

export default App
