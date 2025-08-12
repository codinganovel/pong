# Pong Implementation Plan

## The Vision
**"Leave a note. Check your notes. That's it."**

Digital sticky notes for developers. GitHub OAuth + CLI + ephemeral messaging. The missing piece in developer communication - too small for Issues, too casual for PRs, perfect for "hey, got a sec?"

## Core Architecture

### The Sticky Note Model
- **Server = Office bulletin board** (temporary message queue)
- **Client = Personal pocket** (your pong history lives with you)
- Messages get **immediately deleted** from server when fetched
- 30-day max retention for unfetched pongs
- **Zero permanent server storage**

### Key Constraints (Features!)
- **One message per person** (prevents spam) send a new one? replaces the old one if it hasn't been fetched or deleted.
- **140 characters max** (keeps it light)
- **No replies, no threads** (not a conversation starter)
- **No read receipts** (no pressure)
- **Ephemeral by design** (privacy-first)

## Technical Stack

### Client (Go)
```go
- github.com/spf13/cobra     // CLI framework
- github.com/cli/oauth       // GitHub's own OAuth library
- net/http                   // Basic HTTP client
```

### Server (Go + SQLite)
```sql
-- Minimal schema
CREATE TABLE users (id, github_username, github_id);
CREATE TABLE pongs (id, from_user_id, to_username, message, created_at);
```

### Three Endpoints
```
POST /pong   - Send a pong
GET  /pongs  - Get my pending pongs (immediately deletes them)
POST /clear  - cleans unfetched pongs after a month
```

## User Experience

### Commands
```bash
# First time setup
$ pong login    # GitHub OAuth in browser

# Core usage
$ pong @alice "coffee in 20?"           # Send
$ pong                                  # Check yours
$ pong clear                           # Clear local history

# That's it. Three commands.
```

### Flow
1. **Send**: Client validates GitHub user exists, posts to server
2. **Receive**: Client fetches pongs, server immediately deletes them
3. **History**: Lives locally, user controls retention

## Success Metrics

### Technical
- Single binary deployment
- <200 lines server code
- Zero database migrations needed
- Works offline for local history