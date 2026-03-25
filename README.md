# ISS Tracker CLI

ターミナル上でISS（国際宇宙ステーション）の現在位置をリアルタイムに追跡するCLIツールです。ASCII世界地図上にISSの位置を表示し、現在宇宙にいる宇宙飛行士の情報や自分の位置からの距離も確認できます。

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow)

## 機能

- ASCII世界地図上にISSの現在位置をリアルタイム表示
- IPアドレスから自分の位置を推定し、地図上に表示
- ISSまでの距離を計算
- 次回のISS通過予測時刻を表示
- 現在宇宙にいる宇宙飛行士の一覧（宇宙船別）
- 更新間隔のカスタマイズ

## 必要環境

- Go 1.25 以上
- インターネット接続（API通信のため）
- Unicode対応のターミナル

## インストール

```bash
go install github.com/natori-hrj/iss-tracker-cli@latest
```

またはソースからビルド：

```bash
git clone https://github.com/natori-hrj/iss-tracker-cli.git
cd iss-tracker-cli
go build -o iss-tracker .
```

## 使い方

```bash
# デフォルト設定で起動（5秒間隔で更新）
iss-tracker

# 更新間隔を指定（秒）
iss-tracker -i 10
iss-tracker --interval 10
```

### キー操作

| キー | 動作 |
|------|------|
| `r` | 手動で即時更新 |
| `q` / `Esc` / `Ctrl+C` | 終了 |

## 使用API

| API | 用途 | URL |
|-----|------|-----|
| Open Notify - ISS Now | ISS現在位置の取得 | `http://api.open-notify.org/iss-now.json` |
| Open Notify - Astros | 宇宙飛行士情報の取得 | `http://api.open-notify.org/astros.json` |
| ip-api | IPアドレスからの位置推定 | `http://ip-api.com/json/` |

> **注意**: これらのAPIはHTTP（非HTTPS）で通信します。ip-api.comにはあなたのIPアドレスが送信されます。

## プロジェクト構成

```
.
├── main.go              # エントリーポイント
├── cmd/
│   └── root.go          # CLIコマンド定義（cobra）
└── internal/
    ├── api/
    │   └── client.go    # ISS位置・宇宙飛行士情報のAPI通信
    ├── geo/
    │   └── geo.go       # 位置情報取得・距離計算・通過予測
    ├── tui/
    │   └── model.go     # TUIモデル（Bubble Tea）
    └── ui/
        ├── map.go       # ASCII世界地図の描画
        └── styles.go    # 表示スタイル定義（Lip Gloss）
```

## 使用ライブラリ

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUIフレームワーク
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - ターミナルスタイリング
- [Cobra](https://github.com/spf13/cobra) - CLIフレームワーク

## ライセンス

[MIT License](LICENSE)