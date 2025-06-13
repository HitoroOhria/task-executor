# task-executer

[Task](https://taskfile.dev/) のタスクをインクリメンタルサーチで実行する。
タスクの変数がある場合、その変数の値を指定することができる。

```shell
$ task run:vars
task: [run] go run ./cmd -taskfile=test/Taskfile.vars.yml

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
Enter "REQUIRED"        : foo
--- optional ---
Enter "OPTIONAL1"       : bar
Enter "LOOOOOOOOOOOOOOOOOONG_OPTIONAL2": 
Enter "DEFAULT1"        : 
Enter "DEFAULT2"        : 
---   end   ---

run: task -t test/Taskfile.vars.yml with-all REQUIRED="foo" OPTIONAL1="bar"
VALUE=value
REQUIRED=foo
OPTIONAL1=bar
LOOOOOOOOOOOOOOOOOONG_OPTIONAL2=
DEFAULT1=default1
DEFAULT2=default-base
```

## How to use

```shell
$ task build
$ ./main
```

## Features

- [x] 必須な変数を入力できるようにする
- [x] デフォルト値付きのオプショナルな変数を入力できるようにする
- [ ] 入力プロンプトにデフォルト値の値を入力する
- [ ] 依存先のタスクの変数も入力できるようにする
- [x] include した Taskfile に対応する

### Development

Taskfile 経由で実行する。

```shell
$ task run TF="<taskfile>"
$ task run:vars
$ task run:includes
$ task run:all
```
