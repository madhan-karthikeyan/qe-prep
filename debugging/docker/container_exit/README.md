## Symptoms
The container builds and runs but exits immediately with status 0. `docker ps -a` shows `Exited (0)` right away.

## Root Cause
A container only lives as long as its main process (PID 1). The Python script prints three lines and exits, so the container exits immediately after startup.

## Fix
Run a long-lived process: an HTTP server, a message consumer, a sleep loop, or `tail -f /dev/null` for dev containers.

## Prevention
- Ensure the container's `CMD`/`ENTRYPOINT` runs a process that stays alive.
- Use init systems (tini, dumb-init) for proper signal handling in production.
- Test with `docker run -it` to see output interactively.
