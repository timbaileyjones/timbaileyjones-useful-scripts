# tidy-screenshots

Keeps macOS screenshots and screen recordings from piling up on the Desktop.

- **Renames captures as soon as they're taken.** A launchd agent watches `~/Desktop`,
  and turns `Screenshot 2026-09-30 at 2.28.16 PM.png` into
  `screenshot-20260930-142816-1.png` (screen recordings become `recording-…-1.mov`).
  The names sort chronologically and have no spaces.
- **Archives old captures.** Every 4 hours, cron moves captures older than 24 hours
  into `~/Desktop/old`.

The timestamp comes from the capture time in the macOS filename, converted to 24-hour
time. If it can't be parsed, the file's creation time is used instead. `N` starts at 1
and increments when two captures share a timestamp, so nothing is ever overwritten.

## Usage

```sh
tidy-screenshots rename [dir]   # rename captures in place (default: ~/Desktop)
tidy-screenshots archive [-n]   # move captures older than 24h to ~/Desktop/old (-n = dry run)
```

`rename ~/Desktop/old` is handy for a one-time cleanup of captures archived before
this was set up.

## Install

```sh
cp tidy-screenshots ~/.local/bin/
chmod +x ~/.local/bin/tidy-screenshots

# Renamer: launchd agent that watches the Desktop.
# The plist has hardcoded /Users/tim paths; edit them for another user.
cp com.tim.tidy-screenshots.plist ~/Library/LaunchAgents/
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.tim.tidy-screenshots.plist

# Archiver: every 4 hours via cron.
(crontab -l; echo '0 */4 * * * $HOME/.local/bin/tidy-screenshots archive >> $HOME/Library/Logs/tidy-screenshots.log 2>&1') | crontab -
```

### Full Disk Access (required)

macOS blocks background processes from reading `~/Desktop`. In
**System Settings → Privacy & Security → Full Disk Access**, add:

- `/bin/bash`, for the launchd agent. This is a broad grant: any bash script
  run in the background gets it too.
- `/usr/sbin/cron`, for the archiver.

Without these, the log shows `cannot read /Users/…/Desktop (grant Full Disk Access?)`.
A macOS update can reset these permissions, so check them first if the jobs stop working.

## Logs

Both jobs append to `~/Library/Logs/tidy-screenshots.log`.

## Uninstall

```sh
launchctl bootout gui/$(id -u)/com.tim.tidy-screenshots
rm ~/Library/LaunchAgents/com.tim.tidy-screenshots.plist ~/.local/bin/tidy-screenshots
crontab -l | grep -v tidy-screenshots | crontab -
```

## Notes

- cron skips runs while the Mac is asleep; the next run catches up.
- `old/` is never pruned.
