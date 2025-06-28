# task-executor

[Task](https://taskfile.dev/) のタスクをインクリメンタルサーチで実行する。
タスクの変数がある場合、その変数の値を指定することができる。

![CleanShot_2025-06-16_11-24-01](https://github.com/user-attachments/assets/5617d746-527c-4f26-83f0-d10a39866114)

```shell
$ task run:vars
task: [run:vars] go run ./cmd -taskfile=test/Taskfile.vars.yml

# 1. select task
QUERY>
 simple:                  Simple command
 with-vars:               Command with vars
 with-optional-vars:      Command with optional vars
 with-default-vars:       Command with default vars
 with-requires:           Command with requires
 with-long-vars:          Command with long vars
 with-all:                Command with all pattern

# 2. input value
Input variable values.

Variable            Req.  Def.        Value               
──────────────────────────────────────────────────────────
REQUIRED             ✓                foo                  
OPTIONAL1                             bar                  
LOOOOOOOOOOOOOOOO…                                         
DEFAULT1                  [default1]                       
DEFAULT2                                                   
DUPLICATE                                                  

(enter to finish)

# 3. run task
┌─────────────────────────────────────────────┐
│ [run]                                       │
│ task -t test/Taskfile.vars.yml with-all \   │
│     REQUIRED=foo \                          │
│     OPTIONAL1=bar \                         │
│     LOOOOOOOOOOOOOOOOOONG_OPTIONAL2=baz \   │
│     DEFAULT1=default1                       │
└─────────────────────────────────────────────┘
VALUE=value
REQUIRED=foo
OPTIONAL1=bar
LOOOOOOOOOOOOOOOOOONG_OPTIONAL2=baz
DEFAULT1=default1
DEFAULT2=default-base
DUPLICATE=
```

## How to use

```shell
$ task build
$ ./task-executor
```

## Features and TODOs

- [x] 必須な変数を入力できるようにする
- [x] デフォルト値付きのオプショナルな変数を入力できるようにする
- [x] 入力プロンプトにデフォルト値の値を入力する
- [x] 依存先のタスクの変数も入力できるようにする
- [x] include した Taskfile に対応する
- [x] リッチ TUI で変数を入力する
- [ ] 入力受付の終了後に textinput のカーソルを非表示にする
- [ ] エラーハンドリングとユーザー向けメッセージの追加
- [ ] ANSI エスケープシーケンスの入力を回避する
  - https://chatgpt.com/c/684ed0d9-346c-8003-b05d-7cf2f1e4effc
- [ ] `[run]` をボーダーに埋め込む
  - https://chatgpt.com/c/684ed686-8d50-8003-ad84-694b850ff096
- [ ] run に Copy ボタンを追加する
- [ ] インクリメンタルサーチの TUI を実装する
- [ ] より実践的な形式に対応する

## Development

Run program via Taskfile.

```shell
$ task run TF="<taskfile>"
$ task run:vars
$ task run:includes
$ task run:deps
$ task run:all
```

## Learning programming

- 値オブジェクトが便利
  - 型レベルで値を保証できる
    - 特に「タスク名」と「完全タスク名」の存在がややこしく、それらの区別を型レベルで保証できることは大きい
  - メソッドを追加できる
    - 値オブジェクトに関するロジックを集約することで、コードがシンプルになり見通しが良くなる
- 似たような概念がある時、構造体を分けるかどうか
  - 「必須変数」と「オプショナル変数」
    - 構造体を分けた
      - `RequiredVar` と `OptionalVar` を定義した
    - なぜか
      - 本ツールには、先に必須変数を入力して後からオプショナル変数を入力する、という要件がある
      - その要件を満たすには、最初から構造体を分けて、それぞれ別のスライスに格納すると、開発しやすかった
    - 懸念点
      - ロジックの重複が発生するが、許容できるほどにシンプルな内容だった
  - 「コマンド」と「依存タスクコマンド」
    - 構造体を分けなかった
      - `Cmd` にまとめ、ポインタ型のフィールドを2つ持たせた
    - なぜか
      - 現状、分ける理由がなかった
      - コマンドの方は使用する機会がなく、区別するための概念としてしか登場していない
      - 2個に分けるのは早すぎる抽象化のように思えた
- CLI コマンドの実行を、モデルから分離するか、依存させるか
  - 依存させた
    - モデルの中に `Runner` のインターフェースを持ち、実体を DI して実装した
  - なぜか
    - モデル対象となる Taskfile と、行いたい「変数の入力」「タスクの実行」の CLI 操作は、結合度が高いと感じた
    - ユースケースに独立させるのではなく、モデルに依存させた方が、全体の構造がシンプルになると考えた
  - 懸念点
    - パッケージの依存関係が複雑になった
    - 欲張った結果、`domain/console` -> `domain/value` への依存が発生した
    - モデルと CLI 操作を密結合させたため、`domain/model` -> `domain/console` への依存が発生している
    - オニオンアーキテクチャの `domain/repository` と比べると依存の方向が逆になっており、危ういバランスだと感じる
- TUI 開発が楽しすぎた
  - [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) に感謝
