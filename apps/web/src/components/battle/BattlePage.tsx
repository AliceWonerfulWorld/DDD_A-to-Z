import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useMemo, useState } from "react";
import { PATHS } from "../../constants/paths";
import { GopherSprite } from "../shared/GopherSprite";
import type { BattleOpponent, PetSummary } from "../../features/pet/api";
import { readBattleSession } from "../../features/pet/battleSession";
import {
  buildSampleBattleResult,
  toBattleReplay,
  type BattleReplayTurn,
} from "../../features/pet/battleReplay";
import { sampleCurrentPet, sampleOpponents } from "../../features/pet/sampleData";
import styles from "./BattlePage.module.css";

interface BattlePageProps {
  onNavigate: (path: string) => void;
}

interface BattleCharacterAsset {
  kind: "sprite" | "placeholder";
  label: string;
}

const characterAssets: Record<string, BattleCharacterAsset> = {
  Go: { kind: "sprite", label: "Gopher君" },
  Rust: { kind: "placeholder", label: "Fe" },
  TypeScript: { kind: "placeholder", label: "TS" },
  Python: { kind: "placeholder", label: "Py" },
};

function displayPetName(pet: PetSummary) {
  return pet.attribute === "Go" ? "Gopher君" : pet.name;
}

function fallbackSession() {
  const opponent = sampleOpponents[0]!;
  return {
    playerPet: sampleCurrentPet,
    opponent,
    result: buildSampleBattleResult(sampleCurrentPet, opponent),
  };
}

function clampHP(value: number) {
  return Math.max(0, value);
}

export function BattlePage({ onNavigate }: BattlePageProps) {
  const session = useMemo(() => readBattleSession() ?? fallbackSession(), []);
  const replay = useMemo(
    () => toBattleReplay(session.playerPet, session.opponent, session.result),
    [session],
  );
  const [turnIndex, setTurnIndex] = useState(0);
  const [activeTurn, setActiveTurn] = useState<BattleReplayTurn | null>(null);
  const [logs, setLogs] = useState<string[]>([]);
  const [playerHP, setPlayerHP] = useState(replay.playerPet.maxHp);
  const [enemyHP, setEnemyHP] = useState(replay.opponent.pet.maxHp);
  const isComplete = turnIndex >= replay.turns.length;

  useEffect(() => {
    if (turnIndex >= replay.turns.length) return;

    const timer = window.setTimeout(
      () => {
        const turn = replay.turns[turnIndex];
        if (!turn) return;

        setActiveTurn(turn);
        setLogs((current) => [...current, turn.message]);
        if (turn.targetSide === "player") {
          setPlayerHP((current) => clampHP(current - turn.damage));
        } else {
          setEnemyHP((current) => clampHP(current - turn.damage));
        }
        setTurnIndex((current) => current + 1);

        window.setTimeout(() => setActiveTurn(null), 620);
      },
      turnIndex === 0 ? 650 : 1100,
    );

    return () => window.clearTimeout(timer);
  }, [replay.turns, turnIndex]);

  const replayAgain = () => {
    setTurnIndex(0);
    setActiveTurn(null);
    setLogs([]);
    setPlayerHP(replay.playerPet.maxHp);
    setEnemyHP(replay.opponent.pet.maxHp);
  };

  return (
    <main className={styles.screen}>
      <div className={styles.shell}>
        <header className={styles.topBar}>
          <div>
            <p className={styles.eyebrow}>AUTO BATTLE REPLAY</p>
            <h1 className={styles.title}>Battle Arena</h1>
          </div>
          <div className={styles.navButtons}>
            <button className={styles.button} type="button" onClick={replayAgain}>
              REPLAY
            </button>
            <button className={styles.button} type="button" onClick={() => onNavigate(PATHS.PETS)}>
              PETS
            </button>
          </div>
        </header>

        <section className={styles.hud} aria-label="Battle status">
          <HPPanel hp={playerHP} maxHp={replay.playerPet.maxHp} pet={replay.playerPet} />
          <div className={styles.versus}>VS</div>
          <HPPanel enemy hp={enemyHP} maxHp={replay.opponent.pet.maxHp} pet={replay.opponent.pet} />
        </section>

        <section className={styles.arena} aria-label="Battle arena">
          <AnimatePresence>
            {activeTurn?.isCritical && (
              <motion.div
                animate={{ opacity: [0, 1, 1, 0], scale: [0.7, 1.4, 1.05, 1] }}
                className={styles.effectBanner}
                exit={{ opacity: 0 }}
                initial={{ opacity: 0, scale: 0.65 }}
                key={`critical-${activeTurn.turn}`}
                transition={{ duration: 0.7 }}
              >
                CRITICAL!
              </motion.div>
            )}
            {activeTurn?.combo && !activeTurn.isCritical && (
              <motion.div
                animate={{ opacity: [0, 1, 1, 0], scale: [0.7, 1.25, 1, 1] }}
                className={styles.effectBanner}
                exit={{ opacity: 0 }}
                initial={{ opacity: 0, scale: 0.65 }}
                key={`combo-${activeTurn.turn}`}
                transition={{ duration: 0.7 }}
              >
                x{activeTurn.combo} COMBO
              </motion.div>
            )}
          </AnimatePresence>

          <Fighter activeTurn={activeTurn} pet={replay.playerPet} side="player" />
          <Fighter
            activeTurn={activeTurn}
            opponent={replay.opponent}
            pet={replay.opponent.pet}
            side="enemy"
          />
        </section>

        <section className={styles.terminal} aria-label="Battle log">
          <p className={styles.terminalTitle}>
            {isComplete ? `RESULT: ${replay.result.result.toUpperCase()}` : "BATTLE LOG"}
          </p>
          <div className={styles.logList}>
            {logs.length === 0 && <p className={styles.logLine}>Battle sequence booting...</p>}
            {logs.map((log, index) => (
              <motion.p
                animate={{ opacity: 1, x: 0 }}
                className={styles.logLine}
                initial={{ opacity: 0, x: -8 }}
                key={`${index}-${log}`}
              >
                {String(index + 1).padStart(2, "0")} &gt; {log}
              </motion.p>
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}

function HPPanel({
  enemy = false,
  hp,
  maxHp,
  pet,
}: {
  enemy?: boolean;
  hp: number;
  maxHp: number;
  pet: PetSummary;
}) {
  const pct = Math.max(0, Math.min(100, (hp / maxHp) * 100));

  return (
    <div className={`${styles.hpPanel} ${enemy ? styles.hpPanelEnemy : ""}`}>
      <div className={styles.hpName}>
        <span>{displayPetName(pet)}</span>
        <span>
          HP {hp}/{maxHp}
        </span>
      </div>
      <div className={styles.hpShell}>
        <motion.div
          animate={{ width: `${pct}%` }}
          className={styles.hpFill}
          transition={{ duration: 0.42, ease: "easeOut" }}
        />
      </div>
    </div>
  );
}

function Fighter({
  activeTurn,
  pet,
  side,
}: {
  activeTurn: BattleReplayTurn | null;
  opponent?: BattleOpponent;
  pet: PetSummary;
  side: "player" | "enemy";
}) {
  const asset = characterAssets[pet.attribute] ?? { kind: "placeholder", label: pet.attribute };
  const isAttacking = activeTurn?.actorSide === side;
  const isHit = activeTurn?.targetSide === side;
  const lungeDistance = side === "player" ? 62 : -62;

  return (
    <motion.div
      animate={
        isAttacking
          ? { x: [0, lungeDistance, 0], y: [0, -4, 0] }
          : isHit
            ? { x: [0, -10, 10, -7, 7, 0] }
            : { x: 0, y: 0 }
      }
      className={`${styles.fighter} ${side === "enemy" ? styles.fighterEnemy : ""}`}
      transition={{ duration: isAttacking ? 0.42 : 0.34, ease: "easeOut" }}
    >
      <AnimatePresence>
        {isHit && activeTurn && (
          <motion.div
            animate={{ opacity: 0, y: -62, scale: 1.18 }}
            className={styles.damage}
            exit={{ opacity: 0 }}
            initial={{ opacity: 1, y: -12, scale: 0.82 }}
            key={`${side}-damage-${activeTurn.turn}`}
            transition={{ duration: 0.72, ease: "easeOut" }}
          >
            -{activeTurn.damage}
          </motion.div>
        )}
      </AnimatePresence>

      {asset.kind === "sprite" ? (
        <GopherSprite style={{ transform: "scale(1.35)", transformOrigin: "bottom center" }} />
      ) : (
        <div className={styles.placeholder}>{asset.label}</div>
      )}
      <div className={styles.fighterName}>{displayPetName(pet)}</div>
    </motion.div>
  );
}
