## My Person Blog Source

#### Sync gacha

```bash
cd gosrc
# Arknights/明日方舟
go run ./cmd/arknightsv2
# TODO Genshin Impact/原神
# Honkai: Star Rail/崩坏：星穹铁道
go run ./cmd/star-rail-wish
# Zenless Zone Zero/绝区零
go run ./cmd/zzz
```

## Update theme manually

Run:

```bash
hugo mod get -u github.com/CaiJimmy/hugo-theme-stack/v3
hugo mod tidy
```

> This starter template has been configured with `v3` version of theme. Due to the limitation of Go module, once the `v4` or up version of theme is released, you need to update the theme manually. (Modifying `config/module.toml` file)
