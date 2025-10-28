#!/bin/zsh

# Stratego Test Script
# This script helps you test the game frontend

echo "🎮 Stratego Frontend Test Guide"
echo "================================\n"

echo "📋 Pre-flight Checklist:"
echo "  1. Backend is running on localhost:8080"
echo "  2. Frontend will start on localhost:5173"
echo "  3. Browser console is open (F12)\n"

echo "🧪 What to Test:\n"

echo "1. VISUAL DISPLAY"
echo "   ✓ Your pieces (bottom 4 rows) show icons on RED background"
echo "   ✓ Enemy pieces (top 4 rows) show '?' on BLUE background"
echo "   ✓ Empty cells (middle 2 rows except lakes) are DARK GRAY"
echo "   ✓ Lakes (2x2 squares in middle) show 🌊 on BLUE\n"

echo "2. PIECE SELECTION"
echo "   ✓ Click your movable piece → GOLD border appears"
echo "   ✓ Click your Bomb/Flag → Console: '❌ Piece has no valid moves'"
echo "   ✓ Click enemy piece → Console: '❌ Not your piece'"
echo "   ✓ Click empty cell → Console: '❌ Empty cell'\n"

echo "3. MOVEMENT HIGHLIGHTING"
echo "   ✓ After selecting piece → Green pulsing borders on valid moves"
echo "   ✓ Can move orthogonally (up/down/left/right)"
echo "   ✓ Scouts can move multiple spaces"
echo "   ✓ Cannot move through pieces or lakes\n"

echo "4. MAKING MOVES"
echo "   ✓ Click green highlighted cell → Piece moves"
echo "   ✓ Console: '✓ Making move from {...} to {...}'"
echo "   ✓ Board updates with new position"
echo "   ✓ AI responds with counter-move\n"

echo "5. CONSOLE LOGS TO CHECK"
echo "   Look for these messages:"
echo "   • 'Received gameState' and 'Received boardState'"
echo "   • '🎮 isHumanTurn calculation' should show result: true"
echo "   • Enemy piece sample should show 'ownerName: \"AI Red\"'"
echo "   • No '❌ Click ignored' messages when it's your turn\n"

echo "🐛 Common Issues:\n"

echo "Issue: 'Game not running' error"
echo "  → Check console for 'waitingForInput: true'"
echo "  → Check 'currentPlayerId: 0' (your turn)"
echo "  → Restart backend if game finished immediately\n"

echo "Issue: Can't see enemy pieces"
echo "  → Check console: enemy piece should have 'ownerName: \"AI Red\"'"
echo "  → Should NOT have 'type' field (it's hidden)"
echo "  → Should still show '?' on blue background\n"

echo "Issue: Empty cells are selectable"
echo "  → Check console: empty cell should have 'ownerName: \"\"'"
echo "  → Should see '❌ Empty cell' message\n"

echo "\n🚀 Starting Frontend Development Server...\n"

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo "❌ Error: Not in frontend directory!"
    echo "   Run: cd /home/sem/prog/go/2526-DI-004-Stratego/code/frontend"
    exit 1
fi

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Start the dev server
npm run dev
