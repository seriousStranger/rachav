import { useState } from "react"
import type { ChangeEvent } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { RiEditLine, RiDeleteBinLine, RiSaveLine, RiCloseLine, RiUserLine, RiKey2Line } from "@remixicon/react"

interface UserListProps {
  users: Record<string, string>
  onUsersChange: (users: Record<string, string>) => Promise<boolean>
}

function generatePassword(length: number = 12): string {
  const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
  let password = ""
  for (let i = 0; i < length; i++) {
    const randomIndex = Math.floor(Math.random() * charset.length)
    password += charset[randomIndex]
  }
  return password
}

export function UserList({ users, onUsersChange }: UserListProps) {
  const [editingUser, setEditingUser] = useState<string | null>(null)
  const [newPassword, setNewPassword] = useState("")
  const [newUsername, setNewUsername] = useState("")
  const [newUserPassword, setNewUserPassword] = useState("")

  const handleDelete = async (username: string) => {
    if (!confirm(`Are you sure you want to delete user "${username}"?`)) {
      return
    }
    
    const updatedUsers = { ...users }
    delete updatedUsers[username]
    
    const success = await onUsersChange(updatedUsers)
    if (success) {
      // Success handled by parent
    }
  }

  const handleEdit = (username: string) => {
    setEditingUser(username)
    setNewPassword("")
  }

  const handleSaveEdit = async () => {
    if (!newPassword.trim() || !editingUser) {
      alert("Password cannot be empty")
      return
    }

    const updatedUsers = { ...users }
    updatedUsers[editingUser] = newPassword
    
    const success = await onUsersChange(updatedUsers)
    if (success) {
      setEditingUser(null)
      setNewPassword("")
    }
  }

  const handleCancelEdit = () => {
    setEditingUser(null)
    setNewPassword("")
  }

  const handleAddUser = async () => {
    if (!newUsername.trim() || !newUserPassword.trim()) {
      alert("Username and password cannot be empty")
      return
    }

    if (users[newUsername]) {
      alert("User already exists")
      return
    }

    const updatedUsers = { ...users, [newUsername]: newUserPassword }
    
    const success = await onUsersChange(updatedUsers)
    if (success) {
      setNewUsername("")
      setNewUserPassword("")
    }
  }

  const userEntries = Object.entries(users)

  return (
    <div className="space-y-6">
      {/* Add User Form */}
      <div className="border rounded-lg p-4">
        <h3 className="text-lg font-medium mb-4 flex items-center gap-2">
          <RiUserLine className="h-5 w-5" />
          Add New User
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Username</label>
            <Input
              placeholder="Enter username"
              value={newUsername}
              onChange={(e: ChangeEvent<HTMLInputElement>) => setNewUsername(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Password</label>
            <div className="flex gap-2">
              <Input
                placeholder="Enter password or generate"
                value={newUserPassword}
                onChange={(e: ChangeEvent<HTMLInputElement>) => setNewUserPassword(e.target.value)}
                className="flex-1"
              />
              <Button
                type="button"
                variant="outline"
                onClick={() => setNewUserPassword(generatePassword())}
              >
                <RiKey2Line className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
        <Button className="mt-4" onClick={handleAddUser}>
          Add User
        </Button>
      </div>

      {/* Users Table */}
      <div>
        <h3 className="text-lg font-medium mb-4">User List ({userEntries.length} total)</h3>
        {userEntries.length === 0 ? (
          <div className="text-center py-12 border-2 border-dashed border-border rounded-lg">
            <div className="mx-auto w-12 h-12 text-muted-foreground mb-4">
              <RiUserLine className="w-full h-full" />
            </div>
            <h3 className="text-lg font-medium mb-2">No users configured</h3>
            <p className="text-muted-foreground">Add users using the form above.</p>
          </div>
        ) : (
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Username</TableHead>
                  <TableHead>Password</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {userEntries.map(([username, password]) => (
                  <TableRow key={username}>
                    <TableCell className="font-medium">{username}</TableCell>
                    <TableCell>
                      {editingUser === username ? (
                        <div className="flex items-center gap-1 max-w-md">
                          <Input
                            placeholder="New password"
                            value={newPassword}
                            onChange={(e: ChangeEvent<HTMLInputElement>) => setNewPassword(e.target.value)}
                            className="flex-1 min-w-0"
                          />
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            className="shrink-0 h-9 w-9 p-0"
                            onClick={() => setNewPassword(generatePassword())}
                            title="Generate password"
                          >
                            <RiKey2Line className="h-4 w-4" />
                          </Button>
                        </div>
                      ) : (
                        <span className="font-mono break-all">{password}</span>
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      {editingUser === username ? (
                        <div className="flex justify-end gap-2">
                          <Button size="sm" variant="outline" onClick={handleSaveEdit}>
                            <RiSaveLine className="h-4 w-4 mr-1" />
                            Save
                          </Button>
                          <Button size="sm" variant="ghost" onClick={handleCancelEdit}>
                            <RiCloseLine className="h-4 w-4 mr-1" />
                            Cancel
                          </Button>
                        </div>
                      ) : (
                        <div className="flex justify-end gap-2">
                          <Button size="sm" variant="outline" onClick={() => handleEdit(username)}>
                            <RiEditLine className="h-4 w-4 mr-1" />
                            Edit
                          </Button>
                          <Button size="sm" variant="destructive" onClick={() => handleDelete(username)}>
                            <RiDeleteBinLine className="h-4 w-4 mr-1" />
                            Delete
                          </Button>
                        </div>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </div>
    </div>
  )
}