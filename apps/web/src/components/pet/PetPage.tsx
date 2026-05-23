import { useEffect, useMemo, useRef, useState, type CSSProperties } from "react";
import { useMachine } from "@xstate/react";
import { motion, type Variants } from "framer-motion";
import { SPRITE_ASSETS } from "../../constants/assets";
import { PATHS } from "../../constants/paths";
import { ApiError } from "../../lib/api/client";
import { GopherSprite } from "../shared/GopherSprite";
import {
  fetchBattleOpponents,
  fetchMyPets,
  PET_TRAINING_COSTS,
  startPetBattle,
  trainPet,
  type GrantedPet,
  type PetSummary,
  type PetTrainingStat,
} from "../../features/pet/api";
import { consumeGrantedPet } from "../../features/pet/guildGrant";
import { petPageMachine } from "../../features/pet/petPageMachine";
import { buildSampleBattleResult } from "../../features/pet/battleReplay";
import { saveBattleSession } from "../../features/pet/battleSession";
import {
  sampleBattleResult,
  sampleCurrentPet,
  sampleOwnedPets,
  sampleOpponents,
} from "../../features/pet/sampleData";
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
  return pet.attribute.toLowerCase() === "go" ? "Gopher" : pet.name;
}

function grantedPetAttributeLabel(attribute: string) {
  const normalized = attribute.toLowerCase();
  const labels: Record<string, string> = {
    go: "Go",
    rust: "Rust",
    python: "Python",
    java: "Java",
    typescript: "TypeScript",
    haskell: "Haskell",
    zig: "Zig",
  };
  return labels[normalized] ?? attribute;
}

const petPortraits: Record<string, { label: string; tone: string }> = {
  Rust: { label: "Fe", tone: "#ff9f6e" },
  TypeScript: { label: "TS", tone: "#6bb7ff" },
  Python: { label: "Py", tone: "#ffd966" },
  Java: { label: "Jv", tone: "#ff7b7b" },
  Haskell: { label: "λ", tone: "#b89cff" },
  Zig: { label: "Zg", tone: "#f7a541" },
};

const SPRITE_FRAME_WIDTH = 192;
const SPRITE_FRAME_HEIGHT = 208;
const SPRITE_COLUMNS = 8;
const SPRITE_ROWS = 9;
const PYTHON_IDLE_FRAMES = 6;
const PET_SPRITE_DISPLAY_WIDTH = 132;
const PET_SPRITE_DISPLAY_HEIGHT = 143;

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

function levelProgress(pet: PetSummary) {
  const expToNextLevel = Math.max(100, pet.level * 100);
  const currentExp = Math.max(0, Math.min(pet.exp, expToNextLevel));
  return {
    currentExp,
    expToNextLevel,
    percentage: Math.round((currentExp / expToNextLevel) * 100),
  };
}

function apiWaitingMessage(error: unknown, fallback: string) {
  if (error instanceof ApiError && error.status === 404) {
    return "この操作は API 実装待ちです。画面と client の接続口だけ先に用意しています。";
  }
  return fallback;
}

export function PetPage({ onNavigate }: PetPageProps) {
  const isMountedRef = useRef(false);
  const [grantedPetNotice, setGrantedPetNotice] = useState<GrantedPet | null>(null);
  const [snapshot, send] = useMachine(petPageMachine);
  const { data, selectedPetId, opponents, selectedOpponentId, statusMessage, trainingStat } =
    snapshot.context;
  const isLoading = snapshot.matches("loading");
  const isTraining = snapshot.matches("training");
  const isBattling = snapshot.matches("battling");

  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    const grantedPet = consumeGrantedPet();
    if (grantedPet) {
      setGrantedPetNotice(grantedPet);
    }
  }, []);

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
          : await startPetBattle(selectedPet.id, selectedOpponent.petId);
      if (!isMountedRef.current) return;
      saveBattleSession({ playerPet: selectedPet, opponent: selectedOpponent, result });
      send({ type: "BATTLE_SUCCESS", result });
      onNavigate(PATHS.BATTLE);
    } catch (error) {
      if (!isMountedRef.current) return;
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

  const startDemoBattle = () => {
    const opponent = sampleOpponents[0];
    if (!opponent) return;

    saveBattleSession({ playerPet: sampleCurrentPet, opponent, result: sampleBattleResult });
    onNavigate(PATHS.BATTLE);
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

        {grantedPetNotice && <GrantedPetCelebration grantedPet={grantedPetNotice} />}

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
                      {selectedPet.id === data.currentGuildPet?.id && (
                        <span className={styles.chip}>CURRENT GUILD</span>
                      )}
                    </div>
                    <PetStatusPanel pet={selectedPet} cpBalance={data.cpBalance} />
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
                        <span className={styles.actionLabel}>
                          {trainingStat === stat ? "TRAINING..." : training.label}
                        </span>
                        <span className={styles.actionBoost}>+{training.amount}</span>
                        <span className={styles.actionCost}>{training.cost} CP</span>
                      </button>
                    );
                  })}
                </div>

                {statusMessage && (
                  <div className={styles.trainingFeedback} key={statusMessage} role="status">
                    <span className={styles.trainingFeedbackLabel}>BOOST RESULT</span>
                    <p className={styles.trainingFeedbackText}>{statusMessage}</p>
                  </div>
                )}

                <PetBattleRecordPreview />
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
              {import.meta.env.DEV && (
                <button className={styles.demoBattleButton} type="button" onClick={startDemoBattle}>
                  DEMO BATTLE
                </button>
              )}
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

function GrantedPetCelebration({ grantedPet }: { grantedPet: GrantedPet }) {
  const attribute = grantedPetAttributeLabel(grantedPet.attribute);
  const guildName = grantedPet.guildId === "guild_go" ? "Go Guild" : "New Guild";
  const pet: PetSummary = {
    id: grantedPet.id,
    guildId: grantedPet.guildId,
    guildName,
    name: attribute.toLowerCase() === "go" ? "Gopher" : attribute,
    species: attribute.toLowerCase(),
    attribute,
    level: 1,
    exp: 0,
    maxHp: 1,
    power: 1,
    guard: 1,
    speed: 1,
    acquiredAt: grantedPet.createdAt,
  };

  return (
    <motion.aside
      animate={{ opacity: [0, 1, 1, 0], y: [16, 0, 0, -12], scale: [0.98, 1, 1, 0.99] }}
      className={styles.grantToast}
      initial={{ opacity: 0, y: 18, scale: 0.98 }}
      transition={{ duration: 5.2, times: [0, 0.16, 0.82, 1], ease: "easeOut" }}
    >
      <div className={styles.grantPortrait} aria-hidden="true">
        <PetPortrait pet={pet} />
      </div>
      <div>
        <span className={styles.grantLabel}>NEW COMPANION</span>
        <h3 className={styles.grantTitle}>{petDisplayName(pet)} joined!</h3>
        <p className={styles.grantText}>{guildName} の相棒が仲間になりました。</p>
      </div>
    </motion.aside>
  );
}

function PetBattleRecordPreview() {
  const metrics = [
    { label: "BATTLES", value: "--" },
    { label: "WINS", value: "--" },
    { label: "STREAK", value: "--" },
  ];

  return (
    <section className={styles.recordPanel} aria-labelledby="pet-record-title">
      <div className={styles.recordHeader}>
        <div>
          <h3 className={styles.recordTitle} id="pet-record-title">
            CAREER LOG
          </h3>
          <p className={styles.recordCaption}>戦績データ連携予定</p>
        </div>
        <span className={styles.recordBadge}>COMING SOON</span>
      </div>
      <div className={styles.recordMetrics}>
        {metrics.map((metric) => (
          <div className={styles.recordMetric} key={metric.label}>
            <span className={styles.recordMetricLabel}>{metric.label}</span>
            <strong className={styles.recordMetricValue}>{metric.value}</strong>
          </div>
        ))}
      </div>
    </section>
  );
}

function PetStatusPanel({ pet, cpBalance }: { pet: PetSummary; cpBalance: number }) {
  const progress = levelProgress(pet);

  return (
    <div className={styles.statusGrid} aria-label="ペット育成ステータス">
      <div className={`${styles.statusCard} ${styles.levelCard}`}>
        <div className={styles.statusHeader}>
          <span className={styles.statusLabel}>CURRENT LV</span>
          <span className={styles.levelBadge}>Lv {pet.level}</span>
        </div>
        <div className={styles.expTrack} aria-label={`次のレベルまで ${progress.percentage}%`}>
          <span className={styles.expFill} style={{ width: `${progress.percentage}%` }} />
        </div>
        <p className={styles.statusHelp}>
          NEXT {progress.expToNextLevel - progress.currentExp} EXP
        </p>
      </div>

      <div className={`${styles.statusCard} ${styles.cpCard}`}>
        <span className={styles.statusLabel}>OWNED CP</span>
        <strong className={styles.cpValue}>{cpBalance}</strong>
        <span className={styles.statusHelp}>TRAINING RESOURCE</span>
      </div>
    </div>
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
  if (pet.attribute.toLowerCase() === "go") {
    return <GopherSprite />;
  }
  if (pet.attribute.toLowerCase() === "python") {
    return <PythonPetSprite />;
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

function PythonPetSprite() {
  const scale = PET_SPRITE_DISPLAY_WIDTH / SPRITE_FRAME_WIDTH;
  const displaySheetWidth = Math.round(SPRITE_FRAME_WIDTH * SPRITE_COLUMNS * scale);
  const displaySheetHeight = Math.round(SPRITE_FRAME_HEIGHT * SPRITE_ROWS * scale);
  const frameStep = Math.round(SPRITE_FRAME_WIDTH * scale);
  const totalMoveX = frameStep * PYTHON_IDLE_FRAMES;

  return (
    <motion.div
      animate={{ backgroundPositionX: ["0px", `-${totalMoveX}px`] }}
      className={styles.petSprite}
      transition={{
        duration: 0.9,
        repeat: Infinity,
        ease: steppedEase(PYTHON_IDLE_FRAMES),
      }}
      style={{
        width: `${PET_SPRITE_DISPLAY_WIDTH}px`,
        height: `${PET_SPRITE_DISPLAY_HEIGHT}px`,
        backgroundImage: `url(${SPRITE_ASSETS.PYTHON})`,
        backgroundRepeat: "no-repeat",
        backgroundSize: `${displaySheetWidth}px ${displaySheetHeight}px`,
      }}
    />
  );
}
