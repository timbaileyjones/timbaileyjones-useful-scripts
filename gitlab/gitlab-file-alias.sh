# Add to ~/.bashrc or ~/.zshrc (function; alias can't take the path arg correctly)
# Usage: gitlab-file <relative-path>
# Output: https://your-origin/your-repo/-/blob/<sha>/<path>
# Add '#L<line>' yourself if needed.

gitlab-file() {
  local base path
  base=$(/usr/bin/git ls-remote --get-url origin 2>/dev/null | /usr/bin/sed 's|^git@\([^:]*\):\(.*\)\.git$|https://\1/\2|')
  path="${1:?Usage: gitlab-file <relative-path>}"
  echo "${base}/-/blob/$(/usr/bin/git rev-parse HEAD)/${path}"
}

# One-liner for copy-paste into shell config:
# gitlab-file() { local b path; b=$(/usr/bin/git ls-remote --get-url origin 2>/dev/null | /usr/bin/sed 's|^git@\([^:]*\):\(.*\)\.git$|https://\1/\2|'); path="${1:?Usage: gitlab-file <relative-path>}"; echo "${b}/-/blob/$(/usr/bin/git rev-parse HEAD)/${path}"; }
