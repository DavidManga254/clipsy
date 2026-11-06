# clipsy

A desktop clipboard manager that finds other devices on your local network and is meant to share clipboard items with them directly, with no cloud server in the middle.

It's a work in progress. Clipboard capture, local storage, peer discovery and the connection handshake work. Pushing clipboard items to a peer, and the encryption layer in `lib/crypto`, are not finished.

## What works today

- **Clipboard history.** A watcher saves what you copy into a local SQLite database.
- **Peer discovery.** Each instance broadcasts a `DISCOVER` message over UDP and listens for others. Peers are stored with their IP, TCP port, device name, status and last-seen time.
- **Handshake.** When a peer is found, the two devices open a TCP connection and exchange a handshake that carries their device IDs. The TCP port is picked at startup and sent in the discovery payload.
- **Device identity.** The first launch generates a device ID and keeps it in an `app_meta` table.
- **UI.** A React shell with a sidebar of devices, a header, a list of incoming requests and a device detail view.

## Stack

| Part | Tech |
|---|---|
| Backend | Go 1.24 |
| Desktop shell | [Wails v2](https://wails.io) |
| Storage | SQLite through `modernc.org/sqlite` (pure Go, no CGO) |
| Clipboard | `golang.design/x/clipboard`, plus a Windows-specific implementation |
| Frontend | React, TypeScript, Vite, Tailwind, TanStack Router |

## How it's laid out

```
cmd/clipsy/          entry point
app.go               startup: DB, device ID, UDP discovery, TCP listener
lib/clipboard/       clipboard watcher (Windows implementation included)
lib/peers/discovery/ UDP broadcast and listener
lib/peers/connect/   TCP listener and handshake
lib/DB/              SQLite setup, migrations, controllers
lib/crypto/          placeholder for the transport encryption
frontend/            the Wails UI
```

Discovery and the TCP listener each try a fixed list of ports (55000, 56000, 60000, 61000, 62000, 65000) and use the first free one, so two copies can run on one machine for testing.

## Running it

You need Go 1.24+, Node.js and the Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`).

```bash
wails dev      # live-reload development
wails build    # production build in build/bin
```

The database file `clipsy.db` is created on first run and is git-ignored. It holds your clipboard history, so it should never be committed.

## Where it's going

1. Send clipboard items to trusted peers over the TCP connection.
2. Encrypt that traffic.
3. A trust prompt, so a new device has to be approved before it can send anything.
4. Handle images and files, not only text.
5. macOS and Linux clipboard support.
