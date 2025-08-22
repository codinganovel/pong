# Pong 📝

> **Leave a note. Check your notes. That's it.**

A proof-of-concept CLI tool for ephemeral messaging between GitHub users. Think digital sticky notes for developers.

## 🚨 Important Notice

**This is a proof-of-concept for educational purposes.** The code demonstrates the architecture and implementation of an ephemeral messaging system, but is not production-ready.

⚠️ **Production note:** If you want to actually use this project or build something similar, don’t hard-code a GitHub client secret in the CLI. CLIs are *public clients* and can’t keep secrets. Instead, use [GitHub’s Device Flow](https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#device-flow), which only requires a client ID and lets users log in by entering a one-time code in their browser. This keeps your app secure and avoids leaking secrets in source or binaries.

If you find this idea interesting and decide to deploy or build upon it, I'd appreciate credit for the original concept and implementation. Both the idea and code are freely shared here.

## The Idea

Pong fills the gap between formal communication (GitHub Issues, PRs) and real-time chat. It's for those moments when you want to say "coffee later?" or "nice work on that commit" without starting a whole conversation thread.

### Core Philosophy

- **One message per person** - No spam, no pressure
- **140 characters max** - Keep it light
- **Ephemeral by design** - Messages disappear when read
- **No replies, no threads** - Not a conversation starter
- **GitHub-native** - Built for developers who already live in GitHub

### Key Constraints (Features!)

- Messages are **immediately deleted** from the server when fetched
- **One message per sender-recipient pair** (sending a new one replaces the old unfetched one)
- **7-day auto-cleanup** for unfetched messages
- **No read receipts** - respect attention as a finite resource
- **Local history only** - you control your own data

## How It Works

```bash
# First time setup
pong login              # GitHub OAuth authentication

# Core usage  
pong send alice "coffee in 20?"    # Send a pong
pong                              # Check for new pongs (deletes them from server)
pong history                      # View your local history
pong clear-history               # Clear your local history
```

## Architecture

### Two-Component System

1. **CLI Client** - Handles authentication, sending/receiving, local history
2. **HTTP Server** - Minimal REST API with SQLite backend

### The Sticky Note Model

- **Server** = Office bulletin board (temporary message queue)
- **Client** = Personal notebook (your history lives with you)
- Messages are **immediately deleted** from server when fetched
- Zero permanent server storage

## Setup for Development/Testing

⚠️ **Configuration Required:** Before running, you'll need to set up:

1. **GitHub OAuth App:**
   - Create a GitHub OAuth app
   - Update `cmd/login.go` with your `ClientID` and `ClientSecret`
   - Update callback URI in `cmd/login.go`

2. **Server URL:**
   - Update `cmd/auth.go` with your server URL

### Build and Run

```bash
# Build CLI
go build -o pong .

# Build and run server
cd server
go build -o pong-server .
./pong-server

# Use CLI
./pong login
./pong send someuser "hello world"
./pong
```

## Files Structure

```
├── cmd/                 # CLI commands
│   ├── login.go        # GitHub OAuth flow
│   ├── send.go         # Send pongs
│   ├── root.go         # Check for pongs
│   ├── history.go      # Local history
│   └── auth.go         # Token management
├── server/
│   └── main.go         # HTTP server + SQLite
├── main.go             # CLI entry point
└── go.mod              # Dependencies
```

## Educational Value

This project demonstrates:

- **Ephemeral data patterns** - Immediate deletion as a feature
- **CLI design** with Cobra framework
- **GitHub OAuth integration** using official libraries
- **Constraint-driven development** - limitations as features
- **Client-server separation** - different data ownership models
- **SQLite for simple backends**
- **Go project structure** for CLI tools

## What This Isn't

- Not another chat app
- Not a notification system  
- Not for urgent communication
- Not production-ready (missing rate limiting, advanced error handling, etc.)

## Credit and Attribution

If you deploy this idea or build something similar, I'd appreciate a mention! Both the concept and implementation are shared freely here for the developer community.

Original idea and implementation by codinganovel.

## 📄 License

under ☕️, check out [the-coffee-license](https://github.com/codinganovel/The-Coffee-License)

I've included both licenses with the repo, do what you know is right. The licensing works by assuming you're operating under good faith.

---

*Pong: The sticky note of developer communication.*