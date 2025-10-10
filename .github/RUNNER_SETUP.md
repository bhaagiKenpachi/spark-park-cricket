# Self-Hosted Runner Setup Guide

## Docker Permission Fix

If you see this error:
```
permission denied while trying to connect to the Docker daemon socket
```

### Solution

Run these commands **on your self-hosted runner machine**:

```bash
# 1. Check which user is running the GitHub Actions runner
ps aux | grep "Runner.Listener"

# The user is likely 'runner' or the user who installed the runner

# 2. Add that user to the docker group (replace 'runner' with your actual user)
sudo usermod -aG docker runner

# 3. Restart the runner service
sudo systemctl restart actions.runner.*

# OR if runner is not a systemd service, stop and start it manually:
# Stop the runner (Ctrl+C if running in terminal)
# Then start it again: ./run.sh

# 4. Verify the user is in docker group
groups runner

# 5. Test Docker access
sudo -u runner docker ps
```

## Common Runner Users

- `runner` - Default user name
- `ubuntu` - On Ubuntu systems
- `github` - Some setups use this
- Your username - If you installed runner as yourself

## Alternative: Run Without Sudo

If you cannot restart the runner immediately, you can temporarily fix by changing docker.sock permissions (NOT recommended for production):

```bash
sudo chmod 666 /var/run/docker.sock
```

**Warning:** This makes Docker accessible to all users, which is a security risk. Use the usermod method above instead.

## Verify Setup

After fixing permissions, test that the runner can access Docker:

```bash
# Switch to the runner user
sudo -u runner docker ps

# Should show running containers without errors
```

## Runner Service Management

### If runner is a systemd service:
```bash
# Check status
sudo systemctl status actions.runner.*

# Restart
sudo systemctl restart actions.runner.*

# View logs
sudo journalctl -u actions.runner.* -f
```

### If runner runs as a process:
```bash
# Find the runner process
ps aux | grep Runner.Listener

# Kill it
sudo kill <PID>

# Restart it (as the runner user)
sudo -u runner /opt/github-runner-bhaagi-spark-park-cricket/run.sh
```

## After Fixing

1. Restart the GitHub Actions runner
2. Go to GitHub Actions tab
3. Re-run the failed workflow
4. The Docker permission error should be resolved

