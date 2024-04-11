# github_utils

## Motivation

I've checked my Github repository recently and found a lot of useless repositories(forks).

I've found this repo (yangshun/delete-github-forks)[https://github.com/yangshun/delete-github-forks]
But this is `javascript(node)`.

I've thought if I do this other way with `golang`?
So this is!

## Usage

```
const usageMessage = `
Usage: github_utils <subcommand> [options]

Available Subcommands:
  fetch    Fetches repositories from GitHub and saves to file
  list     Lists repositories from file
  remove   Removes repositories (interactive with confirmation)
`

const removeCmdUsageMessage = `
Usage: github_utils remove <subcommand>

Available Subcommands:
  all      Removes all repositories with confirmation
  					[You can manually remove single repo from json file]
  check    Asks confirmation for each single repository
  					y 	remove
  					n 	skip
  					q 	quit WITHOUT ANY REMOVE
  					s 	Skip ALL NEXT
`
```

## License

MIT
