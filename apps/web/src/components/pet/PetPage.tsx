import { useEffect, useMemo, type CSSProperties } from "react";
import { useMachine } from "@xstate/react";
import { motion, type Variants } from "framer-motion";
import { PATHS } from "../../constants/paths";
import { ApiError } from "../../lib/api/client";
import { GopherSprite } from "../shared/GopherSprite";
import {
  fetchBattleOpponents,
  fetchMyPets,
  PET_TRAINING_COSTS,
  startPetBattle,
  trainPet,
  type PetSummary,
  type PetTrainingStat,
} from "../../features/pet/api";
import { consumeGrantedPet } from "../../features/pet/guildGrant";
import { petPageMachine } from "../../features/pet/petPageMachine";
import { buildSampleBattleResult } from "../../features/pet/battleReplay";
import { saveBattleSession } from "../../features/pet/battleSession";
import { sampleCurrentPet, sampleOwnedPets, sampleOpponents } from "../../features/pet/sampleData";
import { steppedEase } from "../../lib/animationUtils";
import styles from "./PetPage.module.css";

interface PetPageProps {
  onNavigate: (path: string) => void;
}

let petPageBootstrapPromise: Promise<{
  myPets: PromiseSettledResult<Awaited<ReturnType<typeof fetchMyPets>>>;
  opponents: PromiseSettledResult<Awaited<ReturnType<typeof fetchBattleOpponents>>>;
}> | null = null;

async function fetchPetPageBootstrap() {
  petPageBootstrapPromise ??= Promise.allSettled([fetchMyPets(), fetchBattleOpponents()]).then(
    ([myPets, opponents]) => ({ myPets, opponents }),
  );
  return petPageBootstrapPromise;
}

function petDisplayName(pet: PetSummary | null | undefined) {
  if (!pet) return "相棒未選択";
  return pet.attribute === "Go" ? "Gopher" : pet.name;
}

const petPortraits: Record<string, { label: string; tone: string }> = {
  Rust: { label: "Fe", tone: "#ff9f6e" },
  TypeScript: { label: "TS", tone: "#6bb7ff" },
  Python: { label: "Py", tone: "#ffd966" },
  Java: { label: "Jv", tone: "#ff7b7b" },
  Haskell: { label: "λ", tone: "#b89cff" },
  Zig: { label: "Zg", tone: "#f7a541" },
};

const pageVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      duration: 0.24,
      ease: steppedEase(4),
      staggerChildren: 0.08,
      delayChildren: 0.08,
    },
  },
};

const panelVariants: Variants = {
  hidden: { opacity: 0, y: 18, scale: 0.985 },
  visible: {
    opacity: 1,
    y: 0,
    scale: 1,
    transition: { duration: 0.34, ease: steppedEase(6) },
  },
};

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
  const [snapshot, send] = useMachine(petPageMachine);
  const {
    data,
    selectedPetId,
    opponents,
    selectedOpponentId,
    statusMessage,
    noticeMessage,
    trainingStat,
  } = snapshot.context;
  const isLoading = snapshot.matches("loading");
  const isTraining = snapshot.matches("training");
  const isBattling = snapshot.matches("battling");

  useEffect(() => {
    const grantedPet = consumeGrantedPet();
    if (grantedPet) {
      const guildName = grantedPet.guildId === "guild_go" ? "Goギルド" : "所属ギルド";
      const petName = grantedPet.attribute === "go" ? "Gopher" : grantedPet.attribute;
      send({ type: "NOTICE", message: `${guildName}の相棒「${petName}」が仲間になった！` });
    }
  }, [send]);

  useEffect(() => {
    let isMounted = true;

    fetchPetPageBootstrap()
      .then(({ myPets, opponents }) => {
        if (!isMounted) {
          return;
        }
        if (myPets.status === "rejected" || opponents.status === "rejected") {
          petPageBootstrapPromise = null;
        }

        if (myPets.status === "fulfilled") {
          send({ type: "LOAD_SUCCESS", data: myPets.value });
        } else if (import.meta.env.DEV) {
          console.error("failed to fetch pets", myPets.reason);
          send({
            type: "LOAD_SUCCESS",
            data: {
              cpBalance: 120,
              currentGuildPet: sampleCurrentPet,
              pets: sampleOwnedPets,
            },
            statusMessage: "API未接続のため、画面確認用サンプルを表示しています。",
          });
        } else {
          console.error("failed to fetch pets", myPets.reason);
          send({ type: "LOAD_FAILURE", message: "ペット情報を取得できませんでした。" });
        }

        if (opponents.status === "fulfilled") {
          send({ type: "OPPONENTS_SUCCESS", opponents: opponents.value });
        } else if (import.meta.env.DEV) {
          console.info("battle opponents are not available yet", opponents.reason);
          send({ type: "OPPONENTS_SUCCESS", opponents: sampleOpponents });
        } else {
          send({ type: "OPPONENTS_FAILURE" });
        }
      })
      .catch((error: unknown) => {
        if (!isMounted) return;
        console.error("failed to fetch pet page data", error);
        petPageBootstrapPromise = null;
        if (import.meta.env.DEV) {
          send({
            type: "LOAD_SUCCESS",
            data: {
              cpBalance: 120,
              currentGuildPet: sampleCurrentPet,
              pets: sampleOwnedPets,
            },
            statusMessage: "API未接続のため、画面確認用サンプルを表示しています。",
          });
          send({ type: "OPPONENTS_SUCCESS", opponents: sampleOpponents });
          return;
        }
        send({ type: "LOAD_FAILURE", message: "ペット情報を取得できませんでした。" });
        send({ type: "OPPONENTS_FAILURE" });
      });

    return () => {
      isMounted = false;
    };
  }, [send]);

  const selectedPet = useMemo(
    () => data?.pets.find((pet) => pet.id === selectedPetId) ?? data?.currentGuildPet ?? null,
    [data, selectedPetId],
  );
  const selectedOpponent = useMemo(
    () => opponents.find((opponent) => opponent.userId === selectedOpponentId) ?? null,
    [opponents, selectedOpponentId],
  );

  const train = async (stat: PetTrainingStat) => {
    if (!selectedPet || isTraining) return;
    const training = PET_TRAINING_COSTS[stat];
    if ((data?.cpBalance ?? 0) < training.cost) {
      send({ type: "INSUFFICIENT_CP" });
      return;
    }

    send({ type: "TRAIN", stat });
    try {
      const result = await trainPet(selectedPet.id, stat);
      send({
        type: "TRAIN_SUCCESS",
        result,
        message: `${PET_TRAINING_COSTS[result.increasedStat].label} が ${result.increasedBy} 上がった！ CP: ${result.cpBefore} → ${result.cpAfter}`,
      });
    } catch (error) {
      console.error("failed to train pet", error);
      send({
        type: "TRAIN_FAILURE",
        message: apiWaitingMessage(
          error,
          "育成に失敗しました。少し時間を置いて再度お試しください。",
        ),
      });
    }
  };

  const battle = async () => {
    if (!selectedOpponent || !selectedPet || isBattling) return;

    send({ type: "BATTLE" });
    try {
      const result =
        import.meta.env.DEV && selectedOpponent.userId.startsWith("sample_")
          ? await new Promise<ReturnType<typeof buildSampleBattleResult>>((resolve) => {
              window.setTimeout(
                () => resolve(buildSampleBattleResult(selectedPet, selectedOpponent)),
                650,
              );
            })
          : await startPetBattle(selectedOpponent.userId);
      saveBattleSession({ playerPet: selectedPet, opponent: selectedOpponent, result });
      send({ type: "BATTLE_SUCCESS", result });
      onNavigate(PATHS.BATTLE);
    } catch (error) {
      console.error("failed to start pet battle", error);
      send({
        type: "BATTLE_FAILURE",
        message: apiWaitingMessage(
          error,
          "バトルを開始できませんでした。少し時間を置いて再度お試しください。",
        ),
      });
    }
  };

  return (
    <motion.main
      animate="visible"
      className={styles.screen}
      initial="hidden"
      variants={pageVariants}
    >
      <motion.div className={styles.shell} variants={pageVariants}>
        <motion.header className={styles.header} variants={panelVariants}>
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
        </motion.header>

        {noticeMessage && <div className={styles.notice}>{noticeMessage}</div>}
        {statusMessage && <div className={styles.notice}>{statusMessage}</div>}

        <motion.div className={styles.layout} variants={pageVariants}>
          <motion.section
            className={`${styles.panel} ${styles.trainingPanel}`}
            aria-labelledby="current-pet-title"
            variants={panelVariants}
          >
            <h2 className={styles.panelTitle} id="current-pet-title">
              TRAINING PET
            </h2>
            {isLoading && <p className={styles.message}>Loading...</p>}
            {!isLoading && !data && <p className={styles.message}>ペット情報を表示できません。</p>}
            {data && !selectedPet && (
              <p className={styles.message}>
                所持ペットはいません。ギルドに加入すると相棒が配布されます。
              </p>
            )}
            {selectedPet && data && (
              <>
                <div className={styles.currentPet}>
                  <div className={styles.spriteStage} aria-hidden="true">
                    <PetPortrait pet={selectedPet} />
                  </div>
                  <div>
                    <h2 className={styles.petName}>{petDisplayName(selectedPet)}</h2>
                    <div className={styles.meta}>
                      <span className={styles.chip}>{selectedPet.guildName} Guild</span>
                      <span className={styles.chip}>{selectedPet.attribute}</span>
                      <span className={styles.chip}>Lv {selectedPet.level}</span>
                      {selectedPet.id === data.currentGuildPet?.id && (
                        <span className={styles.chip}>CURRENT GUILD</span>
                      )}
                      <span className={styles.chip}>CP {data.cpBalance}</span>
                    </div>
                    <PetStats pet={selectedPet} />
                  </div>
                </div>

                <div className={styles.trainingGrid}>
                  {(Object.keys(PET_TRAINING_COSTS) as PetTrainingStat[]).map((stat) => {
                    const training = PET_TRAINING_COSTS[stat];
                    const lacksCP = data.cpBalance < training.cost;
                    return (
                      <button
                        className={styles.actionButton}
                        disabled={isTraining || lacksCP}
                        key={stat}
                        type="button"
                        onClick={() => void train(stat)}
                      >
                        {trainingStat === stat
                          ? "TRAINING..."
                          : `${training.label} +${training.amount} / ${training.cost} CP`}
                      </button>
                    );
                  })}
                </div>
              </>
            )}
          </motion.section>

          <motion.section
            className={`${styles.panel} ${styles.ownedPanel}`}
            aria-labelledby="owned-pets-title"
            variants={panelVariants}
          >
            <h2 className={styles.panelTitle} id="owned-pets-title">
              OWNED PETS
            </h2>
            <div className={styles.list}>
              {data?.pets.length === 0 && <p className={styles.message}>所持ペットはいません。</p>}
              {data?.pets.map((pet) => (
                <button
                  aria-selected={pet.id === selectedPetId}
                  className={styles.petListItem}
                  key={pet.id}
                  type="button"
                  onClick={() => send({ type: "SELECT_PET", petId: pet.id })}
                >
                  <div>
                    <h3 className={styles.itemTitle}>{petDisplayName(pet)}</h3>
                    <p className={styles.itemText}>
                      {pet.guildName} / HP {pet.maxHp} / Power {pet.power}
                    </p>
                  </div>
                  <div className={styles.petListBadges}>
                    {pet.id === data.currentGuildPet?.id && (
                      <span className={styles.chip}>GUILD</span>
                    )}
                    {pet.id === selectedPetId && <span className={styles.chip}>SELECTED</span>}
                    <span className={styles.chip}>Lv {pet.level}</span>
                  </div>
                </button>
              ))}
            </div>
          </motion.section>

          <motion.section
            className={`${styles.panel} ${styles.battlePanel}`}
            aria-labelledby="battle-title"
            variants={panelVariants}
          >
            <h2 className={styles.panelTitle} id="battle-title">
              AUTO BATTLE
            </h2>
            <div className={styles.opponentList}>
              {opponents.length === 0 && (
                <p className={styles.message}>対戦相手候補はまだありません。</p>
              )}
              {opponents.map((opponent) => (
                <button
                  aria-selected={opponent.userId === selectedOpponentId}
                  className={styles.opponentItem}
                  key={opponent.userId}
                  type="button"
                  onClick={() => send({ type: "SELECT_OPPONENT", userId: opponent.userId })}
                >
                  <span className={styles.itemTitle}>{opponent.playerName}</span>
                  <span className={styles.itemText}>
                    {petDisplayName(opponent.pet)} / Lv {opponent.pet.level}
                  </span>
                </button>
              ))}
            </div>
            <div className={styles.battleActions}>
              <button
                className={styles.battleButton}
                disabled={!selectedOpponent || !selectedPet || isBattling}
                type="button"
                onClick={() => void battle()}
              >
                {isBattling ? "MATCHING..." : "START BATTLE"}
              </button>
            </div>
          </motion.section>
        </motion.div>
      </motion.div>
    </motion.main>
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

function PetPortrait({ pet }: { pet: PetSummary }) {
  if (pet.attribute === "Go") {
    return <GopherSprite />;
  }

  const portrait = petPortraits[pet.attribute] ?? {
    label: pet.attribute.slice(0, 2).toUpperCase(),
    tone: "#74f7a1",
  };

  return (
    <div
      className={styles.placeholderPortrait}
      style={{ "--pet-tone": portrait.tone } as CSSProperties}
    >
      <span className={styles.placeholderAura} />
      <span className={styles.placeholderFace}>{portrait.label}</span>
      <span className={styles.placeholderName}>{petDisplayName(pet)}</span>
    </div>
  );
}
