# Node Flow Documentation

## Endpoints

### `/getTime`

#### Flow:
1. **Client Request**: 
    - The client sends a timestamp (`ts`) to the server.
2. **Validation**: 
    - The server checks if the difference between the node's current time (`node_time`) and the client's timestamp (`client_ts`) is within a predefined limit (5 minutes).
3. **Response**: 
    - If valid, the server responds with a signed timestamp (`ts_sign`).

---

### `/receiveMsg`

#### Flow:
1. **Client Request**: 
    - The client sends a message certificate (`msg_cert`) to the server.
2. **Verification**:
    - The server verifies the structure of the certificate and unwraps it.
    - It ensures that all module certificates (`mod_cert[i].ts_cert`) and module messages (`mod[i].msgs`) are consistent.
3. **Sender Validation**: 
    - The server verifies the sender's signature.
4. **Timestamp Validation**: 
    - The server validates one timestamp certificate (`ts_cert`).
5. **Blockchain Validation**: 
    - The server receives and verifies the blockchain data.
6. **Database Validation**: 
    - The server checks the database consistency using `findDb`.
7. **Module Certificate Validation**: 
    - The server validates each module certificate (`mod_cert[i]`).
8. **Storage and Response**: 
    - The server stores the data and sends an appropriate response back to the client.

---

### `/sendMsg`

#### Flow:
1. **Client Request**: 
    - The client requests messages for a specific timestamp.
2. **Response**: 
    - The server retrieves and returns all messages corresponding to the given timestamp.

