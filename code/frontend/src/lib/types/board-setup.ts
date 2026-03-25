const blueIcons = import.meta.glob('$lib/assets/pieces/blue/*.png', { eager: true, import: 'default' });
const redIcons = import.meta.glob('$lib/assets/pieces/red/*.png', { eager: true, import: 'default' });

const getIcon = (icons: Record<string, any>, name: string) => {
    const key = Object.keys(icons).find(k => k.toLowerCase().endsWith(`/${name.toLowerCase()}.png`));
    return (key ? icons[key] : '') as string;
};

export interface BoardSetup {
    id: number;
    name: string;
    description: string;
    created_at: string;
    updated_at: string;
    is_default: boolean;
    setup_data: string;
}

export interface PieceInfo {
    name: string;
    icon_blue: string;
    icon_red: string;
    rank: string;
    count: number;
}

export const PIECE_INVENTORY: Record<string, PieceInfo> = {
    '0': { name: 'Flag', icon_blue: getIcon(blueIcons, 'f'), icon_red: getIcon(redIcons, 'f'), rank: 'F', count: 1 },
    'B': { name: 'Bomb', icon_blue: getIcon(blueIcons, 'b'), icon_red: getIcon(redIcons, 'b'), rank: 'B', count: 6 },
    '1': { name: 'Spy', icon_blue: getIcon(blueIcons, '1'), icon_red: getIcon(redIcons, '1'), rank: 'S', count: 1 },
    '2': { name: 'Scout', icon_blue: getIcon(blueIcons, '2'), icon_red: getIcon(redIcons, '2'), rank: '2', count: 8 },
    '3': { name: 'Miner', icon_blue: getIcon(blueIcons, '3'), icon_red: getIcon(redIcons, '3'), rank: '3', count: 5 },
    '4': { name: 'Sergeant', icon_blue: getIcon(blueIcons, '4'), icon_red: getIcon(redIcons, '4'), rank: '4', count: 4 },
    '5': { name: 'Lieutenant', icon_blue: getIcon(blueIcons, '5'), icon_red: getIcon(redIcons, '5'), rank: '5', count: 4 },
    '6': { name: 'Captain', icon_blue: getIcon(blueIcons, '6'), icon_red: getIcon(redIcons, '6'), rank: '6', count: 4 },
    '7': { name: 'Major', icon_blue: getIcon(blueIcons, '7'), icon_red: getIcon(redIcons, '7'), rank: '7', count: 3 },
    '8': { name: 'Colonel', icon_blue: getIcon(blueIcons, '8'), icon_red: getIcon(redIcons, '8'), rank: '8', count: 2 },
    '9': { name: 'General', icon_blue: getIcon(blueIcons, '9'), icon_red: getIcon(redIcons, '9'), rank: '9', count: 1 },
    'M': { name: 'Marshal', icon_blue: getIcon(blueIcons, '10'), icon_red: getIcon(redIcons, '10'), rank: '10', count: 1 },
};

export const MAX_BOARD_SETUPS = 10;
export const BOARD_ROWS = 4;
export const BOARD_COLS = 10;
export const BOARD_CELLS = BOARD_ROWS * BOARD_COLS;
