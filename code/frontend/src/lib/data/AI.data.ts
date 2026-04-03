import type { AI } from "$lib/types/game";
import fafoImage from "$lib/assets/ai/fafo.webp";
import fatoImage from "$lib/assets/ai/fato.webp";

export const AIs: AI[] = [
    {
        name: "FAFO",
        id: "fafo",
        description:
            "The Fuck Around & Find Out AI is a simple random-move AI.",
        image: fafoImage,
    },
    {
        name: "FATO",
        id: "fato",
        description:
            "The Fuck Around & Try Out AI is a random-move AI that can remember the board and act on it.",
        image: fatoImage,
    },
];
