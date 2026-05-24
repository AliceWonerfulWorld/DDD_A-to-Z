// Lang War ウォーマップ用のヘックスグリッド領土計算。
// 軸座標 (q, r)、先端上向き六角形。
// 入力 : ギルドのホーム座標と CP 量（set_guild / set_guild_count で設定）
// 出力 : ヘックスグリッド上の全セルについて、どのギルドがどの強さで支配しているか
//         （calc 呼び出し後に get_cell_* で取得）

// ── ギルド入力（最大16）──────────────────────────────────────────────────────
let num_guilds : Ref[Int] = { val: 0 }
let guild_q : Array[Int] = Array::make(16, 0)
let guild_r : Array[Int] = Array::make(16, 0)
let guild_cp : Array[Int] = Array::make(16, 1)

// ── 出力セル（最大500、半径9のヘックスグリッド = 271セル）─────────────────────
let max_cells : Int = 500
let num_cells : Ref[Int] = { val: 0 }
let out_q : Array[Int] = Array::make(500, 0)
let out_r : Array[Int] = Array::make(500, 0)
let out_guild : Array[Int] = Array::make(500, -1)
let out_intensity : Array[Int] = Array::make(500, 0)
let out_contested : Array[Int] = Array::make(500, 0)

// ── ヘルパー ──────────────────────────────────────────────────────────────────

pub fn hex_dist(q1 : Int, r1 : Int, q2 : Int, r2 : Int) -> Int {
  let dq = q1 - q2
  let dr = r1 - r2
  let ds = -(dq + dr)
  (dq.abs() + dr.abs() + ds.abs()) / 2
}

// ── 公開 API ──────────────────────────────────────────────────────────────────

pub fn set_guild_count(n : Int) -> Unit {
  num_guilds.val = n
}

// idx: 0始まりのインデックス; q,r: 軸座標; cp: 貢献ポイント（≥1）
pub fn set_guild(idx : Int, q : Int, r : Int, cp : Int) -> Unit {
  if idx >= 0 && idx < 16 {
    guild_q[idx] = q
    guild_r[idx] = r
    guild_cp[idx] = if cp < 1 { 1 } else { cp }
  }
}

// 指定した半径のヘックスグリッド上で領土計算を実行する。
// 影響力モデル: influence = cp * 1000 / dist²（逆二乗）
// 最も影響力の高いギルドが各セルを支配する。
// 2位の影響力が1位の75%を超えるセルは「争奪中」としてフラグを立てる。
pub fn calc(radius : Int) -> Unit {
  num_cells.val = 0
  let mut q = -radius
  while q <= radius {
    let r_lo = Int::max(-radius, -q - radius)
    let r_hi = Int::min(radius, -q + radius)
    let mut r = r_lo
    while r <= r_hi {
      if num_cells.val < max_cells {
        let mut best_inf = 0
        let mut second_inf = 0
        let mut best_g = -1
        let mut gi = 0
        while gi < num_guilds.val {
          let d = hex_dist(q, r, guild_q[gi], guild_r[gi])
          let denom = Int::max(1, d * d)
          let inf = guild_cp[gi] * 1000 / denom
          if inf > best_inf {
            second_inf = best_inf
            best_inf = inf
            best_g = gi
          } else if inf > second_inf {
            second_inf = inf
          }
          gi = gi + 1
        }
        out_q[num_cells.val] = q
        out_r[num_cells.val] = r
        out_guild[num_cells.val] = best_g
        out_intensity[num_cells.val] = if best_inf > 0 {
          Int::min(100, best_inf / 100)
        } else {
          0
        }
        out_contested[num_cells.val] = if best_g >= 0 && second_inf * 100 > best_inf * 75 {
          1
        } else {
          0
        }
        num_cells.val = num_cells.val + 1
      }
      r = r + 1
    }
    q = q + 1
  }
}

pub fn get_cell_count() -> Int {
  num_cells.val
}

pub fn get_cell_q(i : Int) -> Int {
  out_q[i]
}

pub fn get_cell_r(i : Int) -> Int {
  out_r[i]
}

// ギルドインデックス（0始まり）を返す。未支配の場合は -1
pub fn get_cell_guild(i : Int) -> Int {
  out_guild[i]
}

// 0〜100 の強度値を返す
pub fn get_cell_intensity(i : Int) -> Int {
  out_intensity[i]
}

// 争奪中（2ギルドの影響力がほぼ拮抗）なら 1、そうでなければ 0 を返す
pub fn get_cell_contested(i : Int) -> Int {
  out_contested[i]
}
