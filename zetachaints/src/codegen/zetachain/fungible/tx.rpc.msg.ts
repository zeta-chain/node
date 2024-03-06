import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { MsgDeploySystemContracts, MsgDeploySystemContractsResponse, MsgDeployFungibleCoinZRC20, MsgDeployFungibleCoinZRC20Response, MsgRemoveForeignCoin, MsgRemoveForeignCoinResponse, MsgUpdateSystemContract, MsgUpdateSystemContractResponse, MsgUpdateContractBytecode, MsgUpdateContractBytecodeResponse, MsgUpdateZRC20WithdrawFee, MsgUpdateZRC20WithdrawFeeResponse, MsgUpdateZRC20PausedStatus, MsgUpdateZRC20PausedStatusResponse, MsgUpdateZRC20LiquidityCap, MsgUpdateZRC20LiquidityCapResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
  deploySystemContracts(request: MsgDeploySystemContracts): Promise<MsgDeploySystemContractsResponse>;
  deployFungibleCoinZRC20(request: MsgDeployFungibleCoinZRC20): Promise<MsgDeployFungibleCoinZRC20Response>;
  removeForeignCoin(request: MsgRemoveForeignCoin): Promise<MsgRemoveForeignCoinResponse>;
  updateSystemContract(request: MsgUpdateSystemContract): Promise<MsgUpdateSystemContractResponse>;
  updateContractBytecode(request: MsgUpdateContractBytecode): Promise<MsgUpdateContractBytecodeResponse>;
  updateZRC20WithdrawFee(request: MsgUpdateZRC20WithdrawFee): Promise<MsgUpdateZRC20WithdrawFeeResponse>;
  updateZRC20PausedStatus(request: MsgUpdateZRC20PausedStatus): Promise<MsgUpdateZRC20PausedStatusResponse>;
  updateZRC20LiquidityCap(request: MsgUpdateZRC20LiquidityCap): Promise<MsgUpdateZRC20LiquidityCapResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.deploySystemContracts = this.deploySystemContracts.bind(this);
    this.deployFungibleCoinZRC20 = this.deployFungibleCoinZRC20.bind(this);
    this.removeForeignCoin = this.removeForeignCoin.bind(this);
    this.updateSystemContract = this.updateSystemContract.bind(this);
    this.updateContractBytecode = this.updateContractBytecode.bind(this);
    this.updateZRC20WithdrawFee = this.updateZRC20WithdrawFee.bind(this);
    this.updateZRC20PausedStatus = this.updateZRC20PausedStatus.bind(this);
    this.updateZRC20LiquidityCap = this.updateZRC20LiquidityCap.bind(this);
  }
  deploySystemContracts(request: MsgDeploySystemContracts): Promise<MsgDeploySystemContractsResponse> {
    const data = MsgDeploySystemContracts.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "DeploySystemContracts", data);
    return promise.then(data => MsgDeploySystemContractsResponse.decode(new BinaryReader(data)));
  }
  deployFungibleCoinZRC20(request: MsgDeployFungibleCoinZRC20): Promise<MsgDeployFungibleCoinZRC20Response> {
    const data = MsgDeployFungibleCoinZRC20.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "DeployFungibleCoinZRC20", data);
    return promise.then(data => MsgDeployFungibleCoinZRC20Response.decode(new BinaryReader(data)));
  }
  removeForeignCoin(request: MsgRemoveForeignCoin): Promise<MsgRemoveForeignCoinResponse> {
    const data = MsgRemoveForeignCoin.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "RemoveForeignCoin", data);
    return promise.then(data => MsgRemoveForeignCoinResponse.decode(new BinaryReader(data)));
  }
  updateSystemContract(request: MsgUpdateSystemContract): Promise<MsgUpdateSystemContractResponse> {
    const data = MsgUpdateSystemContract.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "UpdateSystemContract", data);
    return promise.then(data => MsgUpdateSystemContractResponse.decode(new BinaryReader(data)));
  }
  updateContractBytecode(request: MsgUpdateContractBytecode): Promise<MsgUpdateContractBytecodeResponse> {
    const data = MsgUpdateContractBytecode.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "UpdateContractBytecode", data);
    return promise.then(data => MsgUpdateContractBytecodeResponse.decode(new BinaryReader(data)));
  }
  updateZRC20WithdrawFee(request: MsgUpdateZRC20WithdrawFee): Promise<MsgUpdateZRC20WithdrawFeeResponse> {
    const data = MsgUpdateZRC20WithdrawFee.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "UpdateZRC20WithdrawFee", data);
    return promise.then(data => MsgUpdateZRC20WithdrawFeeResponse.decode(new BinaryReader(data)));
  }
  updateZRC20PausedStatus(request: MsgUpdateZRC20PausedStatus): Promise<MsgUpdateZRC20PausedStatusResponse> {
    const data = MsgUpdateZRC20PausedStatus.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "UpdateZRC20PausedStatus", data);
    return promise.then(data => MsgUpdateZRC20PausedStatusResponse.decode(new BinaryReader(data)));
  }
  updateZRC20LiquidityCap(request: MsgUpdateZRC20LiquidityCap): Promise<MsgUpdateZRC20LiquidityCapResponse> {
    const data = MsgUpdateZRC20LiquidityCap.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Msg", "UpdateZRC20LiquidityCap", data);
    return promise.then(data => MsgUpdateZRC20LiquidityCapResponse.decode(new BinaryReader(data)));
  }
}