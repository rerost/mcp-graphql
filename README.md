## mcp-graphql

例: GitHubのGraphQLを利用する場合
```json
{
  "mcpServers": {
    "graphql": {
      "command": "go",
      "args": [
        "run",
        "github.com/rerost/mcp-graphql@latest",
        "--endpoint", "https://api.github.com/graphql",
        "--headers", "Authorization=Bearer <TOKEN>"
      ]
    }
  }
}
```

https://claude.ai/share/c8e86cdf-81f0-4499-9c85-f4b4645a2756

### Debug
Use inspector https://github.com/modelcontextprotocol/inspector

```
npx @modelcontextprotocol/inspector go run main.go --endpoint=https://api.github.com/graphql --headers '"Authorization=Bearer <TOKEN>"'
```
