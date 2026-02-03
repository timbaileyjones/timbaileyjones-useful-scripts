#!/usr/bin/env python3
"""
Prune local branches whose remote tracking branch no longer exists.

Handles branches checked out in worktrees gracefully, with option to remove them.

By default, runs in dry-run mode. Use --approve to actually delete.
"""

import subprocess
import sys
import re


def run_git(*args: str) -> tuple[int, str]:
    """Run a git command and return (returncode, output)."""
    result = subprocess.run(
        ["git", *args],
        capture_output=True,
        text=True,
    )
    return result.returncode, result.stdout.strip()


def get_gone_branches() -> list[dict]:
    """Get branches whose remote tracking branch is gone."""
    _, output = run_git("branch", "-vv")
    branches = []

    for line in output.splitlines():
        if ": gone]" not in line:
            continue

        line = line.strip()
        is_current = line.startswith("*")
        is_worktree = line.startswith("+")

        # Strip the marker and leading whitespace
        if is_current or is_worktree:
            line = line[1:].strip()

        # Extract branch name (first field)
        branch_name = line.split()[0]

        # Extract worktree path if present (shown in parentheses for + branches)
        worktree_path = None
        if is_worktree:
            match = re.search(r"\(([^)]+)\)", line)
            if match:
                worktree_path = match.group(1)

        branches.append({
            "name": branch_name,
            "is_current": is_current,
            "is_worktree": is_worktree,
            "worktree_path": worktree_path,
        })

    return branches


def prune_branches(approve: bool = False, include_worktrees: bool = False) -> None:
    """Prune gone branches, optionally including worktree removal."""
    dry_run = not approve

    # First, fetch and prune remote tracking refs
    print("Fetching and pruning remote tracking refs...")
    run_git("fetch", "--prune")

    branches = get_gone_branches()

    if not branches:
        print("No local branches with gone remotes found.")
        return

    deleted = []
    worktrees_removed = []
    skipped_current = []
    skipped_worktree = []

    for branch in branches:
        name = branch["name"]

        if branch["is_current"]:
            skipped_current.append(branch)
            continue

        if branch["is_worktree"]:
            if include_worktrees:
                # Remove worktree first, then branch
                wt_path = branch["worktree_path"]
                if dry_run:
                    print(f"Would remove worktree: {wt_path}")
                    print(f"Would delete branch: {name}")
                    worktrees_removed.append(branch)
                    deleted.append(branch)
                else:
                    # Remove the worktree
                    rc, output = run_git("worktree", "remove", wt_path)
                    if rc == 0:
                        print(f"Removed worktree: {wt_path}")
                        worktrees_removed.append(branch)
                        # Now delete the branch
                        rc, output = run_git("branch", "-D", name)
                        if rc == 0:
                            print(f"Deleted branch: {name}")
                            deleted.append(branch)
                        else:
                            print(f"Failed to delete {name}: {output}", file=sys.stderr)
                    else:
                        print(f"Failed to remove worktree {wt_path}: {output}", file=sys.stderr)
                        skipped_worktree.append(branch)
            else:
                skipped_worktree.append(branch)
            continue

        # Delete the branch
        if dry_run:
            print(f"Would delete branch: {name}")
            deleted.append(branch)
        else:
            returncode, output = run_git("branch", "-D", name)
            if returncode == 0:
                print(f"Deleted branch: {name}")
                deleted.append(branch)
            else:
                print(f"Failed to delete {name}: {output}", file=sys.stderr)

    # Summary
    print()
    action = "Would delete" if dry_run else "Deleted"

    if deleted:
        print(f"{action} {len(deleted)} branch(es).")

    if worktrees_removed:
        wt_action = "Would remove" if dry_run else "Removed"
        print(f"{wt_action} {len(worktrees_removed)} worktree(s).")

    if skipped_current:
        print(f"\nSkipped {len(skipped_current)} branch(es) (currently checked out):")
        for b in skipped_current:
            print(f"  - {b['name']}")

    if skipped_worktree:
        print(f"\nSkipped {len(skipped_worktree)} branch(es) (checked out in worktrees):")
        for b in skipped_worktree:
            print(f"  - {b['name']}")
            if b["worktree_path"]:
                print(f"    Worktree: {b['worktree_path']}")
        print("\nTo include these, add --include-worktrees flag.")

    if dry_run and (deleted or skipped_worktree):
        print("\n" + "=" * 50)
        print("This was a dry run. To apply changes, use --approve")
        if skipped_worktree:
            print("To also remove worktrees, use --approve --include-worktrees")


def main() -> None:
    args = sys.argv[1:]

    if "--help" in args or "-h" in args:
        print(__doc__)
        print("Usage: prune-gone-branches.py [OPTIONS]")
        print()
        print("Options:")
        print("  --approve            Actually delete branches (default is dry-run)")
        print("  --include-worktrees  Also remove worktrees for gone branches")
        print("  --help, -h           Show this help message")
        sys.exit(0)

    approve = "--approve" in args
    include_worktrees = "--include-worktrees" in args

    if not approve:
        print("=== DRY RUN MODE ===\n")

    prune_branches(approve=approve, include_worktrees=include_worktrees)


if __name__ == "__main__":
    main()
