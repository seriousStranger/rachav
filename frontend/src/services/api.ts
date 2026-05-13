const API_BASE = "api"

const getAuthHeader = (username: string, password: string): string => {
  // Use encodeURIComponent to handle special characters
  try {
    const encoded =
      encodeURIComponent(username) + ":" + encodeURIComponent(password)
    const token = btoa(encoded)
    return `Basic ${token}`
  } catch (error) {
    console.error(
      "getAuthHeader error:",
      error,
      "username:",
      username,
      "password:",
      password
    )
    throw new Error(
      `Failed to create auth header: ${error instanceof Error ? error.message : "Unknown error"}`
    )
  }
}

export const testCredentials = async (
  username: string,
  password: string
): Promise<boolean> => {
  const response = await fetch(`${API_BASE}/user-list`, {
    headers: {
      "Api-Authorization": getAuthHeader(username, password),
    },
  })
  return response.ok
}

export const fetchUsers = async (
  username: string,
  password: string
): Promise<Record<string, string>> => {
  const response = await fetch(`${API_BASE}/user-list`, {
    headers: {
      "Api-Authorization": getAuthHeader(username, password),
    },
  })
  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`HTTP ${response.status}: ${errorText}`)
  }
  return response.json()
}

export const fetchHost = async (
  username: string,
  password: string
): Promise<string> => {
  const response = await fetch(`${API_BASE}/host`, {
    headers: {
      "Api-Authorization": getAuthHeader(username, password),
    },
  })
  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`HTTP ${response.status}: ${errorText}`)
  }
  const data = await response.json()
  return data.host
}

export const saveUsers = async (
  users: Record<string, string>,
  username: string,
  password: string
): Promise<boolean> => {
  try {
    const authHeader = getAuthHeader(username, password)
    const response = await fetch(`${API_BASE}/user-list`, {
      method: "POST",
      headers: {
        "Api-Authorization": authHeader,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(users),
    })
    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`HTTP ${response.status}: ${errorText}`)
    }
    return true
  } catch (error) {
    console.error("saveUsers error:", error)
    throw error
  }
}

