# task-executer

[Task](https://taskfile.dev/) のタスクをインクリメンタルサーチで実行する。
タスクの変数がある場合、その変数の値を指定可能である。

```shell
$ task run
task: [run] go run ./... -taskfile=Taskfile.test.yml

# select task
QUERY>
 simple:                                    Simple command
 with-vars:                                 Command with vars
 with-requires:                             Command with requires
 with-vars-and-requires:                    Command with vars and requires
 with-vars-and-requires-and-long-var:       Command with vars and requires and long var

# e.g. selectd with-vars-and-requires
# input variable value and execute task
Enter "REQUIRED"  (required): foo
Enter "OPTIONAL1" (optional): bar
run: task -t Taskfile.test.yml with-vars-and-requires REQUIRED="foo" OPTIONAL1="bar"
OPTIONAL1=bar VALUE=bar
REQUIRED=foo
```

### 実行方法

ビルドして経由で実行する。

```shell
$ task build
$ ./main
```

Taskfile 経由で実行する。

```shell
$ task run TF="Taskfile.yml"
```

