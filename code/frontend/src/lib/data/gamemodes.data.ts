import type { GameMode } from "$lib/types/game";

export const gamemodes: {
    mode: GameMode;
    icon: string;
    title: string;
    desc: string;
}[] = [
        {
            mode: "human_vs_ai",
            icon: "🧑 vs 🤖",
            title: "Human vs AI",
            desc: "Play against the computer.",
        },
        {
            mode: "ai_vs_ai",
            icon: "🤖 vs 🤖",
            title: "AI vs AI",
            desc: "Watch two AIs battle it out.",
        },
        {
            mode: "human_vs_human",
            icon: "🧑 vs 🧑",
            title: "Human vs Human",
            desc: "Coming soon.",
        },
    ];