# Stratego Frontend

A beautiful, modern frontend for the Stratego board game built with **Svelte 5** and **SvelteKit**.

## 🎮 Features

- **Svelte 5 Runes**: Using the latest Svelte 5 features for reactive state management
- **Beautiful UI**: Modern design with smooth animations and gradients
- **Interactive Gameplay**: Click to select pieces and move them around the board
- **Real-time Combat**: See combat results instantly with visual feedback
- **Game State Management**: Complete game state tracking with move history
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **Piece Information**: Hover and select pieces to see detailed information
- **Visual Indicators**: Clear highlighting for selected pieces and valid moves

## 🚀 Getting Started

### Prerequisites

- Node.js (v18 or higher recommended)
- npm, yarn, pnpm, or bun

### Installation

```sh
# Install dependencies
npm install
```

## Developing

Start a development server:

```sh
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

The application will be available at `http://localhost:5173`

## Building

To create a production version:

```sh
npm run build
```

You can preview the production build with `npm run preview`.

## 🎯 How to Play

1. **Start Game**: Click "Start Game" on the welcome screen
2. **Select Piece**: Click on your pieces (shown with your color) to select them
3. **Move**: Click on highlighted valid move positions to move your piece
4. **Combat**: Attack enemy pieces by moving to their position
5. **Win**: Capture the opponent's flag or eliminate all movable pieces

## 🎨 Tech Stack

- **Framework**: SvelteKit with Svelte 5
- **Language**: TypeScript
- **Styling**: Tailwind CSS + Custom CSS
- **Build Tool**: Vite

## 📁 Project Structure

```
src/
├── lib/
│   ├── types.ts              # TypeScript type definitions
│   ├── gameState.ts          # Game logic and state management
│   └── components/
│       ├── BoardCell.svelte  # Individual board cell component
│       ├── GameBoard.svelte  # Main game board component
│       └── GameInfo.svelte   # Game information sidebar
├── routes/
│   └── +page.svelte          # Main game page
└── app.css                   # Global styles
```

## 🎮 Game Rules

### Objective
Capture your opponent's Flag to win the game!

### Piece Types
- **Flag (🚩)**: Must be captured to win
- **Bomb (💣)**: Cannot move, destroys most attackers
- **Spy (🕵️)**: Can defeat Marshal if attacking
- **Scout (🔭)**: Can move multiple spaces
- **Miner (⛏️)**: Can defuse bombs
- **Sergeant - Marshal**: Various combat pieces with different ranks

### Movement
- Most pieces move one square orthogonally (up, down, left, right)
- Scouts can move any number of spaces in a straight line
- Cannot move through lakes or other pieces

### Combat
- Higher rank usually wins
- Special rules:
  - Spy defeats Marshal (if attacking)
  - Miner defeats Bomb
  - Bomb defeats all except Miner
  - Equal ranks result in both pieces being removed

## 🔧 Development

### Adding New Features

The game state is managed using Svelte 5 runes in `gameState.ts`. To add new features:

1. Update types in `lib/types.ts`
2. Modify game logic in `lib/gameState.ts`
3. Update UI components as needed

### Styling

The project uses a combination of:
- Tailwind CSS for utility classes
- Component-scoped CSS for specific styling
- CSS custom properties for theming

## 📝 License

Part of the Digital Innovation Stratego Project

## 🤝 Contributing

This is an educational project for the Digital Innovation course.
