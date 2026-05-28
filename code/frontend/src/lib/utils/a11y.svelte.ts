import type { Piece } from "$lib/types/game";

export const lakePositions = [
    { x: 2, y: 4 },
    { x: 3, y: 4 },
    { x: 2, y: 5 },
    { x: 3, y: 5 },
    { x: 6, y: 4 },
    { x: 7, y: 4 },
    { x: 6, y: 5 },
    { x: 7, y: 5 },
];

export function isLake(x: number, y: number, rows = 10, isLakeCell?: (x: number, y: number) => boolean) {
    return isLakeCell?.(x, y) ?? (rows === 10 && lakePositions.some((pos) => pos.x === x && pos.y === y));
}

export function getFriendlyPieceName(rank: string | undefined): string {
    if (!rank) return "Unknown Piece";
    if (rank === "0" || rank === "F") return "Flag";
    if (rank === "B") return "Bomb";
    const mapped: Record<string, string> = {
        "1": "Spy",
        "2": "Scout",
        "3": "Miner",
        "4": "Sergeant",
        "5": "Lieutenant",
        "6": "Captain",
        "7": "Major",
        "8": "Colonel",
        "9": "General",
        "10": "Marshal",
        "M": "Marshal"
    };
    return mapped[rank] || `Rank ${rank} Piece`;
}

export function getCellAriaLabel(params: {
    x: number;
    y: number;
    piece: Piece | null;
    isSelected: boolean;
    isValidMove: boolean;
    isLake: boolean;
    viewerId: number;
    isSetupPhase: boolean;
}) {
    const colLetter = String.fromCharCode(65 + params.x);
    const rowNumber = params.y + 1;
    let label = `Column ${colLetter}, Row ${rowNumber}. `;

    if (params.isLake) {
        label += "Lake cell. Impassable.";
        return label;
    }

    if (params.piece) {
        const isOwn = params.piece.ownerId === params.viewerId || (params.isSetupPhase && params.piece.ownerId === 0);
        const ownerStr = isOwn ? "Your" : (params.piece.ownerId === 0 ? "Blue" : "Red");

        const canSee = params.piece.ownerId === params.viewerId || params.piece.revealed || params.isSetupPhase;
        if (canSee) {
            const rank = params.piece.rank || params.piece.type;
            const name = getFriendlyPieceName(rank);
            label += `${ownerStr} ${name} (Rank ${rank}).`;
        } else {
            label += `Hidden ${ownerStr} Piece.`;
        }
    } else {
        label += "Empty cell.";
    }

    if (params.isSelected) {
        label += " Selected.";
    }

    if (params.isValidMove) {
        label += " Valid move target.";
    }

    return label;
}

class A11yStore {
    announcement = $state("");
    focusedCell = $state<{ x: number; y: number } | null>(null);

    announce(message: string) {
        this.announcement = "";
        // browser's screen reader needs a brief delay to recognize sequential duplicate announcements
        setTimeout(() => {
            this.announcement = message;
        }, 50);
    }
}

export const a11yStore = new A11yStore();

export function handleBoardKeyDown(
    e: KeyboardEvent,
    params: {
        x: number;
        y: number;
        cols: number;
        rows: number;
    }
) {
    let newX = params.x;
    let newY = params.y;
    let handled = false;

    switch (e.key) {
        case "ArrowLeft":
            newX = Math.max(0, params.x - 1);
            handled = true;
            break;
        case "ArrowRight":
            newX = Math.min(params.cols - 1, params.x + 1);
            handled = true;
            break;
        case "ArrowUp":
            newY = Math.max(0, params.y - 1);
            handled = true;
            break;
        case "ArrowDown":
            newY = Math.min(params.rows - 1, params.y + 1);
            handled = true;
            break;
        case "Home":
            if (e.ctrlKey) {
                newX = 0;
                newY = 0;
            } else {
                newX = 0;
            }
            handled = true;
            break;
        case "End":
            if (e.ctrlKey) {
                newX = params.cols - 1;
                newY = params.rows - 1;
            } else {
                newX = params.cols - 1;
            }
            handled = true;
            break;
    }

    if (handled) {
        e.preventDefault();
        const element = document.getElementById(`cell-${newX}-${newY}`);
        if (element) {
            element.focus();
        }
    }
}

export function announceMove(move: any, viewerId: number): string {
    if (!move) return "";
    const colFrom = String.fromCharCode(65 + move.fromX);
    const rowFrom = move.fromY + 1;
    const colTo = String.fromCharCode(65 + move.toX);
    const rowTo = move.toY + 1;

    const playerStr = move.playerId === viewerId ? "You" : (move.playerId === 0 ? "Blue" : "Red");
    let msg = `${playerStr} moved from ${colFrom}${rowFrom} to ${colTo}${rowTo}.`;

    if (move.result === "combat" || move.attacker) {
        const attackerName = getFriendlyPieceName(move.attacker?.rank || move.attacker?.type);
        const defenderName = getFriendlyPieceName(move.defender?.rank || move.defender?.type);
        const attackerOwner = move.attacker?.ownerId === viewerId ? "Your" : (move.attacker?.ownerId === 0 ? "Blue" : "Red");
        const defenderOwner = move.defender?.ownerId === viewerId ? "Your" : (move.defender?.ownerId === 0 ? "Blue" : "Red");

        msg += ` Combat! ${attackerOwner} ${attackerName} attacked ${defenderOwner} ${defenderName}.`;
        
        if (move.attackerWon && move.defenderWon) {
            msg += " Both pieces were defeated!";
        } else if (move.attackerWon) {
            msg += ` ${attackerOwner} ${attackerName} won and defeated the ${defenderName}!`;
        } else if (move.defenderWon) {
            msg += ` ${defenderOwner} ${defenderName} won and defeated the ${attackerName}!`;
        }
    }

    return msg;
}
