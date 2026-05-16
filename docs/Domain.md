# ドメイン単語帳

## ドメイン単語帳とは

ドメイン単語帳は、特定のドメインに関連する用語や概念をまとめたリストです。
ユビキタス言語はプロダクト全体で1つではなく、コンテキスト単位で定義します。

---

## 境界づけられたコンテキスト（Bounded Context）

| コンテキスト | 役割 |
|------------|------|
| User / Account | ユーザー、GitHub アカウント連携、認証 |
| GitHub Integration | GitHub OAuth・API連携、活動データの取得と正規化 |
| Game / Faction | 言語勢力・ポイント・ランキング・シーズンの中核ルール |
| Presentation / API | 画面向け API、表示用集約 |

---

## 用語集

### User / Account コンテキスト

| 用語 | 英語表記 | 定義 |
|------|---------|------|
| ユーザー | User | アプリを利用するプレイヤー。GitHub アカウントと紐づく。 |
| GitHub アカウント | GitHubAccount | ユーザーに紐づく GitHub の連携情報（ID・ユーザー名・アバター URL）。 |
| セッション | Session | ログイン中のユーザーを識別するトークン情報。有効期限 24 時間。 |

---

### GitHub Integration コンテキスト

| 用語 | 英語表記 | 定義 |
|------|---------|------|
| コントリビューション | Contribution | GitHub 上で行われる開発活動。Commit・Pull Request・Issue Close などのアクションを指す。ゲーム内で解析され、CP（活動総合）と SP（技術ごと）の両方に変換される。 |
| コミット | Commit | 1ファイル以上の変更を伴う GitHub へのコード提出。MVP では +1 CP。 |
| プルリクエスト | PullRequest | マージ済みの GitHub Pull Request。MVP では +5 CP。 |

---

### Game / Faction コンテキスト

| 用語 | 英語表記 | 定義 |
|------|---------|------|
| コントリビューションポイント | ContributionPoint / CP | GitHub 上の開発活動を数値化した総合的なポイント。Commit・PR など活動全般を単純集計する。技術ごとには分けない。（検討中）建物建設など一部のギルド活動に SP と併用して消費できる。ただしランキング（勢力値）には影響しない。 |
| スキルポイント | SkillPoint / SP | 技術ごとに独立して存在するポイント。コントリビューションに使用された技術を GitHub API で解析して、その技術の SP が貯まる。例：Golang と TS と Java を使って開発していれば Golang SP / TS SP / Java SP がそれぞれ貯まる。自分のギルド以外の技術の SP も獲得できるが、ギルド活動（建築など）に使えるのは自ギルドの技術の SP のみ。技術の判別は MVP では言語レベルを想定し、将来的には package.json 等のマニフェスト解析によりライブラリ・フレームワーク単位への拡張を検討。 |
| 勢力値 | FactionScore | ギルドの強さを表す値。所属メンバー全員の累計 SP（消費した分は差し引かない） + 建築など別途獲得分の合計。ランキングに用いる。 |
| ギルド | Guild | 特定のプログラミング技術・言語を支持するプレイヤーコミュニティ。プレイヤーは推し技術を選択することで対応するギルドへ所属し、GitHub 上の活動によって得た SP（自ギルド技術分）をギルド活動に投資する。Guild 同士はシーズンごとのランキングや技術覇権を争う。ユーザーにつき1つ。加入・脱退・変更はいつでも可能。 脱退しても消費済みの SP/CP はギルドにそのまま残る。保有情報は下表参照。 |
| ギルドプロフィール | GuildProfile | ギルドが持つ外観情報。名前・説明文・アイコンをカスタマイズ可能。 |
| 称号 | Title | ギルドが条件達成（例：シーズン優勝）によって獲得する特別な肩書き。ギルドプロフィールに表示される。複数所持可能。 |
| 推し技術 | FavoriteTech | ユーザーがギルドを選ぶ際に選択するプログラミング言語・技術。 |
| シーズン | Season | ランキング集計の期間区切り。終了時に最も勢力値が高いギルドが勝利し、次シーズンの世界設定・称号・見た目・イベントに反映される。SP・建物などの研究開発成果はシーズンをまたいでリセットされない。 ランキングは累計勢力値とシーズンごとの勢力値の両方を持つ。 |
| ランキング | Ranking | ギルドの勢力値を比較した順位。 |
| ワールドマップ | WorldMap | 全ギルドの勢力値占有率を可視化した GUI。全ギルドの勢力値合計を 100% として各ギルドがどれだけ占有しているかを表す。「領土」はこのマップ上の UI 表現として使われる言葉であり、ドメイン用語としては WorldMap を使う。 |
| キャラクター | Character | 各技術の特性を擬人化したもの（例：Rust = 鎧武者・無敵の守り、Python = 魔法使い・万能型）。ユーザーが自由に選択して育成・バトルに使う。ギルドとは独立して選べる。 |
| キャラ育成 | CharacterGrowth | ユーザーが CP を消費して個人の推しキャラクターを育てる行動。能力値を上げてバトルに備える。 |
| キャラ能力 | CharacterAbility | キャラクターが持つ能力値・スキル。育成によって強化される。 |
| 研究開発 | Research | ユーザーが能動的に実行するギルド強化行動。SP（自ギルド技術分）や CP を消費して発動し、成果物として建物を獲得する。「ギルド強化」「研究開発」は同じ概念。 |
| 建物 | Building | 研究開発によってギルドが保有できる成果物。複数種類あり、同じ種類でも複数建設可能。 特定の能力（現状の想定は勢力値ボーナス）を持つ。例：勢力値 +500（固定）、時間経過で勢力値が上昇するなど（詳細は検討中）。シーズンをまたいでも消えない。 |
| バトル | Battle | ギルド間の戦闘メカニクス（MVP 後の機能）。各ユーザーが育成した個人の推しキャラクターで戦う。 |
| チャット | Chat | ギルド内のメンバー同士で会話できるギルドチャットと、全ギルドをまたいで会話できるグローバルチャットの2種類。 |
| マイページ | MyPage | ユーザー個人の開発情報を振り返れるページ。コード量・開発言語一覧・コントリビューショングラフ（草）・CP解析情報を閲覧できる。余力で推しキャラクターの育成も行える。 |
| 台帳 | Ledger | ContributionPoint の増減を append-only で記録するもの。直接更新・削除は不可。 |
| 台帳エントリ | LedgerEntry | 台帳の1件。Amount（増減量）・Type（獲得/消費/調整）・Reason・SourceType・SourceID・BalanceAfter を持つ。 |

---

## ギルドが保有する情報

| 種別 | 用語 | 内容 |
|------|------|------|
| プロフィール | GuildProfile | 名前・説明文・アイコン（カスタマイズ可能） |
| メンバー | Member | 所属ユーザーの一覧 |
| 勢力値 | FactionScore | 累計SP合計 + 建物ボーナス合計 |
| 建物 | Building | 研究開発で獲得した成果物。複数種類・複数所持可 |
| 称号 | Title | 条件達成で獲得する肩書き。複数所持可 |
| チャット | GuildChat | ギルドメンバー間のチャット履歴 |

---

## 用語の関係

```
User ──── GitHubAccount
  │
  ├── Guild（ユーザーにつき1つ、加入・脱退・変更はいつでも可能）
  │
  └── Character（推しキャラ、自由に選択）
        └── CharacterGrowth（CP を消費して育成）
              └── Battle（ギルド対戦で使用）

Contribution（GitHub活動: Commit / PullRequest）
  │
  ├── ContributionPoint（CP）：活動全般の総合集計ポイント
  │     ├── LedgerEntry（台帳エントリ）
  │     ├── Research（建物建設に消費・検討中）
  │     └── CharacterGrowth（キャラ育成に消費）
  │
  └── SkillPoint（SP）：技術ごとに独立して貯まるポイント
        例）Golang SP: 200 / TS SP: 100 / Java SP: 50 …
        └── Research（自ギルド技術のSPのみ消費可能）

Guild
  ├── GuildProfile（名前・説明文・アイコン、カスタマイズ可能）
  ├── Member[]（所属ユーザー）
  ├── FactionScore = 所属メンバー全員の累計SP合計 + Building ボーナス合計
  ├── Building[]（研究開発の成果物、複数種類・複数所持可）
  │     └── FactionScore にボーナス加算
  ├── Title[]（称号、条件達成で獲得・複数所持可）
  └── GuildChat（ギルド内チャット）

Season ──── Ranking（ランキング）：期間内の勢力値で集計
              ※ SP・Building・Title はシーズンをまたいでリセットされない

WorldMap（ワールドマップ）
  └── 全ギルドの FactionScore を 100% として占有率を可視化

GlobalChat（グローバルチャット）：全ギルドをまたいだチャット

MyPage
  ├── コード量・開発言語一覧
  ├── コントリビューショングラフ（草）
  ├── CP 解析情報
  └── CharacterGrowth（推しキャラ育成）
```

---

## 実装状況

| 用語 | 実装状態 | 実装箇所 |
|------|---------|---------|
| User | 実装済 | `services/api/internal/domain/user/user.go` |
| GitHubAccount | 実装済 | `services/api/internal/domain/user/user.go` |
| ContributionPoint / LedgerEntry | 実装済 | `services/api/internal/domain/contributionpoint/` |
| SkillPoint / LedgerEntry | 実装済（point_ledger に統合） | `services/api/internal/domain/contributionpoint/` |
| Guild | 未実装（企画定義済） | — |
| Season | 未実装（企画定義済） | — |
| WorldMap | 未実装（企画定義済） | — |
| Character / CharacterGrowth | 未実装（企画定義済） | — |
| Battle | 未実装（MVP 後） | — |
| Chat（ギルド / グローバル） | 未実装（企画定義済） | — |
| MyPage | 未実装（企画定義済） | — |

---

## 関連ドキュメント

- [ProductOverview.md](./ProductOverview.md) — プロダクト概要・ゲーム設計
- [Architecture.md](./Architecture.md) — システムアーキテクチャ
- [BackendStructure.md](./BackendStructure.md) — バックエンド構造詳細
