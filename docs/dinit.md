# Dinit

## Environments

- `DINIT_CONSUL_LEAVE_INTERVAL` - The interval (in seconds) between `consul leave` and `app exit`, default: 5
- `DINIT_EXIT` - Whether dinit exit, default: `true`

## Service Restart In-place

### Descriptions

The main process restarts these child processes after receiving the `restart In-place` signal

### How to Use

Provide two methods

one example:
```bash
kill -s SIGHUP $dinit_pid
```

another example:

```bash
curl http://localhost:8888/init/reload
```
