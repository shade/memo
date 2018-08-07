# Memo

### Prerequisites

- Golang (version 1.9)
- MySQL (version 5.5)
- Memcache
- Bitcoin node (ABC, Unlimited, etc)

#### Optional
- Statsd

### Setup

- Get repo
    ```sh
    go get github.com/memocash/memo/...
    ```

- Create MySQL database
  - Use charset `utf8_general_ci`

- Create config.yaml in memo directory

    ```yaml
    MYSQL_HOST: 127.0.0.1
    MYSQL_USER: memo_user
    MYSQL_PASS: memo_password
    MYSQL_DB: memo
    
    MEMCACHE_HOST: 127.0.0.1
    MEMCACHE_PORT: 11211
    
    BITCOIN_NODE_HOST: 127.0.0.1
    BITCOIN_NODE_PORT: 8333

    STATSD_HOST: 127.0.0.1
    STATSD_PORT: 8125
    ```

### Running

```sh
go build

# Run action node to collect all memo actions
./memo action-node

# Separately run web server
./memo web --insecure

# Also run the user node to get funding txns from local users
./memo user-node
```

### Notes
- Can take about 30 minutes for the action node to fully sync
- Node can sometimes disconnect while syncing, just restart
- You may see a few errors, these are usually mal-formed memos and can be ignored


### View

Visit `http://127.0.0.1:8261` in your browser
