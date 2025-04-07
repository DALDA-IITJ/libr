# Node Flow Documentation

## Endpoints


### `/receiveMsg`

#### Flow:
1. **Client Request**: 
    - The client sends a message certificate (`msg_cert`) to the server.
2. **Verification**:
    - The server verifies the structure of the certificate and unwraps it.
    - It ensures that all module timestamps (`mod_cert[i].ts`) and module messages (`mod[i].msgs`) are consistent.
3. **Sender Validation**: 
    - The server verifies the sender's signature.
4. **Blockchain Validation**: 
    - The server receives and verifies the blockchain data.
5. **Database Validation**: 
    - The server checks the database consistency using `findDb`.
7. **Module Certificate Validation**: 
    - The server validates each module certificate (`mod_cert[i]`) and each mod is different.
8. **Server side Time verification with buffer**:
    - The server verifies the time using a buffer.
9. **Storage and Response**: 
    - The server stores the data and sends an appropriate response back to the client.

---

### `/sendMsg`

#### Flow:
1. **Client Request**: 
    - The client requests messages for a specific timestamp.
2. **Response**: 
    - The server retrieves and returns all messages corresponding to the given timestamp.

