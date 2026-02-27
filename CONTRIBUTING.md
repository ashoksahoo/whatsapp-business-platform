# Contributing to WhatsApp Business Platform

Thank you for your interest in contributing! This is an open-source project and contributions are welcome.

## Ways to Contribute

- **Bug reports** — Open an issue describing the bug and how to reproduce it
- **Feature requests** — Open an issue describing the feature and use case
- **Bug fixes** — Submit a pull request with the fix
- **New features** — Open an issue first to discuss the approach, then submit a PR
- **Documentation** — Improve docs, fix typos, add examples
- **Tests** — Add missing tests or improve existing coverage

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
```bash
git clone https://github.com/YOUR_USERNAME/whatsapp-business-platform.git
cd whatsapp-business-platform
```

3. **Create a feature branch**
```bash
git checkout -b feature/your-feature-name
```

4. **Make your changes** and ensure everything works

5. **Run tests**
```bash
go test ./...
cd frontend && npm run lint
```

6. **Commit your changes**
```bash
git commit -m "feat: add your feature description"
```

7. **Push and open a PR**
```bash
git push origin feature/your-feature-name
```

## Development Setup

See the [README](README.md) for full setup instructions.

Quick start:
```bash
cp .env.example .env
# Fill in your WhatsApp credentials
go run cmd/server/main.go
```

## Commit Convention

Use conventional commits:

- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation changes
- `test:` — Adding or updating tests
- `refactor:` — Code change that's not a fix or feature
- `chore:` — Maintenance tasks

## Pull Request Guidelines

- Keep PRs focused — one feature or fix per PR
- Include tests for new functionality
- Update documentation if needed
- Ensure `go test ./...` passes
- Ensure `npm run lint` passes in `frontend/`
- Write a clear PR description explaining what and why

## Reporting Bugs

Open a GitHub issue with:
- Go version and OS
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages

## Feature Requests

Open a GitHub issue with:
- The problem you're trying to solve
- Your proposed solution
- Any alternatives you've considered

## Code Style

**Go**
- Follow standard Go conventions (`gofmt`)
- Run `go vet ./...` before submitting
- Keep functions small and focused

**TypeScript/React**
- Follow the existing code patterns
- Run `npm run lint` before submitting

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for helping make this project better!**
