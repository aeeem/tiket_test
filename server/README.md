# Technical test for Technical test Ayo.co.id


## Run Locally

Clone the project

```bash
  git clone https://link-to-project
```

Go to the project directory

```bash
  cd my-project
```


Rename config-example.json to config.json

```bash
  mv config-example.json config.json
```

tidy up dependency

```bash
  go mod tidy
```

Start the tiket_test/server

```bash
  make dockerbuild
```

the docs can be seen at
```bash
  http://172.17.0.1:9292/swagger/index.html#
```

and can be tested trough there