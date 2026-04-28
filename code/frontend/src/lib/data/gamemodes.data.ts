import type { GameMode } from "$lib/types/game";

const ai_vs_ai: GameMode = {
    mode: "ai_vs_ai",
    icon: "🤖 vs 🤖",
    title: "AI vs AI",
    desc: "Watch two AIs battle it out.",
}

const human_vs_ai: GameMode = {
    mode: "human_vs_ai",
    icon: "🧑 vs 🤖",
    title: "Human vs AI",
    desc: "Play against the computer.",
}

const human_vs_human: GameMode = {
    mode: "human_vs_human",
    icon: "🧑 vs 🧑",
    title: "Human vs Human",
    desc: "Coming soon.",
    disabled: true,
}

const unknown: GameMode = {
    mode: "unknown",
    icon: "",
    title: "Unknown",
    desc: "Unknown.",
    disabled: true,
}

export const gamemodes = {
    ai_vs_ai,
    human_vs_ai,
    human_vs_human,
    unknown,
    fromString: (modeStr: string): GameMode => {
        return (gamemodes as any)[modeStr] || gamemodes.unknown;
    }
} as const;

export const allGamemodes = [gamemodes.ai_vs_ai, gamemodes.human_vs_ai, gamemodes.human_vs_human];


