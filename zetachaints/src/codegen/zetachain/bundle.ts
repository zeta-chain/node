import * as _23 from "./authority/genesis";
import * as _24 from "./authority/policies";
import * as _25 from "./authority/query";
import * as _26 from "./authority/tx";
import * as _27 from "./crosschain/cross_chain_tx";
import * as _28 from "./crosschain/events";
import * as _29 from "./crosschain/gas_price";
import * as _30 from "./crosschain/genesis";
import * as _31 from "./crosschain/in_tx_hash_to_cctx";
import * as _32 from "./crosschain/in_tx_tracker";
import * as _33 from "./crosschain/last_block_height";
import * as _34 from "./crosschain/out_tx_tracker";
import * as _35 from "./crosschain/params";
import * as _36 from "./crosschain/query";
import * as _37 from "./crosschain/tx";
import * as _38 from "./emissions/events";
import * as _39 from "./emissions/genesis";
import * as _40 from "./emissions/params";
import * as _41 from "./emissions/query";
import * as _42 from "./emissions/tx";
import * as _43 from "./emissions/withdrawable_emissions";
import * as _44 from "./fungible/events";
import * as _45 from "./fungible/foreign_coins";
import * as _46 from "./fungible/genesis";
import * as _47 from "./fungible/params";
import * as _48 from "./fungible/query";
import * as _49 from "./fungible/system_contract";
import * as _50 from "./fungible/tx";
import * as _51 from "./observer/ballot";
import * as _52 from "./observer/blame";
import * as _53 from "./observer/block_header";
import * as _54 from "./observer/chain_nonces";
import * as _55 from "./observer/crosschain_flags";
import * as _56 from "./observer/events";
import * as _57 from "./observer/genesis";
import * as _58 from "./observer/keygen";
import * as _59 from "./observer/node_account";
import * as _60 from "./observer/nonce_to_cctx";
import * as _61 from "./observer/observer";
import * as _62 from "./observer/params";
import * as _63 from "./observer/pending_nonces";
import * as _64 from "./observer/query";
import * as _65 from "./observer/tss_funds_migrator";
import * as _66 from "./observer/tss";
import * as _67 from "./observer/tx";
import * as _71 from "./authority/tx.amino";
import * as _72 from "./crosschain/tx.amino";
import * as _73 from "./emissions/tx.amino";
import * as _74 from "./fungible/tx.amino";
import * as _75 from "./observer/tx.amino";
import * as _76 from "./authority/tx.registry";
import * as _77 from "./crosschain/tx.registry";
import * as _78 from "./emissions/tx.registry";
import * as _79 from "./fungible/tx.registry";
import * as _80 from "./observer/tx.registry";
import * as _81 from "./authority/query.rpc.Query";
import * as _82 from "./crosschain/query.rpc.Query";
import * as _83 from "./emissions/query.rpc.Query";
import * as _84 from "./fungible/query.rpc.Query";
import * as _85 from "./observer/query.rpc.Query";
import * as _86 from "./authority/tx.rpc.msg";
import * as _87 from "./crosschain/tx.rpc.msg";
import * as _88 from "./emissions/tx.rpc.msg";
import * as _89 from "./fungible/tx.rpc.msg";
import * as _90 from "./observer/tx.rpc.msg";
import * as _91 from "./rpc.query";
import * as _92 from "./rpc.tx";
export namespace zetachain {
  export namespace zetacore {
    export const authority = {
      ..._23,
      ..._24,
      ..._25,
      ..._26,
      ..._71,
      ..._76,
      ..._81,
      ..._86
    };
    export const crosschain = {
      ..._27,
      ..._28,
      ..._29,
      ..._30,
      ..._31,
      ..._32,
      ..._33,
      ..._34,
      ..._35,
      ..._36,
      ..._37,
      ..._72,
      ..._77,
      ..._82,
      ..._87
    };
    export const emissions = {
      ..._38,
      ..._39,
      ..._40,
      ..._41,
      ..._42,
      ..._43,
      ..._73,
      ..._78,
      ..._83,
      ..._88
    };
    export const fungible = {
      ..._44,
      ..._45,
      ..._46,
      ..._47,
      ..._48,
      ..._49,
      ..._50,
      ..._74,
      ..._79,
      ..._84,
      ..._89
    };
    export const observer = {
      ..._51,
      ..._52,
      ..._53,
      ..._54,
      ..._55,
      ..._56,
      ..._57,
      ..._58,
      ..._59,
      ..._60,
      ..._61,
      ..._62,
      ..._63,
      ..._64,
      ..._65,
      ..._66,
      ..._67,
      ..._75,
      ..._80,
      ..._85,
      ..._90
    };
  }
  export const ClientFactory = {
    ..._91,
    ..._92
  };
}