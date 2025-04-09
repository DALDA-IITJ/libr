const express = require("express");
const bodyParser = require("body-parser");
const crypto = require("crypto");
const dotenv = require("dotenv");
const { ec: EC } = require("elliptic");

dotenv.config();

const app = express();
const PORT = process.env.PORT || 8080;

// Setup secp256k1
const ec = new EC("secp256k1");

// Load private key from .env
const PRIVATE_KEY = process.env.PRIVATE_KEY;
if (!PRIVATE_KEY) throw new Error("Missing PRIVATE_KEY in .env");

// Create key pair from private key
const keyPair = ec.keyFromPrivate(PRIVATE_KEY);

// Get compressed public key (for response)
const publicKey = keyPair.getPublic(true, "hex");

app.use(bodyParser.json());

app.post("/moderate", (req, res) => {
    const { message, timestamp } = req.body;

    // Construct payload and hash it
    const payload = { message, timestamp };
    const dataToSign = JSON.stringify(payload);
    const hash = crypto.createHash("sha256").update(dataToSign).digest("hex");

    // Sign the hash
    const signature = keyPair.sign(hash, "hex");
    const signatureHex = signature.toDER("hex");

    console.log("signing message ", message, " with sign ", signatureHex);

    // Respond with signature and public key
    res.status(200).json({
        public_key: publicKey,
        sign: signatureHex,
    });
});

app.post("/db/savemsg", (req, res) => {
    console.log(req.body);
    res.status(200).json();
})

app.listen(PORT, () => {
    console.log(`Moderator running on http://localhost:${PORT}`);
});
