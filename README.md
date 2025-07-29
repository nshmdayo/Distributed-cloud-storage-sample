# 分散型クラウドストレージ (Distributed Cloud Storage)

ブロックチェーン技術を活用した分散型クラウドストレージシステム

## 概要

このプロジェクトは、ブロックチェーンとP2P技術を組み合わせて、中央集権的なサーバーに依存しない分散型クラウドストレージシステムを構築します。

### 主要機能

- **分散ファイル保存**: ファイルを暗号化して複数ノードに分散保存
- **ブロックチェーン管理**: ファイルメタデータとアクセス権限をブロックチェーンで管理
- **P2Pネットワーク**: ピアツーピア通信による高い可用性
- **インセンティブシステム**: ストレージ提供者への報酬システム
- **冗長性**: データ損失を防ぐための複製機能

## 技術スタック

- **Go**: バックエンド開発言語
- **libp2p**: P2Pネットワーク通信
- **Ethereum/Polygon**: スマートコントラクト
- **IPFS**: 分散ファイルシステム
- **Gin**: REST APIフレームワーク

## クイックスタート

### 前提条件

- Go 1.21以上
- Docker & Docker Compose
- Make

### インストール

```bash
git clone https://github.com/nshmdayo/Distributed-cloud-storage-sample.git
cd Distributed-cloud-storage-sample
make setup
```

### ストレージノード起動

```bash
make run-node
```

### API サーバー起動

```bash
make run-api
```

## プロジェクト構造

```
/
├── cmd/              # 実行可能ファイル
├── internal/         # 内部パッケージ
├── pkg/             # 公開パッケージ
├── contracts/       # スマートコントラクト
├── scripts/         # ビルド・デプロイスクリプト
├── docs/           # ドキュメント
└── tests/          # テストファイル
```

## 貢献

プロジェクトへの貢献を歓迎します。詳細は [CONTRIBUTING.md](CONTRIBUTING.md) をご覧ください。

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルをご覧ください。