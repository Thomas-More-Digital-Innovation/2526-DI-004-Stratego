
    const schema = {
  "asyncapi": "3.0.0",
  "info": {
    "title": "GoStrategy WebSocket API",
    "version": "1.0.0",
    "description": "Real-time communication protocol for the GoStrategy game engine. \nAll messages follow a common wrapper format:\n```json\n{\n  \"type\": \"messageType\",\n  \"data\": { ... }\n}\n```\n",
    "contact": {
      "name": "GoStrategy Team",
      "url": "https://gostrategy.dotsem.be"
    }
  },
  "servers": {
    "local": {
      "host": "localhost:8080",
      "protocol": "ws",
      "description": "Local development server"
    },
    "production": {
      "host": "api.gostrategy.dotsem.be",
      "protocol": "wss",
      "description": "Production server"
    }
  },
  "channels": {
    "game": {
      "address": "/ws/game/{gameID}",
      "parameters": {
        "gameID": {
          "description": "The unique ID of the game session."
        }
      },
      "messages": {
        "move": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "move",
                "x-parser-schema-id": "<anonymous-schema-3>"
              },
              "data": {
                "type": "object",
                "required": [
                  "from",
                  "to"
                ],
                "properties": {
                  "from": {
                    "type": "object",
                    "required": [
                      "x",
                      "y"
                    ],
                    "properties": {
                      "x": {
                        "type": "integer",
                        "minimum": 0,
                        "maximum": 9,
                        "x-parser-schema-id": "<anonymous-schema-5>"
                      },
                      "y": {
                        "type": "integer",
                        "minimum": 0,
                        "maximum": 9,
                        "x-parser-schema-id": "<anonymous-schema-6>"
                      }
                    },
                    "x-parser-schema-id": "Position"
                  },
                  "to": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from"
                },
                "x-parser-schema-id": "<anonymous-schema-4>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-2>"
          },
          "x-parser-unique-object-id": "move",
          "x-parser-message-name": "MoveRequest"
        },
        "getValidMoves": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "getValidMoves",
                "x-parser-schema-id": "<anonymous-schema-8>"
              },
              "data": {
                "type": "object",
                "required": [
                  "position"
                ],
                "properties": {
                  "position": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from"
                },
                "x-parser-schema-id": "<anonymous-schema-9>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-7>"
          },
          "x-parser-unique-object-id": "getValidMoves",
          "x-parser-message-name": "GetValidMovesRequest"
        },
        "swapPieces": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "swapPieces",
                "x-parser-schema-id": "<anonymous-schema-11>"
              },
              "data": {
                "type": "object",
                "required": [
                  "pos1",
                  "pos2"
                ],
                "properties": {
                  "pos1": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from",
                  "pos2": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from"
                },
                "x-parser-schema-id": "<anonymous-schema-12>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-10>"
          },
          "x-parser-unique-object-id": "swapPieces",
          "x-parser-message-name": "SwapPiecesRequest"
        },
        "randomizeSetup": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "randomizeSetup",
                "x-parser-schema-id": "<anonymous-schema-14>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "playerId": {
                    "type": "integer",
                    "description": "Optional player ID for AI vs AI games",
                    "x-parser-schema-id": "<anonymous-schema-16>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-15>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-13>"
          },
          "x-parser-unique-object-id": "randomizeSetup",
          "x-parser-message-name": "RandomizeSetupRequest"
        },
        "loadSetup": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "loadSetup",
                "x-parser-schema-id": "<anonymous-schema-18>"
              },
              "data": {
                "type": "object",
                "required": [
                  "setupData"
                ],
                "properties": {
                  "playerId": {
                    "type": "integer",
                    "x-parser-schema-id": "<anonymous-schema-20>"
                  },
                  "setupData": {
                    "type": "string",
                    "description": "Base64 encoded 40 bytes of setup data",
                    "x-parser-schema-id": "<anonymous-schema-21>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-19>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-17>"
          },
          "x-parser-unique-object-id": "loadSetup",
          "x-parser-message-name": "LoadSetupRequest"
        },
        "startGame": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "startGame",
                "x-parser-schema-id": "<anonymous-schema-23>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "headless": {
                    "type": "boolean",
                    "default": false,
                    "x-parser-schema-id": "<anonymous-schema-25>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-24>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-22>"
          },
          "x-parser-unique-object-id": "startGame",
          "x-parser-message-name": "StartGameRequest"
        },
        "setSpeed": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "setSpeed",
                "x-parser-schema-id": "<anonymous-schema-27>"
              },
              "data": {
                "type": "object",
                "required": [
                  "speedMs"
                ],
                "properties": {
                  "speedMs": {
                    "type": "integer",
                    "minimum": 500,
                    "maximum": 5000,
                    "x-parser-schema-id": "<anonymous-schema-29>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-28>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-26>"
          },
          "x-parser-unique-object-id": "setSpeed",
          "x-parser-message-name": "SetSpeedRequest"
        },
        "animationComplete": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "type": "string",
                "enum": [
                  "animationComplete",
                  "pause",
                  "unpause",
                  "ping",
                  "step"
                ],
                "x-parser-schema-id": "<anonymous-schema-31>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-30>"
          },
          "x-parser-unique-object-id": "animationComplete",
          "x-parser-message-name": "SimpleCommand"
        },
        "pause": "$ref:$.channels.game.messages.animationComplete",
        "unpause": "$ref:$.channels.game.messages.animationComplete",
        "ping": "$ref:$.channels.game.messages.animationComplete",
        "step": "$ref:$.channels.game.messages.animationComplete",
        "gameState": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "gameState",
                "x-parser-schema-id": "<anonymous-schema-33>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "round": {
                    "type": "integer",
                    "x-parser-schema-id": "<anonymous-schema-34>"
                  },
                  "currentPlayerId": {
                    "type": "integer",
                    "x-parser-schema-id": "<anonymous-schema-35>"
                  },
                  "currentPlayerName": {
                    "type": "string",
                    "x-parser-schema-id": "<anonymous-schema-36>"
                  },
                  "isGameOver": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-37>"
                  },
                  "winnerId": {
                    "type": "integer",
                    "nullable": true,
                    "x-parser-schema-id": "<anonymous-schema-38>"
                  },
                  "player1Score": {
                    "type": "integer",
                    "x-parser-schema-id": "<anonymous-schema-39>"
                  },
                  "player2Score": {
                    "type": "integer",
                    "x-parser-schema-id": "<anonymous-schema-40>"
                  },
                  "waitingForInput": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-41>"
                  },
                  "paused": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-42>"
                  },
                  "isSetupPhase": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-43>"
                  },
                  "setupRemainingSecs": {
                    "type": "integer",
                    "x-parser-schema-id": "<anonymous-schema-44>"
                  }
                },
                "x-parser-schema-id": "GameState"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-32>"
          },
          "x-parser-unique-object-id": "gameState",
          "x-parser-message-name": "GameStateUpdate"
        },
        "boardState": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "boardState",
                "x-parser-schema-id": "<anonymous-schema-46>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "board": {
                    "type": "array",
                    "items": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "type": {
                            "type": "string",
                            "example": "Marshal",
                            "x-parser-schema-id": "<anonymous-schema-49>"
                          },
                          "rank": {
                            "type": "string",
                            "example": "10",
                            "x-parser-schema-id": "<anonymous-schema-50>"
                          },
                          "ownerId": {
                            "type": "integer",
                            "x-parser-schema-id": "<anonymous-schema-51>"
                          },
                          "ownerName": {
                            "type": "string",
                            "x-parser-schema-id": "<anonymous-schema-52>"
                          },
                          "revealed": {
                            "type": "boolean",
                            "x-parser-schema-id": "<anonymous-schema-53>"
                          },
                          "icon": {
                            "type": "string",
                            "x-parser-schema-id": "<anonymous-schema-54>"
                          },
                          "position": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from"
                        },
                        "x-parser-schema-id": "Piece"
                      },
                      "x-parser-schema-id": "<anonymous-schema-48>"
                    },
                    "x-parser-schema-id": "<anonymous-schema-47>"
                  },
                  "width": {
                    "type": "integer",
                    "const": 10,
                    "x-parser-schema-id": "<anonymous-schema-55>"
                  },
                  "height": {
                    "type": "integer",
                    "const": 10,
                    "x-parser-schema-id": "<anonymous-schema-56>"
                  }
                },
                "x-parser-schema-id": "BoardState"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-45>"
          },
          "x-parser-unique-object-id": "boardState",
          "x-parser-message-name": "BoardStateUpdate"
        },
        "combat": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "combat",
                "x-parser-schema-id": "<anonymous-schema-58>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "attacker": "$ref:$.channels.game.messages.boardState.payload.properties.data.properties.board.items.items",
                  "defender": "$ref:$.channels.game.messages.boardState.payload.properties.data.properties.board.items.items",
                  "attackerWon": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-59>"
                  },
                  "defenderWon": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-60>"
                  },
                  "attackerDied": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-61>"
                  },
                  "defenderDied": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-62>"
                  }
                },
                "x-parser-schema-id": "CombatResult"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-57>"
          },
          "x-parser-unique-object-id": "combat",
          "x-parser-message-name": "CombatUpdate"
        },
        "moveResult": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "moveResult",
                "x-parser-schema-id": "<anonymous-schema-64>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "success": {
                    "type": "boolean",
                    "x-parser-schema-id": "<anonymous-schema-66>"
                  },
                  "error": {
                    "type": "string",
                    "x-parser-schema-id": "<anonymous-schema-67>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-65>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-63>"
          },
          "x-parser-unique-object-id": "moveResult",
          "x-parser-message-name": "MoveResult"
        },
        "validMoves": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "validMoves",
                "x-parser-schema-id": "<anonymous-schema-69>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "position": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from",
                  "validMoves": {
                    "type": "array",
                    "items": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from",
                    "x-parser-schema-id": "<anonymous-schema-71>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-70>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-68>"
          },
          "x-parser-unique-object-id": "validMoves",
          "x-parser-message-name": "ValidMovesResponse"
        },
        "gameOver": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "gameOver",
                "x-parser-schema-id": "<anonymous-schema-73>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "winnerId": {
                    "type": "integer",
                    "nullable": true,
                    "x-parser-schema-id": "<anonymous-schema-75>"
                  },
                  "winnerName": {
                    "type": "string",
                    "x-parser-schema-id": "<anonymous-schema-76>"
                  },
                  "winCause": {
                    "type": "string",
                    "x-parser-schema-id": "<anonymous-schema-77>"
                  },
                  "round": {
                    "type": "integer",
                    "x-parser-schema-id": "<anonymous-schema-78>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-74>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-72>"
          },
          "x-parser-unique-object-id": "gameOver",
          "x-parser-message-name": "GameOverEvent"
        },
        "error": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "error",
                "x-parser-schema-id": "<anonymous-schema-80>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "error": {
                    "type": "string",
                    "x-parser-schema-id": "<anonymous-schema-82>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-81>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-79>"
          },
          "x-parser-unique-object-id": "error",
          "x-parser-message-name": "ErrorEvent"
        },
        "moveHistory": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "const": "moveHistory",
                "x-parser-schema-id": "<anonymous-schema-84>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "moves": {
                    "type": "array",
                    "items": {
                      "type": "object",
                      "properties": {
                        "from": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from",
                        "to": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from"
                      },
                      "x-parser-schema-id": "Move"
                    },
                    "x-parser-schema-id": "<anonymous-schema-86>"
                  }
                },
                "x-parser-schema-id": "<anonymous-schema-85>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-83>"
          },
          "x-parser-unique-object-id": "moveHistory",
          "x-parser-message-name": "MoveHistoryUpdate"
        },
        "pong": {
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "type": "string",
                "enum": [
                  "pong"
                ],
                "x-parser-schema-id": "<anonymous-schema-88>"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-87>"
          },
          "x-parser-unique-object-id": "pong",
          "x-parser-message-name": "SimpleEvent"
        }
      },
      "x-parser-unique-object-id": "game"
    }
  },
  "operations": {
    "sendCommands": {
      "action": "send",
      "channel": "$ref:$.channels.game",
      "messages": [
        "$ref:$.channels.game.messages.move",
        "$ref:$.channels.game.messages.getValidMoves",
        "$ref:$.channels.game.messages.swapPieces",
        "$ref:$.channels.game.messages.randomizeSetup",
        "$ref:$.channels.game.messages.loadSetup",
        "$ref:$.channels.game.messages.startGame",
        "$ref:$.channels.game.messages.setSpeed",
        "$ref:$.channels.game.messages.animationComplete",
        "$ref:$.channels.game.messages.animationComplete",
        "$ref:$.channels.game.messages.animationComplete",
        "$ref:$.channels.game.messages.animationComplete",
        "$ref:$.channels.game.messages.animationComplete"
      ],
      "x-parser-unique-object-id": "sendCommands"
    },
    "receiveEvents": {
      "action": "receive",
      "channel": "$ref:$.channels.game",
      "messages": [
        "$ref:$.channels.game.messages.gameState",
        "$ref:$.channels.game.messages.boardState",
        "$ref:$.channels.game.messages.combat",
        "$ref:$.channels.game.messages.moveResult",
        "$ref:$.channels.game.messages.validMoves",
        "$ref:$.channels.game.messages.gameOver",
        "$ref:$.channels.game.messages.error",
        "$ref:$.channels.game.messages.moveHistory",
        "$ref:$.channels.game.messages.pong"
      ],
      "x-parser-unique-object-id": "receiveEvents"
    }
  },
  "components": {
    "messages": {
      "MoveRequest": "$ref:$.channels.game.messages.move",
      "GetValidMovesRequest": "$ref:$.channels.game.messages.getValidMoves",
      "SwapPiecesRequest": "$ref:$.channels.game.messages.swapPieces",
      "RandomizeSetupRequest": "$ref:$.channels.game.messages.randomizeSetup",
      "LoadSetupRequest": "$ref:$.channels.game.messages.loadSetup",
      "StartGameRequest": "$ref:$.channels.game.messages.startGame",
      "SetSpeedRequest": "$ref:$.channels.game.messages.setSpeed",
      "SimpleCommand": "$ref:$.channels.game.messages.animationComplete",
      "GameStateUpdate": "$ref:$.channels.game.messages.gameState",
      "BoardStateUpdate": "$ref:$.channels.game.messages.boardState",
      "CombatUpdate": "$ref:$.channels.game.messages.combat",
      "MoveResult": "$ref:$.channels.game.messages.moveResult",
      "ValidMovesResponse": "$ref:$.channels.game.messages.validMoves",
      "GameOverEvent": "$ref:$.channels.game.messages.gameOver",
      "ErrorEvent": "$ref:$.channels.game.messages.error",
      "MoveHistoryUpdate": "$ref:$.channels.game.messages.moveHistory",
      "SimpleEvent": "$ref:$.channels.game.messages.pong"
    },
    "schemas": {
      "Position": "$ref:$.channels.game.messages.move.payload.properties.data.properties.from",
      "Move": "$ref:$.channels.game.messages.moveHistory.payload.properties.data.properties.moves.items",
      "Piece": "$ref:$.channels.game.messages.boardState.payload.properties.data.properties.board.items.items",
      "GameState": "$ref:$.channels.game.messages.gameState.payload.properties.data",
      "BoardState": "$ref:$.channels.game.messages.boardState.payload.properties.data",
      "CombatResult": "$ref:$.channels.game.messages.combat.payload.properties.data"
    }
  },
  "x-parser-spec-parsed": true,
  "x-parser-api-version": 3,
  "x-parser-spec-stringified": true
};
    const config = {"show":{"sidebar":true},"sidebar":{"showOperations":"byDefault"}};
    const appRoot = document.getElementById('root');
    AsyncApiStandalone.render(
        { schema, config, }, appRoot
    );
  