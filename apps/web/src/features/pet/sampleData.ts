import type { BattleOpponent, BattleResult, PetSummary } from "./api";

export const sampleCurrentPet: PetSummary = {
  id: "sample_pet_go",
  guildId: "guild_go",
  guildName: "Go",
  name: "Gopher",
  species: "gopher",
  attribute: "Go",
  level: 4,
  exp: 40,
  maxHp: 35,
  power: 6,
  guard: 5,
  speed: 7,
  acquiredAt: "2026-05-23T00:00:00Z",
};

export const sampleOpponents: BattleOpponent[] = [
  {
    userId: "sample_user_rust",
    playerName: "FerrisBlade",
    pet: {
      id: "sample_pet_rust",
      guildId: "guild_rust",
      guildName: "Rust",
      name: "Ferris",
      species: "crab",
      attribute: "Rust",
      level: 3,
      exp: 20,
      maxHp: 32,
      power: 7,
      guard: 6,
      speed: 4,
      acquiredAt: "2026-05-23T00:00:00Z",
    },
  },
  {
    userId: "sample_user_python",
    playerName: "PyRunner",
    pet: {
      id: "sample_pet_python",
      guildId: "guild_python",
      guildName: "Python",
      name: "Py",
      species: "python",
      attribute: "Python",
      level: 5,
      exp: 80,
      maxHp: 40,
      power: 5,
      guard: 4,
      speed: 8,
      acquiredAt: "2026-05-23T00:00:00Z",
    },
  },
];

export const sampleOwnedPets: PetSummary[] = [
  sampleCurrentPet,
  {
    id: "sample_pet_typescript",
    guildId: "guild_typescript",
    guildName: "TypeScript",
    name: "Scriptie",
    species: "typescript",
    attribute: "TypeScript",
    level: 2,
    exp: 15,
    maxHp: 30,
    power: 6,
    guard: 4,
    speed: 6,
    acquiredAt: "2026-05-23T00:00:00Z",
  },
  {
    id: "sample_pet_rust_owned",
    guildId: "guild_rust",
    guildName: "Rust",
    name: "Ferris",
    species: "crab",
    attribute: "Rust",
    level: 3,
    exp: 25,
    maxHp: 45,
    power: 7,
    guard: 8,
    speed: 3,
    acquiredAt: "2026-05-23T00:00:00Z",
  },
];

export const sampleBattleResult: BattleResult = {
  result: "win",
  turns: [
    {
      turn: 1,
      actorPetId: "sample_pet_go",
      targetPetId: "sample_pet_rust",
      damage: 7,
      message: "Gopher君の先制攻撃！ Ferris に 7 ダメージ。",
    },
    {
      turn: 2,
      actorPetId: "sample_pet_rust",
      targetPetId: "sample_pet_go",
      damage: 4,
      message: "Ferris の反撃。Gopher君は 4 ダメージを受けた。",
    },
    {
      turn: 3,
      actorPetId: "sample_pet_go",
      targetPetId: "sample_pet_rust",
      damage: 9,
      message: "Gopher君の会心アタック！ 勝負あり。",
    },
  ],
};
