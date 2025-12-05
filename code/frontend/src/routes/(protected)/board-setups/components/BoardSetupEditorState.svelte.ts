import { PIECE_INVENTORY, type BoardSetup } from "$lib/types/boardSetup";

export function createBoardSetupState(
    initialSetup: string | undefined,
    onSave: (setup: BoardSetup) => void
) {
    let board = $state<string[]>(initializeBoard());
    let selectedPiece = $state<string | null>(null);
    let selectedPieceIndex = $state<number | null>(null);

    let setupName = $state("");
    let setupDescription = $state("");
    let isDefault = $state(false);

    function initializeBoard(): string[] {
        if (initialSetup) {
            return parseSetup(initialSetup);
        }
        const pieces: string[] = [];
        Object.entries(PIECE_INVENTORY).forEach(([char, info]) => {
            for (let i = 0; i < info.count; i++) {
                pieces.push(char);
            }
        });
        return shuffleArray(pieces);
    }

    function parseSetup(setup: string): string[] {
        if (setup.length === 40) {
            return Array.from(setup);
        }
        const chars = Array.from("0B123456789M");
        return Array(40)
            .fill(".")
            .map((_, i) => chars[i % chars.length]);
    }

    function shuffleArray(array: string[]): string[] {
        const arr = [...array];
        for (let i = arr.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [arr[i], arr[j]] = [arr[j], arr[i]];
        }
        return arr;
    }

    function getRemainingPieces(): Record<string, number> {
        const remaining: Record<string, number> = {};
        Object.entries(PIECE_INVENTORY).forEach(([char, info]) => {
            remaining[char] = info.count;
        });
        board.forEach((piece) => {
            if (piece !== "." && remaining[piece] !== undefined) {
                remaining[piece]--;
            }
        });
        return remaining;
    }

    function handleCellClick(index: number) {
        if (selectedPiece) {
            if (selectedPieceIndex !== null) {
                // Piece selected from board - swap or move
                if (board[index] === ".") {
                    // Move to empty cell
                    board[index] = selectedPiece;
                    board[selectedPieceIndex] = ".";
                } else {
                    // Swap with another piece
                    [board[selectedPieceIndex], board[index]] = [
                        board[index],
                        board[selectedPieceIndex],
                    ];
                }
            } else {
                // Piece selected from bank - only place in empty cells
                if (board[index] === ".") {
                    board[index] = selectedPiece;
                }
            }
            selectedPiece = null;
            selectedPieceIndex = null;
        } else if (board[index] !== ".") {
            // Pick up piece from board
            selectedPiece = board[index];
            selectedPieceIndex = index;
        }
    }

    function handlePieceSelect(piece: string) {
        const remaining = getRemainingPieces();
        if (remaining[piece] > 0) {
            selectedPiece = piece;
            selectedPieceIndex = null; // From bank, not from board
        }
    }

    function swapCells(fromIndex: number, toIndex: number) {
        [board[fromIndex], board[toIndex]] = [board[toIndex], board[fromIndex]];
    }

    function createBoardSetup() {
        return {
            name: setupName,
            description: setupDescription,
            isDefault: isDefault,
            setupData: board.join(""),
        };
    }

    function handleSave() {
        onSave(createBoardSetup());
    }

    return {
        get board() {
            return board;
        },
        set board(value) {
            board = value;
        },
        get selectedPiece() {
            return selectedPiece;
        },
        set selectedPiece(value) {
            selectedPiece = value;
        },
        get selectedPieceIndex() {
            return selectedPieceIndex;
        },
        set selectedPieceIndex(value) {
            selectedPieceIndex = value;
        },
        get setupName() {
            return setupName;
        },
        set setupName(value) {
            setupName = value;
        },
        get setupDescription() {
            return setupDescription;
        },
        set setupDescription(value) {
            setupDescription = value;
        },
        get isDefault() {
            return isDefault;
        },
        set isDefault(value) {
            isDefault = value;
        },
        getRemainingPieces,
        handleCellClick,
        handlePieceSelect,
        swapCells,
        handleSave,
    };
}
