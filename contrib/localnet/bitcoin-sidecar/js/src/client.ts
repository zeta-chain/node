import { initEccLib, payments, Psbt } from "bitcoinjs-lib";
import { bitcoin, Network, regtest } from "bitcoinjs-lib/src/networks";
import BIP32Factory, { BIP32Interface } from 'bip32';
import * as ecc from 'tiny-secp256k1';
import randomBytes from "randombytes";
import { ScriptBuilder } from "./script";
import { Taptree } from "bitcoinjs-lib/src/types";
import { toXOnly } from "./util";

const LEAF_VERSION_TAPSCRIPT = 0xc0;

initEccLib(ecc);
const bip32 = BIP32Factory(ecc);
const rng = randomBytes;

/// The evm address type, a 20 bytes hex string
export type Address = String;
export type BtcAddress = String;

/// The BTC transaction hash returned
export type BtcTxnHash = String;
export interface BtcInput {
    txn: BtcTxnHash,
    idx: number,
}

/**
 * The example client for interacting with ZetaChain in BTC. There are currently two ways
 * of calling a smart contract on ZetaChain from BTC:
 *
 * - Using OP_RETURN
 * - Using Witness
 *
 * The method used is now based on the data size. Within 80 bytes, `OP_RETURN` is used, else
 * the data is written to Witness.
 *
 * This class handles only the case where data is more than 80 bytes.
 */
export class ZetaBtcClient {
    /** The BTC network interracting with */
    readonly network: Network;

    private reveal: RevealTxnBuilder | null;

    private constructor(network: Network) {
        this.network = network;
    }

    public static regtest(): ZetaBtcClient {
        return new ZetaBtcClient(regtest);
    }

    public static mainnet(): ZetaBtcClient {
        return new ZetaBtcClient(bitcoin);
    }

    /**
     * Call a target address and passing the data call.
     *
     * @param address The target zetachain evm address
     * @param calldata The calldata that will be invoked on Zetachain
     */
    public call(
        address: Address,
        calldata: Buffer,
    ): Address {
        if (calldata.length <= 80) {
            throw Error("Use op return instead");
        }

        if (address.startsWith("0x")) {
            address = address.substring(2);
        }

        return this.callWithWitness(Buffer.concat([Buffer.from(address, "hex"), calldata]));
    }

    private callWithWitness(
        data: Buffer,
    ): Address {
        const internalKey = bip32.fromSeed(rng(64), this.network);

        const leafScript = this.genLeafScript(internalKey.publicKey, data);

        const scriptTree: Taptree = { output: leafScript };

        const { address: commitAddress } = payments.p2tr({
            internalPubkey: toXOnly(internalKey.publicKey),
            scriptTree,
            network: this.network,
        });

        this.reveal = new RevealTxnBuilder(internalKey, leafScript, this.network);

        return commitAddress;
    }

    public buildRevealTxn(to: string, commitTxn: BtcInput, commitAmount: number, feeRate: number): Buffer {
        if (this.reveal === null) {
            throw new Error("commit txn not built yet");
        }

        this.reveal.with_commit_tx(to, commitTxn, commitAmount, feeRate);
        return this.reveal.dump();
    }

    private genLeafScript(publicKey: Buffer, data: Buffer,): Buffer {
        const builder = ScriptBuilder.new(publicKey);
        builder.pushData(data);
        return builder.build();
    }
}

class RevealTxnBuilder {
    private psbt: Psbt;
    private key: BIP32Interface;
    private leafScript: Buffer;
    private network: Network

    constructor(key: BIP32Interface, leafScript: Buffer, network: Network) {
        this.psbt = new Psbt({ network });;
        this.key = key;
        this.leafScript = leafScript;
        this.network = network;
    }

    public with_commit_tx(to: string, commitTxn: BtcInput, commitAmount: number, feeRate: number): RevealTxnBuilder {
        const scriptTree: Taptree = { output: this.leafScript };

        const { output, witness } = payments.p2tr({
            internalPubkey: toXOnly(this.key.publicKey),
            scriptTree,
            redeem: {
                output: this.leafScript,
                redeemVersion: LEAF_VERSION_TAPSCRIPT,
            },
            network: this.network,
        });

        this.psbt.addInput({
            hash: commitTxn.txn.toString(),
            index: commitTxn.idx,
            witnessUtxo: { value: commitAmount, script: output! },
            tapLeafScript: [
                {
                    leafVersion: LEAF_VERSION_TAPSCRIPT,
                    script: this.leafScript,
                    controlBlock: witness![witness!.length - 1],
                },
            ],
        });

        this.psbt.addOutput({
            value: commitAmount - this.estimateFee(to, commitAmount, feeRate),
            address: to,
        });

        this.psbt.signAllInputs(this.key);
        this.psbt.finalizeAllInputs();

        return this;
    }

    public dump(): Buffer {
        return this.psbt.extractTransaction(true).toBuffer();
    }

    private estimateFee(to: string, amount: number, feeRate: number): number {
        const cloned = this.psbt.clone();

        cloned.addOutput({
            value: amount,
            address: to,
        });

        // should have a way to avoid signing but just providing mocked signautre
        cloned.signAllInputs(this.key);
        cloned.finalizeAllInputs();

        const size = cloned.extractTransaction().virtualSize();
        return size * feeRate;
    }
}