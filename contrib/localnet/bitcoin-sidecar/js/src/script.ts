import { opcodes, script, Stack } from "bitcoinjs-lib";
import { toXOnly } from "./util";

const MAX_SCRIPT_ELEMENT_SIZE = 520;

/** The tapscript builder for zetaclient spending script */
export class ScriptBuilder {
    private script: Stack;

    private constructor(initialScript: Stack) {
        this.script = initialScript;
    }

    public static new(publicKey: Buffer): ScriptBuilder {
        const stack = [
            toXOnly(publicKey),
            opcodes.OP_CHECKSIG,
        ];
        return new ScriptBuilder(stack);
    }

    public pushData(data: Buffer) {
        if (data.length <= 80) {
            throw new Error("data length should be more than 80 bytes");
        }

        this.script.push(
            opcodes.OP_FALSE,
            opcodes.OP_IF
        );

        const chunks = chunkBuffer(data, MAX_SCRIPT_ELEMENT_SIZE);
        for (const chunk of chunks) {
            this.script.push(chunk);
        }

        this.script.push(opcodes.OP_ENDIF);
    }

    public build(): Buffer {
        return script.compile(this.script);
    }
}

function chunkBuffer(buffer: Buffer, chunkSize: number): Buffer[] {
    const chunks = [];
    for (let i = 0; i < buffer.length; i += chunkSize) {
      const chunk = buffer.slice(i, i + chunkSize);
      chunks.push(chunk);
    }
    return chunks;
  }