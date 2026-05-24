// ── hex_dist のテスト ──────────────────────────────────────────────────────────

test "hex_dist: 同じセルは0" {
  assert_eq(hex_dist(0, 0, 0, 0), 0)
}

test "hex_dist: 隣接セルは1" {
  assert_eq(hex_dist(0, 0, 1, 0), 1)
  assert_eq(hex_dist(0, 0, 0, 1), 1)
  assert_eq(hex_dist(0, 0, -1, 1), 1)
}

test "hex_dist: 距離2" {
  assert_eq(hex_dist(0, 0, 2, 0), 2)
  assert_eq(hex_dist(0, 0, 1, 1), 2)
}

test "hex_dist: 対称性" {
  assert_eq(hex_dist(3, -1, 0, 0), hex_dist(0, 0, 3, -1))
}

// ── calc のテスト ─────────────────────────────────────────────────────────────

test "calc: 半径1のグリッドは7セル" {
  set_guild_count(1)
  set_guild(0, 0, 0, 100)
  calc(1)
  assert_eq(get_cell_count(), 7)
}

test "calc: ホームセルは自ギルドが支配" {
  set_guild_count(1)
  set_guild(0, 0, 0, 100)
  calc(3)
  let mut found = false
  let mut i = 0
  while i < get_cell_count() {
    if get_cell_q(i) == 0 && get_cell_r(i) == 0 {
      assert_eq(get_cell_guild(i), 0)
      found = true
    }
    i = i + 1
  }
  assert_eq(found, true)
}

test "calc: 強度はホームセルで最大（100）" {
  set_guild_count(1)
  set_guild(0, 0, 0, 100)
  calc(3)
  let mut home_intensity = 0
  let mut i = 0
  while i < get_cell_count() {
    if get_cell_q(i) == 0 && get_cell_r(i) == 0 {
      home_intensity = get_cell_intensity(i)
    }
    i = i + 1
  }
  assert_eq(home_intensity, 100)
}

test "calc: 2ギルドが隣接すると争奪セルが存在する" {
  set_guild_count(2)
  set_guild(0, -3, 0, 100)
  set_guild(1, 3, 0, 100)
  calc(5)
  let mut contested_count = 0
  let mut i = 0
  while i < get_cell_count() {
    if get_cell_contested(i) == 1 {
      contested_count = contested_count + 1
    }
    i = i + 1
  }
  assert_eq(contested_count > 0, true)
}

test "calc: 2ギルドが離れていると両ギルドの支配セルが存在する" {
  set_guild_count(2)
  set_guild(0, -4, 0, 100)
  set_guild(1, 4, 0, 100)
  calc(5)
  let mut guild0_found = false
  let mut guild1_found = false
  let mut i = 0
  while i < get_cell_count() {
    if get_cell_guild(i) == 0 { guild0_found = true }
    if get_cell_guild(i) == 1 { guild1_found = true }
    i = i + 1
  }
  assert_eq(guild0_found, true)
  assert_eq(guild1_found, true)
}
