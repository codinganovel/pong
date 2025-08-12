# Pong System Behavior Documentation

## Core Commands

### `pong login`
**What it does:** Authenticates with GitHub using OAuth device flow
- Shows a code to enter at github.com/login/device
- Saves token to `~/.pong/token` 
- Required before sending or receiving pongs

### `pong send username "message"`
**What it does:** Sends a pong to a GitHub user
- Validates that target GitHub username exists
- Enforces 140 character message limit
- Replaces any existing unsent pong from you to that user (one-message-per-recipient rule)
- Requires valid login token

**Example:** `pong send alice "coffee later?"`

### `pong` (check for new pongs)
**What it does:** Fetches and displays your pending pongs
- **Server behavior:** Immediately deletes pongs after sending them to you (ephemeral)
- **Client behavior:** Saves fetched pongs to local history before displaying
- Shows count and lists all new pongs
- If no pongs waiting: displays "No pongs waiting for you!"

**Important:** Only shows **NEW** pongs from server, not your local history

### `pong history`
**What it does:** Shows your complete local pong history
- Displays all pongs you've ever received (that were saved locally)
- Shows most recent first
- Includes timestamp of when you fetched each pong
- Purely local - doesn't contact server

### `pong clear-history`
**What it does:** Deletes your local history file
- Removes `~/.pong/history.json`
- Does not affect server or other users
- Cannot be undone

**Example output:** "✓ History cleared!" or "No history file to clear."

## System Constraints & Rules

### Message Constraints
- **140 character limit** per message
- **One message per sender-recipient pair** 
  - If Alice sends Bob a pong, then sends another, the first is replaced
  - Only the most recent unsent pong is kept

### Ephemeral Server Behavior
- **Immediate deletion:** Pongs deleted from server the moment you fetch them
- **7-day cleanup:** Unfetched pongs automatically deleted after 7 days
- **No server history:** Server never retains pongs long-term

### Local History Behavior
- **Automatic saving:** Every time you run `pong`, new pongs are saved locally before display
- **Persistent storage:** Lives in `~/.pong/history.json` until you delete it
- **User controlled:** Only you can view or delete your history

## Data Flow Examples

### Sending a Pong
```
You: pong send alice "meeting at 3?"
├── CLI validates "alice" exists on GitHub ✓
├── CLI sends to server
├── Server validates your token ✓
├── Server deletes any existing pong from you to alice
├── Server stores new pong
└── Returns success ✓
```

### Receiving Pongs
```
You: pong
├── CLI fetches from server
├── Server finds 2 pongs for you
├── Server sends pongs to CLI
├── Server immediately deletes those 2 pongs (ephemeral!)
├── CLI saves pongs to ~/.pong/history.json
└── CLI displays: "You have 2 pongs: ..."
```

### Viewing History
```
You: pong history
├── CLI reads ~/.pong/history.json
├── Shows all pongs you've ever received
└── No server contact needed
```

## Automatic Cleanup

### Server Cleanup
- **Frequency:** Every 24 hours automatically
- **Targets:** Unfetched pongs older than 7 days
- **Reasoning:** Prevents abandoned pongs from accumulating
- **User impact:** None (only affects pongs you never checked)

### Manual Server Cleanup
- **Endpoint:** `POST /clear` (for admin use)
- **Same logic:** Removes unfetched pongs older than 7 days

## File Locations

### Client Files
- **Auth token:** `~/.pong/token` (created by `pong login`)
- **History:** `~/.pong/history.json` (created automatically when receiving pongs)

### Server Files
- **Database:** `pongs.db` (SQLite file in server directory)
- **Schema:** Simple table with from_username, to_username, message, created_at

## Fetch vs History Distinction

**Key Behavioral Difference:**

- `pong` = "Check for NEW pongs" (contacts server, shows only unread messages)
- `pong history` = "Show ALL my pongs" (local only, shows everything I've ever received)

**Example Timeline:**
1. Day 1: Alice sends you "hello"
2. Day 1: You run `pong` → Shows "hello", saves to history, server deletes it
3. Day 2: Bob sends you "sup?"  
4. Day 2: You run `pong` → Shows only "sup?" (new), saves to history
5. Day 2: You run `pong history` → Shows both "hello" and "sup?" with timestamps

## Error Handling

### Authentication Errors
- Invalid/expired token → "Error: not logged in. Run 'pong login' first"
- GitHub API issues → Specific error message with HTTP status

### Validation Errors
- Username doesn't exist → "Error: GitHub user 'username' not found"
- Message too long → "Message too long (X chars). Max 140 characters."

### Server Errors
- Server unreachable → "Failed to send pong: connection refused"
- Server errors → Displays server error message

### History Errors
- Cannot save history → Warning message (non-fatal)
- Cannot load history → Error message and exit

## Privacy & Data Ownership

### What's Ephemeral (Server)
- All pongs are deleted immediately when fetched
- Unfetched pongs deleted after 7 days
- No permanent server-side user data storage

### What's Persistent (Local)
- Your authentication token (until you delete `~/.pong/token`)
- Your pong history (until you run `pong clear-history`)
- Your choice when to clear local data

### What's Never Stored
- Other users' histories
- Read receipts or timestamps of when others read your pongs
- Any metadata beyond basic message content