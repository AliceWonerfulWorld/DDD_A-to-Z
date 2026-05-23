import { useEffect, useMemo, useState } from "react";
import { PATHS } from "../../constants/paths";
import { ApiError } from "../../lib/api/client";
import { GopherSprite } from "../shared/GopherSprite";
import {
  fetchBattleOpponents,
  fetchMyPets,
  PET_TRAINING_COSTS,
  startPetBattle,
  trainPet,
  type BattleOpponent,
  type BattleResult,
  type MyPetsResponse,
  type PetSummary,
  type PetTrainingStat,
} from "../../features/pet/api";
import { consumeGrantedPet } from "../../features/pet/guildGrant";
import styles from "./PetPage.module.css";

interface PetPageProps {
  onNavigate: (path: string) => void;
}

const sampleCurrentPet: PetSummary = {
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

const sampleOpponents: BattleOpponent[] = [
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

const sampleBattleResult: BattleResult = {
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

function petDisplayName(pet: PetSummary | null | undefined) {
  if (!pet) return "相棒未選択";
  return pet.attribute === "Go" ? "Gopher君" : pet.name;
}

function statValue(pet: PetSummary, stat: PetTrainingStat) {
  if (stat === "hp") return pet.maxHp;
  return pet[stat];
}

function apiWaitingMessage(error: unknown, fallback: string) {
  if (error instanceof ApiError && error.status === 404) {
    return "この操作は API 実装待ちです。画面と client の接続口だけ先に用意しています。";
  }
  return fallback;
}

export function PetPage({ onNavigate }: PetPageProps) {
  const [data, setData] = useState<MyPetsResponse | null>(null);
  const [opponents, setOpponents] = useState<BattleOpponent[]>([]);
  const [selectedOpponentId, setSelectedOpponentId] = useState<string | null>(null);
  const [battleResult, setBattleResult] = useState<BattleResult | null>(null);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [noticeMessage, setNoticeMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isTraining, setIsTraining] = useState<PetTrainingStat | null>(null);
  const [isBattling, setIsBattling] = useState(false);

  useEffect(() => {
    const grantedPet = consumeGrantedPet();
    if (grantedPet) {
      const guildName = grantedPet.guildId === "guild_go" ? "Goギルド" : "所属ギルド";
      const petName = grantedPet.attribute === "go" ? "Gopher君" : grantedPet.attribute;
      setNoticeMessage(`${guildName}の相棒「${petName}」が仲間になった！`);
    }
  }, []);

  useEffect(() => {
    let isMounted = true;

    fetchMyPets()
      .then((result) => {
        if (!isMounted) return;
        setData(result);
      })
      .catch((error) => {
        if (!isMounted) return;
        console.error("failed to fetch pets", error);
        if (import.meta.env.DEV) {
          setData({
            cpBalance: 120,
            currentGuildPet: sampleCurrentPet,
            pets: [sampleCurrentPet],
          });
          setStatusMessage("API未接続のため、画面確認用サンプルを表示しています。");
          return;
        }
        setStatusMessage("ペット情報を取得できませんでした。");
      })
      .finally(() => {
        if (isMounted) setIsLoading(false);
      });

    fetchBattleOpponents()
      .then((result) => {
        if (!isMounted) return;
        setOpponents(result);
        setSelectedOpponentId(result[0]?.userId ?? null);
      })
      .catch((error) => {
        if (!isMounted) return;
        console.info("battle opponents are not available yet", error);
        if (import.meta.env.DEV) {
          setOpponents(sampleOpponents);
          setSelectedOpponentId(sampleOpponents[0]?.userId ?? null);
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  const currentPet = data?.currentGuildPet ?? null;
  const selectedOpponent = useMemo(
    () => opponents.find((opponent) => opponent.userId === selectedOpponentId) ?? null,
    [opponents, selectedOpponentId],
  );

  const train = async (stat: PetTrainingStat) => {
    if (!currentPet || isTraining) return;
    const training = PET_TRAINING_COSTS[stat];
    if ((data?.cpBalance ?? 0) < training.cost) {
      setStatusMessage("CPが足りません");
      return;
    }

    setIsTraining(stat);
    setStatusMessage(null);
    try {
      const result = await trainPet(currentPet.id, stat);
      setData((prev) => {
        if (!prev) return prev;
        return {
          ...prev,
          cpBalance: result.cpAfter,
          currentGuildPet: result.pet,
          pets: prev.pets.map((pet) => (pet.id === result.pet.id ? result.pet : pet)),
        };
      });
      setStatusMessage(
        `${PET_TRAINING_COSTS[result.increasedStat].label} が ${result.increasedBy} 上がった！ CP: ${result.cpBefore} → ${result.cpAfter}`,
      );
    } catch (error) {
      console.error("failed to train pet", error);
      setStatusMessage(
        apiWaitingMessage(error, "育成に失敗しました。少し時間を置いて再度お試しください。"),
      );
    } finally {
      setIsTraining(null);
    }
  };

  const battle = async () => {
    if (!selectedOpponent || isBattling) return;

    setIsBattling(true);
    setStatusMessage(null);
    setBattleResult(null);
    try {
      if (import.meta.env.DEV && selectedOpponent.userId.startsWith("sample_")) {
        await new Promise((resolve) => window.setTimeout(resolve, 650));
        setBattleResult(sampleBattleResult);
        return;
      }
      const result = await startPetBattle(selectedOpponent.userId);
      setBattleResult(result);
    } catch (error) {
      console.error("failed to start pet battle", error);
      setStatusMessage(
        apiWaitingMessage(
          error,
          "バトルを開始できませんでした。少し時間を置いて再度お試しください。",
        ),
      );
    } finally {
      setIsBattling(false);
    }
  };

  return (
    <main className={styles.screen}>
      <div className={styles.shell}>
        <header className={styles.header}>
          <div>
            <p className={styles.eyebrow}>PET TERMINAL</p>
            <h1 className={styles.title}>マイペット</h1>
          </div>
          <button
            className={styles.backButton}
            type="button"
            onClick={() => onNavigate(PATHS.HOME)}
          >
            HOME
          </button>
        </header>

        {noticeMessage && <div className={styles.notice}>{noticeMessage}</div>}
        {statusMessage && <div className={styles.notice}>{statusMessage}</div>}

        <div className={styles.layout}>
          <section className={styles.panel} aria-labelledby="current-pet-title">
            <h2 className={styles.panelTitle} id="current-pet-title">
              CURRENT GUILD PET
            </h2>
            {isLoading && <p className={styles.message}>Loading...</p>}
            {!isLoading && !data && <p className={styles.message}>ペット情報を表示できません。</p>}
            {data && !currentPet && (
              <p className={styles.message}>
                現在所属ギルド由来のペットはいません。ギルドに加入すると相棒が配布されます。
              </p>
            )}
            {currentPet && data && (
              <>
                <div className={styles.currentPet}>
                  <div className={styles.spriteStage} aria-hidden="true">
                    <GopherSprite />
                  </div>
                  <div>
                    <h2 className={styles.petName}>{petDisplayName(currentPet)}</h2>
                    <div className={styles.meta}>
                      <span className={styles.chip}>{currentPet.guildName} Guild</span>
                      <span className={styles.chip}>{currentPet.attribute}</span>
                      <span className={styles.chip}>Lv {currentPet.level}</span>
                      <span className={styles.chip}>CP {data.cpBalance}</span>
                    </div>
                    <PetStats pet={currentPet} />
                  </div>
                </div>

                <div className={styles.trainingGrid}>
                  {(Object.keys(PET_TRAINING_COSTS) as PetTrainingStat[]).map((stat) => {
                    const training = PET_TRAINING_COSTS[stat];
                    const lacksCP = data.cpBalance < training.cost;
                    return (
                      <button
                        className={styles.actionButton}
                        disabled={isTraining !== null || lacksCP}
                        key={stat}
                        type="button"
                        onClick={() => void train(stat)}
                      >
                        {training.label} +{training.amount} / {training.cost} CP
                      </button>
                    );
                  })}
                </div>
              </>
            )}
          </section>

          <section className={styles.panel} aria-labelledby="owned-pets-title">
            <h2 className={styles.panelTitle} id="owned-pets-title">
              OWNED PETS
            </h2>
            <div className={styles.list}>
              {data?.pets.length === 0 && <p className={styles.message}>所持ペットはいません。</p>}
              {data?.pets.map((pet) => (
                <article className={styles.petListItem} key={pet.id}>
                  <div>
                    <h3 className={styles.itemTitle}>{petDisplayName(pet)}</h3>
                    <p className={styles.itemText}>
                      {pet.guildName} / HP {pet.maxHp} / Power {pet.power}
                    </p>
                  </div>
                  <span className={styles.chip}>Lv {pet.level}</span>
                </article>
              ))}
            </div>
          </section>

          <section className={styles.panel} aria-labelledby="battle-title">
            <h2 className={styles.panelTitle} id="battle-title">
              AUTO BATTLE
            </h2>
            <div className={styles.list}>
              {opponents.length === 0 && (
                <p className={styles.message}>対戦相手候補はまだありません。</p>
              )}
              {opponents.map((opponent) => (
                <button
                  aria-selected={opponent.userId === selectedOpponentId}
                  className={styles.opponentItem}
                  key={opponent.userId}
                  type="button"
                  onClick={() => setSelectedOpponentId(opponent.userId)}
                >
                  <span className={styles.itemTitle}>{opponent.playerName}</span>
                  <span className={styles.itemText}>
                    {petDisplayName(opponent.pet)} / Lv {opponent.pet.level}
                  </span>
                </button>
              ))}
              <button
                className={styles.actionButton}
                disabled={!selectedOpponent || !currentPet || isBattling}
                type="button"
                onClick={() => void battle()}
              >
                {isBattling ? "BATTLE..." : "BATTLE"}
              </button>
            </div>

            {battleResult && (
              <div className={styles.battleResult}>
                <h3 className={styles.itemTitle}>RESULT: {battleResult.result.toUpperCase()}</h3>
                <div className={styles.list}>
                  {battleResult.turns.map((turn) => (
                    <p className={styles.logItem} key={turn.turn}>
                      {turn.turn}. {turn.message}
                    </p>
                  ))}
                </div>
              </div>
            )}
          </section>
        </div>
      </div>
    </main>
  );
}

function PetStats({ pet }: { pet: PetSummary }) {
  return (
    <div className={styles.statsGrid}>
      {(Object.keys(PET_TRAINING_COSTS) as PetTrainingStat[]).map((stat) => (
        <div className={styles.statCell} key={stat}>
          <span className={styles.statLabel}>{PET_TRAINING_COSTS[stat].label}</span>
          <span className={styles.statValue}>{statValue(pet, stat)}</span>
        </div>
      ))}
    </div>
  );
}
