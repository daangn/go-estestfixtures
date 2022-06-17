
<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

## About The Project
It's inspired [go-testfixtures](https://github.com/go-testfixtures/testfixtures) project.
This module can dump schemas and documents to json of target elasticsearch. Also, can load to target elasticsearch.
It's purpose to create unit test code with elasticsearch process.

Let's see what this got.

## Build With

### 1. go module olivere/elastic
[olivere/elastic](https://github.com/olivere/elastic/search).
This module is good for service server.
If you were use elasticsearch with golang in your workload, this module will help you.
this module give you agility for your workloads when you get or put documents.

### 2. go module go-testcontainer
[testcontainers/testcontainers-go](https://github.com/). This module is good for creating test code with docker.
go-estestfixtures project needs to insurance generated files work on real elasticsearch process.
You can see `esfixture_test.go`.

## Usage

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
- 파일은 `WithDirectory()`에다가 파일을 만듭니다
- `__schema`, `__document` 이름이 달려있게 됩니다.
- **`__schema`의 내용에서는 `.mappings._meta` 에서 `__(prefix)` 가 달려있는 값은 변경하지 마세요.**

### `Load()`

- 파일은 `WithDirectory()`에 존재하는 내용대로 target elasticsearch에 인덱스를 생성합니다.
- **해당 디렉토리에 있는내용을 전부 ES에 만드는게 아니라 targetNames `WithTargetNames` 에 적어주신 이름에 기반하여 만들도록 합니다.**

### `ClearElasticsearch()`

- **해당 디렉토리에 있는내용을 전부 ES에 만드는게 아니라 targetNames `WithTargetNames` 에 적어주신 이름에 기반하여 지웁니다.**
-  지울때에는 `mappings._meta`에 있는 내용을 확인하고 esfixture에서 만든 index의 경우만 제거합니다. 그게 아니거나 의도와 다르게 동작하는거 같은 경우 에러를 내보냅니다.


## Roadmap
- [x] Release `v0.0.1`
- [ ] Add Additional Templates w/ examples
- [ ] Coverage 90% about to v0.0.1
- [ ] Multi-language Support
  - [ ] Korean

See the [open issues](https://github.com/daangn/go-estestfixtures/issues) for a full list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Write some code
4. Run test (it's required docker for running go unit test)
5. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
6. Push to the Branch (`git push origin feature/AmazingFeature`)
7. Open a Pull Request

## License
Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/daangn/go-estestfixtures.svg?style=for-the-badge
[contributors-url]: https://github.com/daangn/go-estestfixtures/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/daangn/go-estestfixtures.svg?style=for-the-badge
[forks-url]: https://github.com/daangn/go-estestfixtures/network/members
[stars-shield]: https://img.shields.io/github/stars/daangn/go-estestfixtures.svg?style=for-the-badge
[stars-url]: https://github.com/daangn/go-estestfixtures/stargazers
[issues-shield]: https://img.shields.io/github/issues/daangn/go-estestfixtures.svg?style=for-the-badge
[issues-url]: https://github.com/daangn/go-estestfixtures/issues
[license-shield]: https://img.shields.io/github/license/daangn/go-estestfixtures.svg?style=for-the-badge
[license-url]: https://github.com/daangn/go-estestfixtures/blob/master/LICENSE.txt
