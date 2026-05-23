import { useEffect, useMemo } from "react";
import { useMachine } from "@xstate/react";
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
import {
  sampleBattleResult,
  sampleCurrentPet,
  sampleOpponents,
} from "../../features/pet/sampleData";
import styles from "./PetPage.module.css";

interface PetPageProps {
  onNavigate: (path: string) => void;
}

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
  const [snapshot, send] = useMachine(petPageMachine);
  const {
    data,
    opponents,
    selectedOpponentId,
    battleResult,
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
      const petName = grantedPet.attribute === "go" ? "Gopher君" : grantedPet.attribute;
      send({ type: "NOTICE", message: `${guildName}の相棒「${petName}」が仲間になった！` });
    }
  }, [send]);

  useEffect(() => {
    let isMounted = true;

    fetchMyPets()
      .then((result) => {
        if (!isMounted) return;
        send({ type: "LOAD_SUCCESS", data: result });
      })
      .catch((error) => {
        if (!isMounted) return;
        console.error("failed to fetch pets", error);
        if (import.meta.env.DEV) {
          send({
            type: "LOAD_SUCCESS",
            data: {
              cpBalance: 120,
              currentGuildPet: sampleCurrentPet,
              pets: [sampleCurrentPet],
            },
            statusMessage: "API未接続のため、画面確認用サンプルを表示しています。",
          });
          return;
        }
        send({ type: "LOAD_FAILURE", message: "ペット情報を取得できませんでした。" });
      });

    fetchBattleOpponents()
      .then((result) => {
        if (!isMounted) return;
        send({ type: "OPPONENTS_SUCCESS", opponents: result });
      })
      .catch((error) => {
        if (!isMounted) return;
        console.info("battle opponents are not available yet", error);
        if (import.meta.env.DEV) {
          send({ type: "OPPONENTS_SUCCESS", opponents: sampleOpponents });
          return;
        }
        send({ type: "OPPONENTS_FAILURE" });
      });

    return () => {
      isMounted = false;
    };
  }, [send]);

  const currentPet = data?.currentGuildPet ?? null;
  const selectedOpponent = useMemo(
    () => opponents.find((opponent) => opponent.userId === selectedOpponentId) ?? null,
    [opponents, selectedOpponentId],
  );

  const train = async (stat: PetTrainingStat) => {
    if (!currentPet || isTraining) return;
    const training = PET_TRAINING_COSTS[stat];
    if ((data?.cpBalance ?? 0) < training.cost) {
      send({ type: "INSUFFICIENT_CP" });
      return;
    }

    send({ type: "TRAIN", stat });
    try {
      const result = await trainPet(currentPet.id, stat);
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
    if (!selectedOpponent || isBattling) return;

    send({ type: "BATTLE" });
    try {
      if (import.meta.env.DEV && selectedOpponent.userId.startsWith("sample_")) {
        await new Promise((resolve) => window.setTimeout(resolve, 650));
        send({ type: "BATTLE_SUCCESS", result: sampleBattleResult });
        return;
      }
      const result = await startPetBattle(selectedOpponent.userId);
      send({ type: "BATTLE_SUCCESS", result });
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
                  onClick={() => send({ type: "SELECT_OPPONENT", userId: opponent.userId })}
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
