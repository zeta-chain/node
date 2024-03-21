import * as _12 from "./abci/types";
import * as _13 from "./crypto/keys";
import * as _14 from "./crypto/proof";
import * as _15 from "./libs/bits/types";
import * as _16 from "./p2p/types";
import * as _17 from "./types/block";
import * as _18 from "./types/evidence";
import * as _19 from "./types/params";
import * as _20 from "./types/types";
import * as _21 from "./types/validator";
import * as _22 from "./version/types";
export namespace tendermint {
  export const abci = {
    ..._12
  };
  export const crypto = {
    ..._13,
    ..._14
  };
  export namespace libs {
    export const bits = {
      ..._15
    };
  }
  export const p2p = {
    ..._16
  };
  export const types = {
    ..._17,
    ..._18,
    ..._19,
    ..._20,
    ..._21
  };
  export const version = {
    ..._22
  };
}