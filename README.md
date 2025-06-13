# task-executer

[Task](https://taskfile.dev/) のタスクをインクリメンタルサーチで実行する。
タスクの変数がある場合、その変数の値を指定することができる。

```shell
$ task run
task: [run] go run ./cmd -taskfile=Taskfile.test.yml

# select task
QUERY>
 simple:                  Simple command                                                                                                                                                                                                                   
 with-vars:               Command with vars
 with-optional-vars:      Command with optional vars
 with-default-vars:       Command with default vars
 with-requires:           Command with requires
 with-long-vars:          Command with long vars
 with-all:                Command with all pattern

# select with-all
--- required ---
Enter "REQUIRED"         (required): foo
--- optional ---
Enter "OPTIONAL1"        (optional): bar
Enter "LOOOOOOOOOOOOOOOOOONG_OPTIONAL2" (optional): baz
Enter "DEFAULT1"         (optional): 
Enter "DEFAULT2"         (optional): 
---   end   ---

run: task -t Taskfile.test.yml with-all REQUIRED="foo" OPTIONAL1="bar" LOOOOOOOOOOOOOOOOOONG_OPTIONAL2="baz"
VALUE=value
OPTIONAL1=bar
LOOOOOOOOOOOOOOOOOONG_OPTIONAL2=baz
REQUIRED=foo
DEFAULT1=default1
DEFAULT2=default-base
```

### 実行方法

ビルドして経由で実行する。

```shell
$ task build
$ ./main
```

### 動作確認方法

Taskfile 経由で実行する。

```shell
$ task run
```

## Features

- [x] 必須な変数を入力できるようにする
- [x] デフォルト値付きのオプショナルな変数を入力できるようにする
- [ ] 入力プロンプトにデフォルト値の値を入力する
- [ ] 依存先のタスクの変数も入力できるようにする
- [ ] include した Taskfile に対応する
