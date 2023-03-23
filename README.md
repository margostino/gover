# go-to-vercel-bootstrap
Bootstrap a GO Serverless Function hosted by Vercel

ToDo:
- [ ] Async and parallel calls
- [ ] Status for deploy and show go-to-link
- [ ] Proper docu
- [ ] Clear logging main status in one place (command). Impro http calls

# Usage

```bash
> gover --name my-first-serverless-app
```

# Configuration

Gover will search config in home directory first with name `.gover.yml`
```yaml
GITHUB_ACCESS_TOKEN: { *** }
GITHUB_USERNAME: { *** }
VERCEL_ACCESS_TOKEN: { *** }
```