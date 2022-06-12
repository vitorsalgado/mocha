# Contributing

Contribution guidelines for this repository.

## Commit Message Format

This project uses the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard for commit
messages.  
Each commit message consists of a header, a body, and a footer.  
The header is mandatory and must conform to the following format:

```
<type>(<scope>): <short summary>
```

The `<type>` and `<summary>` fields are mandatory, the `(<scope>)` field is optional.

### Type

Must be one of the following:

* build: Changes that affect the build system or external dependencies
* ci: Changes to our CI configuration files and scripts
* docs: Documentation only changes
* feat: A new feature
* fix: A bug fix
* perf: A code change that improves performance
* refactor: A code change that neither fixes a bug nor adds a feature
* test: Adding missing tests or correcting existing tests
