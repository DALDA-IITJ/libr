# 📜 State Transactions

## 🔥 1. Genesis of Community
```json
{
    "sender": "creatorPublicKey",
    "recipient": "0x0000000000000",
    "amt": 0,
    "data": {
        "type": "GENESIS",
        "name": "communityName",
        "metadata": "",
        "traffic": "",
        "modFaultTolerance": "",
        "dbReplicationFactor": ""
    },
    "nonce": senderNonce,
    "sign": "" // Treated as community public key
}
```

---

## 🏛 2. DB Joined
```json
{
    "sender": "dbPublicKey",
    "recipient": "communityPublicKey",
    "amt": 0,
    "data": {
        "type": "DB_JOINED",
        "metadata": {
            "ip": "",
            "port": "",
            "other": ""
        }
    },
    "nonce": senderNonce,
    "sign": ""
}
```

---

## 🚪 3. DB Left
```json
{
    "sender": "discovererPK",
    "recipient": "communityPublicKey",
    "amt": 0,
    "data": {
        "type": "DB_LEFT",
        "leaver": "dbPublicKey",
        "metadata": {
            "ip": "",
            "port": "",
            "other": ""
        }
    },
    "nonce": senderNonce,
    "sign": ""
}
```

---

## 🛡 4. Mod Joined
```json
{
    "sender": "modPublicKey",
    "recipient": "communityPublicKey",
    "amt": 0,
    "data": {
        "type": "MOD_JOINED",
        "metadata": {
            "ip": "",
            "port": "",
            "other": ""
        }
    },
    "nonce": senderNonce,
    "sign": ""
}
```

---

## 🚷 5. Mod Left
```json
{
    "sender": "discovererPK",
    "recipient": "communityPublicKey",
    "amt": 0,
    "data": {
        "type": "MOD_LEFT",
        "leaver": "modPublicKey",
        "metadata": {
            "ip": "",
            "port": "",
            "other": ""
        }
    },
    "nonce": senderNonce,
    "sign": ""
}
```

---

## ⚙️ Configurations

### 📡 Traffic (T)
`T = f(no. of msgs/time)`  
➝ Divide all **Unix Timestamps** by `T`

### 🛡 Mod Fault Tolerance (M)
`M = No. of Mods a message needs to be signed from`  
➝ Send to **2M+1** or **All < 2M+1**

### 💾 DB Replication Factor (R)
`R = No. of DBs a message is stored in initially`  
➝ Send to **R** or **All < R**

---

### ⚖️ Governance & Smart Contracts
For **editing community configuration**, only the **creator** has permission. However, this raises a **centralization problem**. The proposed **solution** is to implement **Governance & Smart Contracts** to allow decentralized decision-making.

