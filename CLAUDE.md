# massive-cli

A comprehensive Go CLI for the [Massive](https://massive.com) financial data API covering stocks, crypto, forex, futures, indices, and options with REST, WebSocket streaming, and S3 flat file support.

## Quick Reference

- **Language:** Go 1.24.1
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
- **Build:** `go build -o massive .`
- **Test:** `go test ./...`
- **Config file:** `~/.config/massive/config.json`
- **Base API URL:** `https://api.massive.com`
- **WebSocket URLs:** `wss://delayed.massive.com` (default), `wss://socket.massive.com` (real-time)
- **S3 Endpoint:** `https://files.massive.com`

## Project Structure

```
massive-cli/
├── main.go                     # Entry point
├── cmd/                        # Cobra commands (45 files)
│   ├── root.go                 # Root command, --output flag (table/json)
│   ├── config.go               # Config init/show subcommands
│   ├── helpers.go              # newClient(), maskString(), printJSON()
│   ├── stocks*.go              # 14 files: bars, snapshots, fundamentals, etc.
│   ├── crypto.go               # Crypto REST commands
│   ├── forex.go                # Forex REST commands
│   ├── futures.go              # Futures REST commands
│   ├── indices*.go             # 5 files: bars, snapshots, tickers, etc.
│   ├── options*.go             # 5 files: contracts, snapshots, trades, etc.
│   ├── ws_stocks.go            # WS core: connectAndStreamAsset, formatters, helpers
│   ├── ws_options.go           # WS options commands
│   ├── ws_indices.go           # WS indices commands
│   ├── ws_crypto.go            # WS crypto commands
│   ├── ws_forex.go             # WS forex commands
│   ├── ws_futures.go           # WS futures commands
│   ├── flatfiles.go            # S3 flat file list/download
│   ├── benzinga.go             # Benzinga partner data
│   ├── economy.go              # Economic indicators
│   ├── etfglobal.go            # ETF analytics
│   └── tmx.go                  # TMX Canadian market data
├── internal/
│   ├── api/                    # REST API client (56 files)
│   │   ├── client.go           # HTTP client, apiKey query param auth
│   │   ├── stocks.go           # Stock API methods
│   │   ├── crypto.go           # Crypto API methods
│   │   ├── forex.go            # Forex API methods
│   │   ├── futures.go          # Futures API methods
│   │   ├── indices_*.go        # Index API methods (5 files)
│   │   ├── options_*.go        # Options API methods (5 files)
│   │   ├── benzinga.go         # Benzinga API methods
│   │   ├── economy.go          # Economy API methods
│   │   ├── etfglobal.go        # ETF Global API methods
│   │   ├── tmx.go              # TMX API methods
│   │   └── *_test.go           # One test file per API file
│   ├── config/                 # Config load/save (~/.config/massive/config.json)
│   │   ├── config.go
│   │   └── config_test.go
│   ├── ws/                     # WebSocket client library
│   │   ├── client.go
│   │   └── client_test.go
│   └── flatfiles/              # S3 flat file client
│       ├── client.go
│       └── client_test.go
```

## Configuration

**Environment variables** (take priority over config file):
- `MASSIVE_API_KEY` - API key for REST and WebSocket auth
- `MASSIVE_S3_ACCESS_KEY` - S3 access key for flat files
- `MASSIVE_S3_SECRET_KEY` - S3 secret key for flat files

**Config struct** (`internal/config/config.go`):
```go
type Config struct {
    APIKey      string
    BaseURL     string // default: https://api.massive.com
    S3AccessKey string
    S3SecretKey string
    S3Endpoint  string // default: https://files.massive.com
}
```

## Architecture Patterns

### REST API Client
- Base client in `internal/api/client.go` with 30s HTTP timeout
- Auth via `?apiKey=` query parameter on every request
- All methods return typed response structs
- `SetBaseURL()` for test overrides
- Method naming: `Get{AssetClass}{Operation}()` (e.g., `GetStocksBars()`)
- Parameter structs with optional fields for query params

### Cobra Command Pattern
- Parent commands group by asset class (e.g., `stocks`, `crypto`)
- Child commands for specific operations (e.g., `stocks bars`, `stocks snapshots ticker`)
- Persistent flag `--output` on root (table or json, default table)
- Table output uses `text/tabwriter`
- JSON output uses `json.MarshalIndent` with 2-space indent

### WebSocket Streaming
- All WS commands live in `cmd/ws_*.go` files
- Core connection logic is in `cmd/ws_stocks.go` via `connectAndStreamAsset()`
- Auth flow: connect -> read "connected" -> send `{"action":"auth","params":"API_KEY"}` -> read auth_success -> send `{"action":"subscribe","params":"CHANNEL.TICKER"}` -> stream
- URL format: `wss://{host}/{assetClass}` (no channel in URL, no apiKey in query)
- `--realtime` persistent flag switches from delayed to real-time endpoint
- `--all` flag on each subcommand subscribes to wildcard (`CHANNEL.*`)
- Context cancellation closes the connection via goroutine (not read deadlines)
- Each asset class has its own formatter functions and `buildXxxSubscriptionParams()` helper

### WebSocket Channels by Asset Class
| Asset    | Channels                        |
|----------|---------------------------------|
| Stocks   | T, Q, AM, A, LULD, FMV         |
| Options  | T, Q, AM, A, FMV               |
| Indices  | AM, A, V                        |
| Crypto   | XT, XQ, XA, XAS, FMV           |
| Forex    | C, CA, CAS, FMV                |
| Futures  | T, Q, AM, A                     |

### Flat Files (S3)
- S3-compatible storage at `https://files.massive.com`, bucket `flatfiles`
- Asset prefixes: `us_stocks_sip`, `us_options_opra`, `us_indices`, `global_crypto`, `global_forex`
- Data types: `trades_v1`, `quotes_v1`, `day_aggs_v1`, `minute_aggs_v1`
- Commands: `files list`, `files download`, `files assets`, `files types`

### Testing Pattern
- One test file per API module (e.g., `stocks_test.go`)
- Mock HTTP server with route-based responses
- `NewTestClient()` helper creates client pointing at mock server
- Test data defined as JSON string constants
- Config tests use `SetConfigDir()` with `t.TempDir()`

## Key Files for Common Tasks

| Task | Files |
|------|-------|
| Add new REST endpoint | `internal/api/{asset}.go`, `internal/api/{asset}_test.go`, `cmd/{asset}.go` |
| Add new WS channel | `cmd/ws_{asset}.go`, update `printTableHeader()` in `cmd/ws_stocks.go` |
| Add new asset class | New files in both `internal/api/` and `cmd/`, register in `init()` |
| Change auth logic | `internal/api/client.go` (REST), `cmd/ws_stocks.go` (WebSocket) |
| Change config | `internal/config/config.go` |
| Add flat file asset | `internal/flatfiles/client.go` constants |

## Command Tree Overview

```
massive
├── config [init|show]
├── stocks [bars|open-close|market|snapshots|quotes|trades|news|tickers|
│           exchanges|fundamentals|corporate-actions|filings|indicators|market-ops]
├── crypto [bars|previous-day-bar|daily-market-summary|daily-ticker-summary|
│           snapshots|unified-snapshot|tickers|ticker-overview|trades|last-trade|
│           conditions|exchanges|market-holidays|market-status|indicators|quotes]
├── forex  [bars|previous-day-bar|daily-market-summary|convert|quotes|last-quote|
│           snapshots|unified-snapshot|tickers|ticker-overview|exchanges|
│           market-holidays|market-status|indicators]
├── futures [bars|contracts|products|schedules|exchanges|snapshot|trades|quotes]
├── indices [bars|previous-day-bar|daily-ticker-summary|snapshots|tickers|
│            market-holidays|market-status|indicators]
├── options [bars|contracts|snapshots|previous-day-bar|daily-ticker-summary|
│            trades|quotes|last-trade|last-quote|market-holidays|market-status|indicators]
├── ws
│   ├── stocks  [trades|quotes|agg-minute|agg-second|luld|fmv]
│   ├── options [trades|quotes|agg-minute|agg-second|fmv]
│   ├── indices [agg-minute|agg-second|value]
│   ├── crypto  [trades|quotes|agg-minute|agg-second|fmv]
│   ├── forex   [quotes|agg-minute|agg-second|fmv]
│   └── futures [trades|quotes|agg-minute|agg-second]
├── files [list|download|assets|types]
├── benzinga [news|ratings|earnings|guidance|analysts]
├── economy [inflation|labor-market|treasury-yields]
├── etf-global [analytics|constituents]
└── tmx [corporate-events]
```

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/gorilla/websocket` - WebSocket client
- `github.com/aws/aws-sdk-go-v2` - S3 flat file access
- `github.com/joho/godotenv` - .env file loading

## Code Conventions

- Every new file gets a copyright header with current date
- Detailed comments above every function (public and private)
- One test file per source file
- Table output via `text/tabwriter`, JSON via `json.MarshalIndent`
- Error wrapping with `fmt.Errorf("context: %w", err)`
- Config priority: env vars > config file > defaults
