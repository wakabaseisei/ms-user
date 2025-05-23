# ms-user

このリポジトリは、MSプロジェクトにおける **ユーザー情報を管理するマイクロサービス** です。  
Aurora MySQL Serverless v2 と IAM 認証を活用し、セキュアでスケーラブルなユーザーデータの登録・取得を行います。

---

## 📦 主な責務

- **ユーザー作成 / 取得 / ヘルスチェック API の提供**  
  Connect RPC を用いてエンドポイントを提供

- **永続化処理**  
  Aurora MySQL Serverless v2 へ IAM 認証で接続

- **DBマイグレーションの自動化**  
  `db/migrations/` に変更があった場合のみ、CI経由で Lambda にデプロイ・実行

- **環境変数による柔軟な設定**  
  `DB_HOST` / `DB_NAME` / `DB_USER` / `DB_PORT` などの変数で、環境に応じて接続先を切り替え

---

## 📁 ディレクトリ構成（抜粋）

```
.
├── db/
│ └── migrations/ # マイグレーションSQL
├── internal/
│ ├── cmd/
│ │ ├── app/ # HTTPサーバー（Connect対応）
│ │ └── cli/ # マイグレーション実行用Lambda CLI
│ ├── config/ # DBやAWS設定のenv定義と初期化
│ ├── domain/ # ドメインモデルとリポジトリインターフェース
│ ├── driver/ # gRPCハンドラの実装（UserService）
│ ├── repository/ # MySQLとの接続・永続化処理（IAM認証）
│ └── usecase/ # ユースケース層の実装（CreateUser）
├── Dockerfile # アプリケーション用
├── Dockerfile.migrate # マイグレーション用Lambda向け
└── .github/workflows/ # CI（ECRビルド + マイグレーション自動実行）
```

---

## 🔁 CI/CD

- GitHub Actions による ECR ビルド
- `db/migrations/` に差分がある場合のみ:
  - 専用イメージ（`-migrate` タグ）をビルド
  - Lambda にデプロイ & 実行

---

## ✅ 補足

- IAM 認証により Aurora への接続情報は Secrets Manager に依存せず、ローテーションも容易
- 今後、スキーマ変更時はマイグレーションLambdaを通じて自動反映
