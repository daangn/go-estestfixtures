# go-esfixture
elasticsearch 에 대한 golang용 fixture package.

# requirements
1. docker:  `go test ./...` 할 때에는 docker를 필요로 합니다.

# dependencies

### 1. go module olivere/elastic
elasticsearch로 서비스를 하는경우에 olivere/elastic go module package를 자주 쓰기에 dependency로 넣어서 사용함.

### 2. go module go-testcontainer
service source code 작성시에, docker container로 unit test를 만들거나 integration test를 돌릴때가 많은데

# features

### `NewLoader()`

- 객체를 생성합니다. 대략 이런느낌으로 작성하시면 됩니다.
- 필수값
  - WithTargetNames는 필수값입니다.
``` go
l, err := NewLoader(
    tt.args.ctx,
    "http://localhost9200,
    WithSearchFunc(func(c *elastic.Client, targetNames []string) *elastic.SearchService {
        return c.Search(targetNames...).Size(0).From(10)
    }),
    WithLimit(tt.fields.limit),
    WithTargetNames(tt.fields.targetNames...),
    WithDirectory("./fixturedata"), //  설정하지 않으면 default의 디렉토리 `./testdata/esfixtures`
)
if (err != nil) != tt.wantErr {
    t.Fatalf("Loader.Load() error = %v", err)
}
```

### `Dump()`

- target elasticsearch로 부터 targetNames의 내용대로 가져옵니다.
- 파일은 WithDirectory에다가 파일을 만듭니다
- __schema, __document 이름이 달려있게 됩니다.
- **__schema의 내용에서는 `.mappings._meta` 에서 `__(prefix)` 가 달려있는 값은 변경하지 마세요.**

### `Load()`

- 파일은 WithDirectory에 존재하는 내용대로 target elasticsearch에 인덱스를 생성합니다.
- **해당 디렉토리에 있는내용을 전부 ES에 만드는게 아니라 targetNames `WithTargetNames` 에 적어주신 이름에 기반하여 만들도록 합니다.**

### `ClearElasticsearch()`

- **해당 디렉토리에 있는내용을 전부 ES에 만드는게 아니라 targetNames `WithTargetNames` 에 적어주신 이름에 기반하여 지웁니다.**
-  지울때에는 `mappings._meta`에 있는 내용을 확인하고 esfixture에서 만든 index의 경우만 제거합니다. 그게 아니거나 의도와 다르게 동작하는거 같은 경우 에러를 내보냅니다.


# There's not a release version yet.
- maybe you guys want it I'll immediately release version 0.0.1 version.

# License
My car driver license from korean.