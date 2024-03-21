import { ObservationType, observationTypeFromJSON } from "./observer";
import { BinaryReader, BinaryWriter } from "../../binary";
import { Decimal } from "@cosmjs/math";
export enum VoteType {
  SuccessObservation = 0,
  /** FailureObservation - Failure observation means , the the message that this voter is observing failed / reverted . It does not mean it was unable to observe. */
  FailureObservation = 1,
  NotYetVoted = 2,
  UNRECOGNIZED = -1,
}
export const VoteTypeSDKType = VoteType;
export const VoteTypeAmino = VoteType;
export function voteTypeFromJSON(object: any): VoteType {
  switch (object) {
    case 0:
    case "SuccessObservation":
      return VoteType.SuccessObservation;
    case 1:
    case "FailureObservation":
      return VoteType.FailureObservation;
    case 2:
    case "NotYetVoted":
      return VoteType.NotYetVoted;
    case -1:
    case "UNRECOGNIZED":
    default:
      return VoteType.UNRECOGNIZED;
  }
}
export function voteTypeToJSON(object: VoteType): string {
  switch (object) {
    case VoteType.SuccessObservation:
      return "SuccessObservation";
    case VoteType.FailureObservation:
      return "FailureObservation";
    case VoteType.NotYetVoted:
      return "NotYetVoted";
    case VoteType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export enum BallotStatus {
  BallotFinalized_SuccessObservation = 0,
  BallotFinalized_FailureObservation = 1,
  BallotInProgress = 2,
  UNRECOGNIZED = -1,
}
export const BallotStatusSDKType = BallotStatus;
export const BallotStatusAmino = BallotStatus;
export function ballotStatusFromJSON(object: any): BallotStatus {
  switch (object) {
    case 0:
    case "BallotFinalized_SuccessObservation":
      return BallotStatus.BallotFinalized_SuccessObservation;
    case 1:
    case "BallotFinalized_FailureObservation":
      return BallotStatus.BallotFinalized_FailureObservation;
    case 2:
    case "BallotInProgress":
      return BallotStatus.BallotInProgress;
    case -1:
    case "UNRECOGNIZED":
    default:
      return BallotStatus.UNRECOGNIZED;
  }
}
export function ballotStatusToJSON(object: BallotStatus): string {
  switch (object) {
    case BallotStatus.BallotFinalized_SuccessObservation:
      return "BallotFinalized_SuccessObservation";
    case BallotStatus.BallotFinalized_FailureObservation:
      return "BallotFinalized_FailureObservation";
    case BallotStatus.BallotInProgress:
      return "BallotInProgress";
    case BallotStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface Ballot {
  index: string;
  ballotIdentifier: string;
  voterList: string[];
  votes: VoteType[];
  observationType: ObservationType;
  ballotThreshold: string;
  ballotStatus: BallotStatus;
  ballotCreationHeight: bigint;
}
export interface BallotProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.Ballot";
  value: Uint8Array;
}
export interface BallotAmino {
  index?: string;
  ballot_identifier?: string;
  voter_list?: string[];
  votes?: VoteType[];
  observation_type?: ObservationType;
  ballot_threshold?: string;
  ballot_status?: BallotStatus;
  ballot_creation_height?: string;
}
export interface BallotAminoMsg {
  type: "/zetachain.zetacore.observer.Ballot";
  value: BallotAmino;
}
export interface BallotSDKType {
  index: string;
  ballot_identifier: string;
  voter_list: string[];
  votes: VoteType[];
  observation_type: ObservationType;
  ballot_threshold: string;
  ballot_status: BallotStatus;
  ballot_creation_height: bigint;
}
export interface BallotListForHeight {
  height: bigint;
  ballotsIndexList: string[];
}
export interface BallotListForHeightProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.BallotListForHeight";
  value: Uint8Array;
}
export interface BallotListForHeightAmino {
  height?: string;
  ballots_index_list?: string[];
}
export interface BallotListForHeightAminoMsg {
  type: "/zetachain.zetacore.observer.BallotListForHeight";
  value: BallotListForHeightAmino;
}
export interface BallotListForHeightSDKType {
  height: bigint;
  ballots_index_list: string[];
}
function createBaseBallot(): Ballot {
  return {
    index: "",
    ballotIdentifier: "",
    voterList: [],
    votes: [],
    observationType: 0,
    ballotThreshold: "",
    ballotStatus: 0,
    ballotCreationHeight: BigInt(0)
  };
}
export const Ballot = {
  typeUrl: "/zetachain.zetacore.observer.Ballot",
  encode(message: Ballot, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    if (message.ballotIdentifier !== "") {
      writer.uint32(18).string(message.ballotIdentifier);
    }
    for (const v of message.voterList) {
      writer.uint32(26).string(v!);
    }
    writer.uint32(34).fork();
    for (const v of message.votes) {
      writer.int32(v);
    }
    writer.ldelim();
    if (message.observationType !== 0) {
      writer.uint32(40).int32(message.observationType);
    }
    if (message.ballotThreshold !== "") {
      writer.uint32(50).string(Decimal.fromUserInput(message.ballotThreshold, 18).atomics);
    }
    if (message.ballotStatus !== 0) {
      writer.uint32(56).int32(message.ballotStatus);
    }
    if (message.ballotCreationHeight !== BigInt(0)) {
      writer.uint32(64).int64(message.ballotCreationHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Ballot {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBallot();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        case 2:
          message.ballotIdentifier = reader.string();
          break;
        case 3:
          message.voterList.push(reader.string());
          break;
        case 4:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.votes.push((reader.int32() as any));
            }
          } else {
            message.votes.push((reader.int32() as any));
          }
          break;
        case 5:
          message.observationType = (reader.int32() as any);
          break;
        case 6:
          message.ballotThreshold = Decimal.fromAtomics(reader.string(), 18).toString();
          break;
        case 7:
          message.ballotStatus = (reader.int32() as any);
          break;
        case 8:
          message.ballotCreationHeight = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Ballot>): Ballot {
    const message = createBaseBallot();
    message.index = object.index ?? "";
    message.ballotIdentifier = object.ballotIdentifier ?? "";
    message.voterList = object.voterList?.map(e => e) || [];
    message.votes = object.votes?.map(e => e) || [];
    message.observationType = object.observationType ?? 0;
    message.ballotThreshold = object.ballotThreshold ?? "";
    message.ballotStatus = object.ballotStatus ?? 0;
    message.ballotCreationHeight = object.ballotCreationHeight !== undefined && object.ballotCreationHeight !== null ? BigInt(object.ballotCreationHeight.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: BallotAmino): Ballot {
    const message = createBaseBallot();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.ballot_identifier !== undefined && object.ballot_identifier !== null) {
      message.ballotIdentifier = object.ballot_identifier;
    }
    message.voterList = object.voter_list?.map(e => e) || [];
    message.votes = object.votes?.map(e => voteTypeFromJSON(e)) || [];
    if (object.observation_type !== undefined && object.observation_type !== null) {
      message.observationType = observationTypeFromJSON(object.observation_type);
    }
    if (object.ballot_threshold !== undefined && object.ballot_threshold !== null) {
      message.ballotThreshold = object.ballot_threshold;
    }
    if (object.ballot_status !== undefined && object.ballot_status !== null) {
      message.ballotStatus = ballotStatusFromJSON(object.ballot_status);
    }
    if (object.ballot_creation_height !== undefined && object.ballot_creation_height !== null) {
      message.ballotCreationHeight = BigInt(object.ballot_creation_height);
    }
    return message;
  },
  toAmino(message: Ballot): BallotAmino {
    const obj: any = {};
    obj.index = message.index;
    obj.ballot_identifier = message.ballotIdentifier;
    if (message.voterList) {
      obj.voter_list = message.voterList.map(e => e);
    } else {
      obj.voter_list = [];
    }
    if (message.votes) {
      obj.votes = message.votes.map(e => e);
    } else {
      obj.votes = [];
    }
    obj.observation_type = message.observationType;
    obj.ballot_threshold = message.ballotThreshold;
    obj.ballot_status = message.ballotStatus;
    obj.ballot_creation_height = message.ballotCreationHeight ? message.ballotCreationHeight.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: BallotAminoMsg): Ballot {
    return Ballot.fromAmino(object.value);
  },
  fromProtoMsg(message: BallotProtoMsg): Ballot {
    return Ballot.decode(message.value);
  },
  toProto(message: Ballot): Uint8Array {
    return Ballot.encode(message).finish();
  },
  toProtoMsg(message: Ballot): BallotProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.Ballot",
      value: Ballot.encode(message).finish()
    };
  }
};
function createBaseBallotListForHeight(): BallotListForHeight {
  return {
    height: BigInt(0),
    ballotsIndexList: []
  };
}
export const BallotListForHeight = {
  typeUrl: "/zetachain.zetacore.observer.BallotListForHeight",
  encode(message: BallotListForHeight, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.height !== BigInt(0)) {
      writer.uint32(8).int64(message.height);
    }
    for (const v of message.ballotsIndexList) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BallotListForHeight {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBallotListForHeight();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.height = reader.int64();
          break;
        case 2:
          message.ballotsIndexList.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BallotListForHeight>): BallotListForHeight {
    const message = createBaseBallotListForHeight();
    message.height = object.height !== undefined && object.height !== null ? BigInt(object.height.toString()) : BigInt(0);
    message.ballotsIndexList = object.ballotsIndexList?.map(e => e) || [];
    return message;
  },
  fromAmino(object: BallotListForHeightAmino): BallotListForHeight {
    const message = createBaseBallotListForHeight();
    if (object.height !== undefined && object.height !== null) {
      message.height = BigInt(object.height);
    }
    message.ballotsIndexList = object.ballots_index_list?.map(e => e) || [];
    return message;
  },
  toAmino(message: BallotListForHeight): BallotListForHeightAmino {
    const obj: any = {};
    obj.height = message.height ? message.height.toString() : undefined;
    if (message.ballotsIndexList) {
      obj.ballots_index_list = message.ballotsIndexList.map(e => e);
    } else {
      obj.ballots_index_list = [];
    }
    return obj;
  },
  fromAminoMsg(object: BallotListForHeightAminoMsg): BallotListForHeight {
    return BallotListForHeight.fromAmino(object.value);
  },
  fromProtoMsg(message: BallotListForHeightProtoMsg): BallotListForHeight {
    return BallotListForHeight.decode(message.value);
  },
  toProto(message: BallotListForHeight): Uint8Array {
    return BallotListForHeight.encode(message).finish();
  },
  toProtoMsg(message: BallotListForHeight): BallotListForHeightProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.BallotListForHeight",
      value: BallotListForHeight.encode(message).finish()
    };
  }
};