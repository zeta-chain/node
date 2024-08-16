import { ZetaBtcClient } from "./client";
import express, { Request, Response } from 'express';

const app = express();
const PORT = process.env.PORT || 3000;
let zetaClient = ZetaBtcClient.regtest();

app.use(express.json());

// Middleware to parse URL-encoded bodies
app.use(express.urlencoded({ extended: true }));

// Route to handle JSON POST requests
app.post('/commit', (req: Request, res: Response) => {
    const memo: string = req.body.memo;
    const address = zetaClient.call("", Buffer.from(memo, "hex"));
    res.json({ address });
});

// Route to handle URL-encoded POST requests
app.post('/reveal', (req: Request, res: Response) => {
    const { txn, idx, amount, feeRate, to } = req.body;
    console.log(txn, idx, amount, feeRate);

    const rawHex = zetaClient.buildRevealTxn(to,{ txn, idx }, Number(amount), feeRate).toString("hex");
    zetaClient = ZetaBtcClient.regtest();
    res.json({ rawHex });
});

// Start the server
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});

/**
 * curl --request POST --header "Content-Type: application/json" --data '{"memo":"72f080c854647755d0d9e6f6821f6931f855b9acffd53d87433395672756d58822fd143360762109ab898626556b1c3b8d3096d2361f1297df4a41c1b429471a9aa2fc9be5f27c13b3863d6ac269e4b587d8389f8fd9649859935b0d48dea88cdb40f20c"}' http://localhost:3000/commit
 * curl --request POST --header "Content-Type: application/json" --data '{"txn": "7a57f987a3cb605896a5909d9ef2bf7afbf0c78f21e4118b85d00d9e4cce0c2c", "idx": 0, "amount": 1000, "feeRate": 10}' http://localhost:3000/reveal
 */