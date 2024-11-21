# **Windows Startup Registry Management Library**

## **Overview**
The `winstartupreg` library provides a straightforward interface to manage Windows startup registry entries. It allows applications to add, remove, and list programs configured to run automatically at Windows startup. This library supports multiple registry locations, including both user-specific and machine-wide scopes.

### **Features**
- Add startup entries to different registry locations.
- Remove startup entries safely from all known locations.
- List startup entries for specific or all registry locations.
- Comprehensive error handling and input validation.

---

## **Installation**
To use the `winstartupreg` package, ensure you have Go installed and set up on your system. 

Install the library:
```bash
go get github.com/nishansanjuka/winstartupreg
```

Import the package:
```go
import "github.com/nishansanjuka/winstartupreg"
```

---

## **API Documentation**

### **Types**

#### **`StartupRegistryType`**
Enumeration defining different registry locations for startup entries:
- `CurrentUserRun`: Current user’s startup entries.
- `CurrentUserRunOnce`: Current user’s one-time startup entries.
- `AllUsersRun`: All users’ startup entries.
- `AllUsersRunOnce`: All users’ one-time startup entries.

#### **`StartupEntry`**
Structure representing a Windows startup registry entry:
```go
type StartupEntry struct {
    Name    string // The name of the startup entry
    Command string // The executable command to run at startup
}
```

---

### **Functions**

#### **`AddStartupEntry`**
Adds an application to a specified Windows startup registry location.

**Signature:**
```go
func AddStartupEntry(entry StartupEntry, registryType StartupRegistryType) error
```

**Parameters:**
- `entry` (StartupEntry): The startup entry to add.
- `registryType` (StartupRegistryType): The target registry location.

**Returns:**
- `error`: Describes any failure, or `nil` on success.

**Usage Example:**
```go
entry := winstartupreg.StartupEntry{
    Name:    "MyApp",
    Command: "C:\\path\\to\\MyApp.exe",
}
err := winstartupreg.AddStartupEntry(entry, winstartupreg.CurrentUserRun)
if err != nil {
    fmt.Println("Error adding startup entry:", err)
}
```

---

#### **`RemoveStartupEntry`**
Removes a startup entry from a specific registry location.

**Signature:**
```go
func RemoveStartupEntry(entryName string, registryType StartupRegistryType) error
```

**Parameters:**
- `entryName` (string): The name of the startup entry to remove.
- `registryType` (StartupRegistryType): The registry location to target.

**Returns:**
- `error`: Describes any failure, or `nil` on success.

**Usage Example:**
```go
err := winstartupreg.RemoveStartupEntry("MyApp", winstartupreg.CurrentUserRun)
if err != nil {
    fmt.Println("Error removing startup entry:", err)
}
```

---

#### **`SafeRemoveStartupEntry`**
Attempts to remove a startup entry from all known registry locations.

**Signature:**
```go
func SafeRemoveStartupEntry(entryName string) error
```

**Parameters:**
- `entryName` (string): The name of the startup entry to remove.

**Returns:**
- `error`: Describes any failure, or `nil` if removed from at least one location.

**Usage Example:**
```go
err := winstartupreg.SafeRemoveStartupEntry("MyApp")
if err != nil {
    fmt.Println("Error safely removing startup entry:", err)
}
```

---

#### **`ListStartupEntries`**
Retrieves all startup entries from a specific registry location.

**Signature:**
```go
func ListStartupEntries(registryType StartupRegistryType) (map[string]string, error)
```

**Parameters:**
- `registryType` (StartupRegistryType): The target registry location.

**Returns:**
- `map[string]string`: A map of startup entry names to commands.
- `error`: Describes any failure, or `nil` on success.

**Usage Example:**
```go
entries, err := winstartupreg.ListStartupEntries(winstartupreg.CurrentUserRun)
if err != nil {
    fmt.Println("Error listing startup entries:", err)
} else {
    for name, command := range entries {
        fmt.Printf("Name: %s, Command: %s\n", name, command)
    }
}
```

---

#### **`ListAllStartupEntries`**
Retrieves all startup entries from all known registry locations.

**Signature:**
```go
func ListAllStartupEntries() (map[StartupRegistryType]map[string]string, error)
```

**Parameters:**
- None.

**Returns:**
- `map[StartupRegistryType]map[string]string`: A nested map containing startup entries for each registry location.
- `error`: Describes any failure, or `nil` on success.

**Usage Example:**
```go
allEntries, err := winstartupreg.ListAllStartupEntries()
if err != nil {
    fmt.Println("Error listing all startup entries:", err)
} else {
    for regType, entries := range allEntries {
        fmt.Printf("Registry Type: %v\n", regType)
        for name, command := range entries {
            fmt.Printf("  Name: %s, Command: %s\n", name, command)
        }
    }
}
```

---

### **Testing**
The library includes comprehensive unit tests using [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/).

#### **Running Tests**
1. Ensure you have the Ginkgo CLI installed:
   ```bash
   go install github.com/onsi/ginkgo/v2/ginkgo
   ```
2. Run the tests:
   ```bash
   ginkgo ./...
   ```

---

### **Error Handling**
The library uses detailed error messages to indicate:
- Missing or invalid entry names.
- Non-existent executable paths.
- Registry access issues.

### **Best Practices**
- Always use absolute paths for the `Command` field in `StartupEntry`.
- Verify that the executable exists before adding it to the registry.

---

### **License**
This library is open source and distributed under the MIT License.

For contributions and issues, visit the [GitHub Repository](https://github.com/nishansanjuka/winstartupreg).