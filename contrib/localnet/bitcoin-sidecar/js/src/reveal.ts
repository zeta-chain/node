import { ZetaBtcClient } from "./client";

function main() {
    const client = ZetaBtcClient.regtest();

    const data = Buffer.alloc(600);

    const d = client.call("", data);

    console.log(d);
}

main();