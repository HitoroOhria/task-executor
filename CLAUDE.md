# Task Executor プロジェクト概要

## 概要
[Task](https://taskfile.dev/) のタスクをインクリメンタルサーチで実行するCLIツール。  
タスクの変数がある場合、リッチなTUIで変数の値を指定できます。

## 技術スタック
- **言語**: Go 1.24.3
- **主要ライブラリ**:
  - [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUI フレームワーク
  - [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) - TUI コンポーネント
  - [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - スタイリング
  - [go-task/task/v3](https://github.com/go-task/task) - Taskfile パーサー

## プロジェクト構造
```
├── cmd/main.go                     # エントリポイント
├── adapter/                        # 実装アダプター
│   ├── printer_impl.go
│   ├── runner_impl.go
│   └── variable_inputter_impl.go
├── domain/                         # ドメインロジック
│   ├── console/                    # コンソール関連インターフェース
│   ├── model/                      # Taskfile モデル
│   └── value/                      # 値オブジェクト
├── io/                            # I/O処理
├── poc/                           # プロトタイプ・学習用コード
└── test/                          # テスト用Taskfile群
```

## 機能
1. **インクリメンタルサーチ**: タスクをインクリメンタルサーチで選択
2. **リッチな変数入力**: TUIで変数値を視覚的に入力
3. **Taskfile仕様対応**:
   - 必須変数とオプショナル変数
   - デフォルト値付き変数
   - 依存タスクの変数
   - 包含されたTaskfileの変数
   - 内部タスクの除外

## アーキテクチャの特徴
- **ELTパターン**: Taskfileデータを未加工で格納し、参照時に変換
- **値オブジェクト**: 型レベルでの値保証（タスク名・完全タスク名の区別など）
- **依存注入**: ドメインモデルにコンソール操作を注入

## 使用方法
```shell
$ task build
$ ./task-executor
# または
$ ./task-executor -taskfile=path/to/Taskfile.yml
```

## 開発・テスト
```shell
$ task run TF="<taskfile>"          # 指定Taskfileで実行
$ task run:vars                     # 変数テスト用
$ task run:includes                 # 包含テスト用
$ task run:another                  # 別パターンテスト用
$ task run:all                      # 全パターンテスト用
```

## コマンド
- **lint**: `go fmt`, `go vet` などのリント処理
- **test**: `go test ./...` でテスト実行
- **build**: `go build -o task-executor ./cmd` でビルド