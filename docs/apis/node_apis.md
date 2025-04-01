# **DB Node API Documentation**

## **Base URL**
```
http://<dbnode-ip>:8080
```

---

## **1. Save Message**
### **Endpoint:**
```
POST /savemsg
```
### **Description:**
This endpoint receives a message with its associated certificates, validates it, and saves it to the database if valid.

### **Request Body:**
```json
{
  "sender": "senderpublickey",
  "mod_cert": [
    {
      "mod_sign": {
        "msg": "msg",
        "ts_cert": [
          {
            "ts_sign": {
              "client_ts": "cts",
              "db_ts": "nts",
              "sign": "db_sign"
            }
          }
        ],
        "sign": "mod_sign"
      }
    }
  ],
  "sign": "client_sign"
}
```

### **Response:**
#### **Success (Message Saved):**
```json
{
  "status": "success",
  "message": "Message stored successfully."
}
```
#### **Failure (Invalid Certificate):**
```json
{
  "status": "error",
  "message": "Invalid certificate."
}
```

---

## **2. Retrieve Messages by Timestamp**
### **Endpoint:**
```
GET /getmsg?time=<timestamp>
```
### **Description:**
Fetches all messages stored in the database that match the given timestamp.

### **Query Parameters:**
| Parameter  | Type     | Required | Description                       |
|------------|---------|----------|-----------------------------------|
| `time`     | String  | Yes      | Timestamp for filtering messages |

### **Response:**
#### **Success:**
```json
{
  "status": "success",
  "messages": [
    {
      "sender": "senderpublickey",
      "content": "msg",
      "timestamp": "timestamp"
    }
  ]
}
```
#### **Failure (No Messages Found):**
```json
{
  "status": "error",
  "message": "No messages found for the given timestamp."
}
```

---

## **3. Check Node Health**
### **Endpoint:**
```
GET /isalive
```
### **Description:**
Checks if the DB Node is active and running.

### **Response:**
#### **Success:**
- **HTTP Status Code:** `200 OK`
- **Response Body:** None (empty response)

---

## **4. Get Time Signature**
### **Endpoint:**
```
GET /timesign?time=<cts>
```
### **Description:**
Provides a signed timestamp from the DB Node for a given client timestamp.

### **Query Parameters:**
| Parameter  | Type     | Required | Description                          |
|------------|---------|----------|--------------------------------------|
| `time`     | String  | Yes      | The client timestamp (`cts`)        |

### **Response:**
#### **Success:**
```json
{
  "ts_sign": {
    "client_ts": "cts",
    "node_ts": "nts",
    "sign": "sign"
  }
}
```

#### **Failure (Invalid Request):**
```json
{
  "status": "error",
  "message": "Invalid timestamp format."
}
```

---
## **Error Codes & Handling**

| Error Code | HTTP Status | Description                        |
|------------|------------|------------------------------------|
| `400`      | Bad Request | Invalid input format or parameters |
| `401`      | Unauthorized | Invalid signature or certificate  |
| `404`      | Not Found   | Requested resource not found      |
| `500`      | Server Error | Unexpected server failure         |

---

## **Security Considerations**
- All messages must be signed and verified before storage.
- Timestamp signing ensures message authenticity.
- API should implement **rate limiting** and **authentication** to prevent abuse.

---

This document outlines the core API functionality for the **DB Node** module in the **LIBR Project**. 🚀

