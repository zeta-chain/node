#!/bin/bash
set -euo pipefail

ENVIRONMENT=${ENVIRONMENT:-"athens"}
CONFIG_FILE=${CONFIG_FILE:-"$(dirname "$0")/config/default.json"}
API_KEY=${API_KEY:-""}
OUTPUT_DIR=${OUTPUT_DIR:-"."}
MAX_BLOCKS_OVERRIDE=${MAX_BLOCKS_OVERRIDE:-""}
MAX_TXS_OVERRIDE=${MAX_TXS_OVERRIDE:-""}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --environment)
      ENVIRONMENT="$2"
      shift 2
      ;;
    --config)
      CONFIG_FILE="$2"
      shift 2
      ;;
    --api-key)
      API_KEY="$2"
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --max-blocks)
      MAX_BLOCKS_OVERRIDE="$2"
      shift 2
      ;;
    --max-txs)
      MAX_TXS_OVERRIDE="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [[ -z "$API_KEY" ]]; then
  echo "ERROR: API_KEY is required"
  exit 1
fi

if ! command -v jq &> /dev/null; then
  echo "ERROR: jq is not installed. Please install jq to continue."
  exit 1
fi

if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "ERROR: Config file not found: $CONFIG_FILE"
  exit 1
fi

if [[ "$ENVIRONMENT" == "athens" ]]; then
  MAX_BLOCKS=$(jq -r '.athens.max_blocks' "$CONFIG_FILE")
  MAX_TXS=$(jq -r '.athens.max_transactions' "$CONFIG_FILE")
elif [[ "$ENVIRONMENT" == "mainnet" ]]; then
  MAX_BLOCKS=$(jq -r '.mainnet.max_blocks' "$CONFIG_FILE")
  MAX_TXS=$(jq -r '.mainnet.max_transactions' "$CONFIG_FILE")
else
  echo "ERROR: Unknown environment: $ENVIRONMENT. Must be 'athens' or 'mainnet'."
  exit 1
fi

if [[ -n "$MAX_BLOCKS_OVERRIDE" ]]; then
  MAX_BLOCKS="$MAX_BLOCKS_OVERRIDE"
  echo "Overriding max_blocks from command line: $MAX_BLOCKS"
fi

if [[ -n "$MAX_TXS_OVERRIDE" ]]; then
  MAX_TXS="$MAX_TXS_OVERRIDE"
  echo "Overriding max_transactions from command line: $MAX_TXS"
fi

if [[ "$ENVIRONMENT" != "athens" && "$ENVIRONMENT" != "mainnet" ]]; then
  echo "ERROR: Invalid environment: $ENVIRONMENT. Must be 'athens' or 'mainnet'."
  exit 1
fi
NODE_URL="https://${ENVIRONMENT}.rpc.zetachain.com:443/${API_KEY}/rpc"

echo "Environment: $ENVIRONMENT"
echo "Max blocks: $MAX_BLOCKS"
echo "Max transactions: $MAX_TXS"
echo "Config file: $CONFIG_FILE"
echo "Node URL: $NODE_URL"
echo "Output directory: $OUTPUT_DIR"

mkdir -p "$OUTPUT_DIR"

SUMMARY_FILE="$OUTPUT_DIR/tx_summary.csv"
FAILED_TX_FILE="$OUTPUT_DIR/failed_transactions.json"
ERROR_LOG_FILE="$OUTPUT_DIR/error.log"

echo "block_height,total_txs,success_txs,failed_txs,timestamp" > "$SUMMARY_FILE"

echo "[" > "$FAILED_TX_FILE"

TOTAL_BLOCKS=0
SUCCESSFUL_BLOCKS=0
TOTAL_TXS=0
SUCCESSFUL_TXS=0
FAILED_TXS=0
FIRST_JSON_ENTRY=true

echo "Getting current block height..."
CURRENT_HEIGHT=$(zetacored query block --node="$NODE_URL" --output=json latest | jq -r '.block.header.height')
echo "Current block height: $CURRENT_HEIGHT"

START_HEIGHT=$((CURRENT_HEIGHT - MAX_BLOCKS + 1))
if ((START_HEIGHT < 1)); then
  START_HEIGHT=1
fi
echo "Start block height: $START_HEIGHT"

for ((height=CURRENT_HEIGHT; height>=START_HEIGHT; height--)); do
  if ((TOTAL_TXS >= MAX_TXS)); then
    echo "Reached maximum transaction limit of $MAX_TXS"
    break
  fi

  echo "Processing block $height..."
  
  BLOCK_DATA=$(zetacored query block --node="$NODE_URL" --type=height --output=json "$height")
  TIMESTAMP=$(echo "$BLOCK_DATA" | jq -r '.block.header.time')
  TXS_IN_BLOCK=$(echo "$BLOCK_DATA" | jq '.block.data.txs | length')
  
  if [[ "$TXS_IN_BLOCK" -eq 0 ]]; then
    echo "Block $height has no transactions, skipping"
    continue
  fi

  echo "Block $height has $TXS_IN_BLOCK transactions"
  
  BLOCK_SUCCESS=0
  BLOCK_FAILURE=0
  
  for ((i=0; i<TXS_IN_BLOCK; i++)); do
    if ((TOTAL_TXS >= MAX_TXS)); then
      echo "Reached maximum transaction limit of $MAX_TXS within block $height"
      break
    fi
    
    TX_BASE64=$(echo "$BLOCK_DATA" | jq -r ".block.data.txs[$i]")
    TX_HASH=$(echo "$TX_BASE64" | base64 -d | sha256sum | awk '{print $1}' | tr '[:lower:]' '[:upper:]')
    
    echo "Querying transaction $TX_HASH from block $height (${i}/${TXS_IN_BLOCK})"
    
    if zetacored query tx "$TX_HASH" --node="$NODE_URL" --output=json > /dev/null 2>> "$ERROR_LOG_FILE"; then
      BLOCK_SUCCESS=$((BLOCK_SUCCESS + 1))
      SUCCESSFUL_TXS=$((SUCCESSFUL_TXS + 1))
    else
      BLOCK_FAILURE=$((BLOCK_FAILURE + 1))
      FAILED_TXS=$((FAILED_TXS + 1))
      
      if [[ "$FIRST_JSON_ENTRY" == "true" ]]; then
        FIRST_JSON_ENTRY=false
      else
        echo "," >> "$FAILED_TX_FILE"
      fi
      
      echo "  {\"block_height\":$height,\"tx_hash\":\"$TX_HASH\",\"tx_index\":$i}" >> "$FAILED_TX_FILE"
    fi
    
    TOTAL_TXS=$((TOTAL_TXS + 1))
  done
  
  echo "$height,$TXS_IN_BLOCK,$BLOCK_SUCCESS,$BLOCK_FAILURE,$TIMESTAMP" >> "$SUMMARY_FILE"
  
  if [[ "$BLOCK_FAILURE" -eq 0 ]]; then
    SUCCESSFUL_BLOCKS=$((SUCCESSFUL_BLOCKS + 1))
  fi
  
  TOTAL_BLOCKS=$((TOTAL_BLOCKS + 1))
done

echo "]" >> "$FAILED_TX_FILE"

SUMMARY_REPORT="$OUTPUT_DIR/summary.md"
{
  echo "# Transaction Query Summary"
  echo "- Network: $ENVIRONMENT"
  echo "- Run Time: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
  echo "- Total Blocks Processed: $TOTAL_BLOCKS"
  echo "- Blocks with All Transactions Successful: $SUCCESSFUL_BLOCKS"
  echo "- Total Transactions Processed: $TOTAL_TXS"
  echo "- Successful Transactions: $SUCCESSFUL_TXS"
  echo "- Failed Transactions: $FAILED_TXS"
  if [[ "$TOTAL_TXS" -gt 0 ]]; then
    echo "- Success Rate: $(( SUCCESSFUL_TXS * 100 / TOTAL_TXS ))%"
  else
    echo "- Success Rate: N/A (no transactions processed)"
  fi
} > "$SUMMARY_REPORT"

echo "Transaction validation complete. Summary written to $SUMMARY_REPORT"
