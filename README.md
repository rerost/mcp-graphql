## mcp-graphql
名称 | 必須か？ | 説明 | 例
-- | -- | --
URL | Required | GraphQLを接続するパス | http://localhost:3000/graphql
Schema File | Optional | スキーマファイル。Introspection が使えない時などに利用 | ~/go/src/github.com/rerost/schema.json | schema.graphql
Default Header | Optional | LLMから上書き可能。同じキーを参照する場合のみ上書きされる。LLMは読み取り不可。 | {"Authorization": "Bearer ..."}
Disable Mutation | Optional | ミューテーションを無効化する。本番のGraphQLにアクセスするときなど | true

```
mcp-graphql http://localhost:3000 --headers={}
```
