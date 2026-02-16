# Massive CLI

A command-line interface for the [Massive](https://massive.com) financial data API. Access real-time and historical market data for stocks, options, indices, crypto, forex, and futures -- all from your terminal.

## Why a CLI?

There are plenty of ways to interact with financial data APIs. You could write scripts, use an SDK, or wire up an MCP server for your AI assistant. We built Massive CLI because we think a good CLI tool hits a sweet spot that other approaches miss:

**A context-friendly alternative to MCP servers.** MCP servers are great, but they can flood your AI agent's context window with massive JSON payloads, eating up tokens and degrading response quality. A CLI tool lets you (or your agent) run targeted queries and pipe just the relevant output into the conversation. Smaller, more focused context means better results.

**Better for debugging when agents get it wrong.** When an AI agent calls an API through an MCP server and something breaks, you're left staring at opaque error messages buried in tool call logs. With a CLI, you can run the exact same command yourself, see the raw output, tweak flags, and figure out what went wrong in seconds. It's the fastest path from "that's not right" to "here's what actually happened."

**Works everywhere, no integration required.** No SDK to install, no server to run, no configuration files to wire into your editor. Just a single binary. It works in your terminal, in shell scripts, in CI pipelines, and as a tool that any AI coding assistant can shell out to. Any agent that can run a bash command can use Massive CLI.

**Human-readable by default, machine-readable when you need it.** Every command outputs clean, aligned tables for scanning by eye, and switches to JSON with a single `--output json` flag for piping into `jq`, scripts, or AI context.

## Installation

### From Source

```bash
git clone https://github.com/cloudmanic/massive.git
cd massive
go build -o massive .
```

Move the binary somewhere in your `$PATH`:

```bash
mv massive /usr/local/bin/
```

### Requirements

- Go 1.24+ (for building from source)
- A [Massive](https://massive.com) API key

## Quick Start

### 1. Configure your API key

```bash
massive config init
```

This prompts for your API key and saves it to `~/.config/massive/config.json`. You can also set it via environment variable:

```bash
export MASSIVE_API_KEY=your_api_key_here
```

### 2. Start querying

```bash
# Get daily bars for Apple
massive stocks bars AAPL --from 2025-01-01 --to 2025-01-31

# Stream real-time stock aggregates
massive ws stocks agg-minute --all

# Get the latest news
massive stocks news --ticker AAPL --limit 5

# Everything outputs JSON too
massive stocks bars AAPL --from 2025-01-01 --to 2025-01-31 --output json
```

## Configuration

Massive CLI looks for credentials in this order:

1. **Environment variables** (highest priority)
2. **Config file** at `~/.config/massive/config.json`

### Environment Variables

| Variable | Description |
|----------|-------------|
| `MASSIVE_API_KEY` | API key for REST and WebSocket authentication |
| `MASSIVE_S3_ACCESS_KEY` | S3 access key for flat file downloads |
| `MASSIVE_S3_SECRET_KEY` | S3 secret key for flat file downloads |

You can also put these in a `.env` file in your working directory. See `.env.example` for the template.

### Config Commands

```bash
# Interactive setup (API key + optional S3 credentials)
massive config init

# Show current config (keys are masked)
massive config show
```

## Output Formats

Every command supports two output formats via the `--output` (`-o`) flag:

```bash
# Table output (default) -- human-readable, aligned columns
massive stocks bars AAPL --from 2025-01-01 --to 2025-01-31

# JSON output -- machine-readable, pipe to jq or feed to an AI agent
massive stocks bars AAPL --from 2025-01-01 --to 2025-01-31 -o json
```

## Commands

### Stocks

```bash
# OHLC aggregate bars with configurable timespan
massive stocks bars AAPL --from 2025-01-01 --to 2025-01-31
massive stocks bars AAPL --from 2025-01-01 --to 2025-01-31 --timespan week --multiplier 1

# Daily open/close
massive stocks open-close AAPL --date 2025-01-15

# Market summary
massive stocks market --date 2025-01-15

# Snapshots
massive stocks snapshots ticker AAPL
massive stocks snapshots all
massive stocks snapshots gainers
massive stocks snapshots losers

# Quotes and trades
massive stocks quotes AAPL --date 2025-01-15
massive stocks trades AAPL --date 2025-01-15

# News
massive stocks news --ticker AAPL --limit 10
massive stocks news --published-from 2025-01-01 --published-to 2025-01-31

# Reference data
massive stocks tickers --search apple
massive stocks exchanges

# Fundamentals
massive stocks fundamentals short-interest --ticker AAPL
massive stocks fundamentals short-volume --ticker AAPL
massive stocks fundamentals float --ticker AAPL
massive stocks fundamentals balance-sheet AAPL
massive stocks fundamentals income-statement AAPL
massive stocks fundamentals cash-flow AAPL
massive stocks fundamentals ratios AAPL

# Corporate actions
massive stocks corporate-actions dividends AAPL
massive stocks corporate-actions splits AAPL

# SEC filings
massive stocks filings sections AAPL
massive stocks filings risk-factors AAPL
massive stocks filings risk-categories AAPL

# Technical indicators
massive stocks indicators sma AAPL --from 2025-01-01 --to 2025-01-31
massive stocks indicators ema AAPL --from 2025-01-01 --to 2025-01-31
massive stocks indicators rsi AAPL --from 2025-01-01 --to 2025-01-31
massive stocks indicators macd AAPL --from 2025-01-01 --to 2025-01-31

# Market operations
massive stocks market-ops holidays
massive stocks market-ops status
```

### Options

```bash
# Aggregate bars
massive options bars O:SPY241220P00720000 --from 2024-12-01 --to 2024-12-20

# Contracts
massive options contracts list --underlying AAPL
massive options contracts get O:SPY241220P00720000

# Snapshots
massive options snapshots contract O:SPY241220P00720000
massive options snapshots chain SPY

# Previous day and daily summary
massive options previous-day-bar O:SPY241220P00720000
massive options daily-ticker-summary O:SPY241220P00720000 --date 2024-12-15

# Trades and quotes
massive options trades O:SPY241220P00720000
massive options quotes O:SPY241220P00720000
massive options last-trade O:SPY241220P00720000
massive options last-quote O:SPY241220P00720000

# Technical indicators
massive options indicators sma O:SPY241220P00720000 --from 2024-12-01 --to 2024-12-20
massive options indicators ema O:SPY241220P00720000 --from 2024-12-01 --to 2024-12-20
massive options indicators rsi O:SPY241220P00720000 --from 2024-12-01 --to 2024-12-20
massive options indicators macd O:SPY241220P00720000 --from 2024-12-01 --to 2024-12-20

# Market operations
massive options market-holidays
massive options market-status
```

### Indices

```bash
# Aggregate bars
massive indices bars I:SPX --from 2025-01-01 --to 2025-01-31

# Previous day and daily summary
massive indices previous-day-bar I:SPX
massive indices daily-ticker-summary I:SPX --date 2025-01-15

# Snapshots
massive indices snapshots ticker I:SPX
massive indices snapshots all

# Reference data
massive indices tickers

# Technical indicators
massive indices indicators sma I:SPX --from 2025-01-01 --to 2025-01-31
massive indices indicators ema I:SPX --from 2025-01-01 --to 2025-01-31
massive indices indicators rsi I:SPX --from 2025-01-01 --to 2025-01-31
massive indices indicators macd I:SPX --from 2025-01-01 --to 2025-01-31

# Market operations
massive indices market-holidays
massive indices market-status
```

### Crypto

```bash
# Aggregate bars
massive crypto bars X:BTC-USD --from 2025-01-01 --to 2025-01-31

# Previous day bar
massive crypto previous-day-bar X:BTC-USD

# Market summaries
massive crypto daily-market-summary --date 2025-01-15
massive crypto daily-ticker-summary X:BTC-USD --date 2025-01-15

# Snapshots
massive crypto snapshots market
massive crypto snapshots ticker X:BTC-USD
massive crypto snapshots gainers
massive crypto snapshots losers
massive crypto unified-snapshot X:BTC-USD

# Trades
massive crypto trades X:BTC-USD
massive crypto last-trade BTC USD

# Reference data
massive crypto tickers
massive crypto ticker-overview X:BTC-USD
massive crypto conditions
massive crypto exchanges

# Technical indicators
massive crypto indicators sma X:BTC-USD --from 2025-01-01 --to 2025-01-31
massive crypto indicators ema X:BTC-USD --from 2025-01-01 --to 2025-01-31
massive crypto indicators rsi X:BTC-USD --from 2025-01-01 --to 2025-01-31
massive crypto indicators macd X:BTC-USD --from 2025-01-01 --to 2025-01-31

# Market operations
massive crypto market-holidays
massive crypto market-status
```

### Forex

```bash
# Aggregate bars
massive forex bars C:EURUSD --from 2025-01-01 --to 2025-01-31

# Previous day bar
massive forex previous-day-bar C:EURUSD

# Market summary
massive forex daily-market-summary --date 2025-01-15

# Currency conversion
massive forex convert EUR USD --amount 1000

# Quotes
massive forex quotes C:EURUSD
massive forex last-quote EUR USD

# Snapshots
massive forex snapshots market
massive forex snapshots ticker C:EURUSD
massive forex snapshots gainers
massive forex snapshots losers
massive forex unified-snapshot C:EURUSD

# Reference data
massive forex tickers
massive forex ticker-overview C:EURUSD
massive forex exchanges

# Technical indicators
massive forex indicators sma C:EURUSD --from 2025-01-01 --to 2025-01-31
massive forex indicators ema C:EURUSD --from 2025-01-01 --to 2025-01-31
massive forex indicators rsi C:EURUSD --from 2025-01-01 --to 2025-01-31
massive forex indicators macd C:EURUSD --from 2025-01-01 --to 2025-01-31

# Market operations
massive forex market-holidays
massive forex market-status
```

### Futures

```bash
# Aggregate bars
massive futures bars ESZ4 --from 2025-01-01 --to 2025-01-31

# Reference data
massive futures contracts
massive futures products
massive futures schedules
massive futures exchanges

# Snapshots, trades, and quotes
massive futures snapshot ESZ4
massive futures trades ESZ4
massive futures quotes ESZ4
```

### WebSocket Streaming

Stream real-time (or 15-minute delayed) market data directly to your terminal. All WebSocket commands support `--all` to subscribe to every available ticker, or you can specify individual tickers as arguments.

```bash
# Stocks
massive ws stocks trades AAPL MSFT
massive ws stocks quotes AAPL
massive ws stocks agg-minute --all
massive ws stocks agg-second AAPL
massive ws stocks luld --all
massive ws stocks fmv AAPL

# Options
massive ws options trades O:SPY241220P00720000
massive ws options quotes O:SPY241220P00720000
massive ws options agg-minute --all
massive ws options agg-second --all
massive ws options fmv --all

# Indices
massive ws indices agg-minute I:SPX I:DJI
massive ws indices agg-second --all
massive ws indices value I:SPX I:COMP

# Crypto
massive ws crypto trades X:BTC-USD X:ETH-USD
massive ws crypto quotes X:BTC-USD
massive ws crypto agg-minute --all
massive ws crypto agg-second --all
massive ws crypto fmv X:BTC-USD

# Forex
massive ws forex quotes C:EURUSD C:GBPUSD
massive ws forex agg-minute --all
massive ws forex agg-second --all
massive ws forex fmv --all

# Futures
massive ws futures trades ESZ4
massive ws futures quotes ESZ4
massive ws futures agg-minute --all
massive ws futures agg-second --all
```

By default, WebSocket commands connect to the **delayed** (15-minute) endpoint. Add `--realtime` to connect to the real-time feed (requires a real-time data subscription):

```bash
massive ws stocks agg-minute --all --realtime
```

Press `Ctrl+C` to disconnect gracefully.

### Flat Files (S3)

Download bulk historical data as gzipped CSV files from Massive's S3-compatible storage.

```bash
# List available asset classes and data types
massive files assets
massive files types

# List files for a specific asset class, data type, and year
massive files list stocks trades --year 2024
massive files list stocks trades --year 2024 --month 06

# Download a specific file
massive files download stocks trades 2024-06-15
massive files download stocks trades 2024-06-15 --output-dir ./data
```

**Available asset classes:** `stocks`, `options`, `indices`, `crypto`, `forex`

**Available data types:** `trades`, `quotes`, `day-aggs`, `minute-aggs`

S3 access requires separate credentials. Set them via environment variables or `massive config init`.

### Benzinga (Partner Data)

```bash
massive benzinga news --tickers AAPL --limit 10
massive benzinga ratings --ticker AAPL
massive benzinga earnings --ticker AAPL
massive benzinga guidance --ticker AAPL
massive benzinga analysts --analyst-id 12345
```

### Economy

```bash
massive economy inflation --date-gte 2025-01-01 --date-lte 2025-12-31
massive economy labor-market --date-gte 2025-01-01
massive economy treasury-yields --date-gte 2025-01-01
```

### ETF Global

```bash
massive etf-global analytics --ticker SPY
massive etf-global constituents --ticker SPY
```

### TMX (Canadian Markets)

```bash
massive tmx corporate-events --ticker RY
```

## Using with AI Agents

Massive CLI is designed to work well as a tool for AI coding assistants. Any agent that can execute shell commands can use it.

### Pipe JSON into context

```bash
# Get structured data an agent can reason about
massive stocks bars AAPL --from 2025-01-01 --to 2025-01-31 -o json | head -100

# Get snapshot data for analysis
massive stocks snapshots gainers -o json

# Check market status before making trading decisions
massive stocks market-ops status -o json
```

### Use in scripts

```bash
#!/bin/bash
# Compare technical indicators across multiple tickers
for ticker in AAPL MSFT GOOGL AMZN; do
  echo "=== $ticker ==="
  massive stocks indicators rsi "$ticker" --from 2025-01-01 --to 2025-01-31
  echo
done
```

### Stream data into a file for analysis

```bash
# Capture 5 minutes of market data
massive ws stocks agg-minute --all -o json > market_data.jsonl &
sleep 300
kill %1
```

## Development

```bash
# Build
go build -o massive .

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/api/...
```

## License

Copyright (c) 2026. All rights reserved.
