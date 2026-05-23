import { assign, setup } from "xstate";
import type {
  BattleOpponent,
  BattleResult,
  MyPetsResponse,
  PetTrainingStat,
  TrainingResult,
} from "./api";

export interface PetPageContext {
  data: MyPetsResponse | null;
  selectedPetId: string | null;
  opponents: BattleOpponent[];
  selectedOpponentId: string | null;
  battleResult: BattleResult | null;
  statusMessage: string | null;
  noticeMessage: string | null;
  trainingStat: PetTrainingStat | null;
}

export type PetPageEvent =
  | { type: "NOTICE"; message: string }
  | { type: "LOAD_SUCCESS"; data: MyPetsResponse; statusMessage?: string | null }
  | { type: "LOAD_FAILURE"; message: string }
  | { type: "SELECT_PET"; petId: string }
  | { type: "OPPONENTS_SUCCESS"; opponents: BattleOpponent[] }
  | { type: "OPPONENTS_FAILURE" }
  | { type: "SELECT_OPPONENT"; userId: string }
  | { type: "TRAIN"; stat: PetTrainingStat }
  | { type: "INSUFFICIENT_CP" }
  | { type: "TRAIN_SUCCESS"; result: TrainingResult; message: string }
  | { type: "TRAIN_FAILURE"; message: string }
  | { type: "BATTLE" }
  | { type: "BATTLE_SUCCESS"; result: BattleResult }
  | { type: "BATTLE_FAILURE"; message: string };

export const petPageMachine = setup({
  types: {
    context: {} as PetPageContext,
    events: {} as PetPageEvent,
  },
  actions: {
    setNotice: assign({
      noticeMessage: ({ event }) => (event.type === "NOTICE" ? event.message : null),
    }),
    setLoadedData: assign({
      data: ({ event }) => (event.type === "LOAD_SUCCESS" ? event.data : null),
      selectedPetId: ({ event }) => {
        if (event.type !== "LOAD_SUCCESS") return null;
        return event.data.currentGuildPet?.id ?? event.data.pets[0]?.id ?? null;
      },
      statusMessage: ({ event }) =>
        event.type === "LOAD_SUCCESS" ? (event.statusMessage ?? null) : null,
    }),
    setLoadFailure: assign({
      data: null,
      selectedPetId: null,
      statusMessage: ({ event }) => (event.type === "LOAD_FAILURE" ? event.message : null),
    }),
    selectPet: assign({
      selectedPetId: ({ event }) => (event.type === "SELECT_PET" ? event.petId : null),
      battleResult: null,
      statusMessage: null,
    }),
    setOpponents: assign({
      opponents: ({ event }) => (event.type === "OPPONENTS_SUCCESS" ? event.opponents : []),
      selectedOpponentId: ({ event }) =>
        event.type === "OPPONENTS_SUCCESS" ? (event.opponents[0]?.userId ?? null) : null,
    }),
    selectOpponent: assign({
      selectedOpponentId: ({ event }) => (event.type === "SELECT_OPPONENT" ? event.userId : null),
    }),
    setTraining: assign({
      trainingStat: ({ event }) => (event.type === "TRAIN" ? event.stat : null),
      statusMessage: null,
    }),
    setInsufficientCP: assign({
      statusMessage: "CPが足りません",
    }),
    setTrainingSuccess: assign({
      data: ({ context, event }) => {
        if (event.type !== "TRAIN_SUCCESS" || !context.data) return context.data;

        return {
          ...context.data,
          cpBalance: event.result.cpAfter,
          currentGuildPet: event.result.pet,
          pets: context.data.pets.map((pet) =>
            pet.id === event.result.pet.id ? event.result.pet : pet,
          ),
        };
      },
      statusMessage: ({ event }) => (event.type === "TRAIN_SUCCESS" ? event.message : null),
      trainingStat: null,
    }),
    setTrainingFailure: assign({
      statusMessage: ({ event }) => (event.type === "TRAIN_FAILURE" ? event.message : null),
      trainingStat: null,
    }),
    startBattle: assign({
      statusMessage: null,
      battleResult: null,
    }),
    setBattleSuccess: assign({
      battleResult: ({ event }) => (event.type === "BATTLE_SUCCESS" ? event.result : null),
    }),
    setBattleFailure: assign({
      statusMessage: ({ event }) => (event.type === "BATTLE_FAILURE" ? event.message : null),
    }),
  },
}).createMachine({
  id: "petPage",
  initial: "loading",
  context: {
    data: null,
    selectedPetId: null,
    opponents: [],
    selectedOpponentId: null,
    battleResult: null,
    statusMessage: null,
    noticeMessage: null,
    trainingStat: null,
  },
  on: {
    NOTICE: {
      actions: "setNotice",
    },
    OPPONENTS_SUCCESS: {
      actions: "setOpponents",
    },
    OPPONENTS_FAILURE: {},
    SELECT_PET: {
      actions: "selectPet",
    },
    SELECT_OPPONENT: {
      actions: "selectOpponent",
    },
  },
  states: {
    loading: {
      on: {
        LOAD_SUCCESS: {
          target: "ready",
          actions: "setLoadedData",
        },
        LOAD_FAILURE: {
          target: "failed",
          actions: "setLoadFailure",
        },
      },
    },
    failed: {
      on: {
        LOAD_SUCCESS: {
          target: "ready",
          actions: "setLoadedData",
        },
      },
    },
    ready: {
      on: {
        TRAIN: {
          target: "training",
          actions: "setTraining",
        },
        INSUFFICIENT_CP: {
          actions: "setInsufficientCP",
        },
        BATTLE: {
          target: "battling",
          actions: "startBattle",
        },
      },
    },
    training: {
      on: {
        TRAIN_SUCCESS: {
          target: "ready",
          actions: "setTrainingSuccess",
        },
        TRAIN_FAILURE: {
          target: "ready",
          actions: "setTrainingFailure",
        },
      },
    },
    battling: {
      on: {
        BATTLE_SUCCESS: {
          target: "ready",
          actions: "setBattleSuccess",
        },
        BATTLE_FAILURE: {
          target: "ready",
          actions: "setBattleFailure",
        },
      },
    },
  },
});
