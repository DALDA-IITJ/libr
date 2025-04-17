<div align="center">

# LIBR

[<img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" width="60">](https://golang.org)

[![Open in Visual Studio Code](https://img.shields.io/badge/Open%20in%20VS%20Code-007ACC?logo=visual-studio-code&logoColor=white)](https://vscode.dev/)
[![Contributors](https://img.shields.io/github/contributors/DALDA-IITJ/libr)](https://github.com/DALDA-IITJ/libr/graphs/contributors)
[![Forks](https://img.shields.io/github/forks/DALDA-IITJ/libr?style=social)](https://github.com/DALDA-IITJ/libr/network/members)
[![Stars](https://img.shields.io/github/stars/DALDA-IITJ/libr?style=social)](https://github.com/DALDA-IITJ/libr/stargazers)
[![License](https://img.shields.io/github/license/DALDA-IITJ/libr)](https://github.com/DALDA-IITJ/libr/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-v1.16+-blue?logo=go&logoColor=white)](https://golang.org/)

*A Decentralized Messaging Platform with Content Moderation* 📢🚫

[Key Features](#key-features) • [Installation](#installation) • [Documentation](#api-endpoints) • [Contributing](#contributing)

</div>

## 🔄 Overview

LIBR is a decentralized messaging platform that combines blockchain principles with content moderation to create a secure, distributed communication system. Built with Go for backend services and React for the user interface, LIBR provides a robust framework for verified message exchange with cryptographic proof of moderation and authenticity.

## 🔐 Key Features

- **Decentralized Messaging** 🌐: Exchange messages across a distributed network with cryptographic verification.
- **Content Moderation** 🪖: Multi-level content moderation with configurable thresholds for different categories of inappropriate content.
- **Signature Verification** ✅: Messages are cryptographically signed by both users and moderators, ensuring authenticity and approval.
- **Distributed Architecture** 🧰: Multiple moderator nodes validate content independently to prevent censorship and single points of failure.
- **Cryptographic Security** 🔒: RSA-based cryptography for message signing and verification.
- **Time-Based Message Bucketing** ⏰: Messages are grouped into time buckets for efficient retrieval and organization.

> 🚀 **Note**: LIBR is designed for environments where content moderation is necessary while maintaining the benefits of decentralization. It's ideal for educational institutions, private organizations, and communities seeking a balance between free expression and responsible content management.

## ⚙️ Prerequisites

- **Go** (v1.16 or higher)
- **Node.js** (v14 or higher)
- **Google Cloud Platform API key** (for content moderation services)
- **npm** or **yarn** package manager

## 🚀 Running the System

1. Start moderator nodes:
```bash
cd modules/mod
./mod
# Default port is 4000, can be configured in .env
```

2. Start client application:
```bash
cd modules/client
go run .
```

3. Start UI application:
```bash
cd modules/UI
npm start
```

## 📂 API Endpoints

### 🔧 Moderator Service

#### 🔒 Content Moderation
Send a message for moderation approval.
```bash
curl -X POST http://localhost:4000/moderate \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Your message content here",
    "timestamp": "1681234567"
  }'
```

#### 📊 Response Format
```json
{
  "public_key": "moderator_public_key",
  "sign": "digital_signature_of_approved_content"
}
```

### 📲 Client Service

#### 📢 Send Message
Submit a new message to the network.
```bash
# Endpoint implementation details in client module
```

#### 📃 Retrieve Messages
Get messages from a specific time bucket.
```bash
# Endpoint implementation details in client module
```

## 🚧 System Architecture

LIBR consists of several interconnected modules:

1. **Moderator Nodes** 🔧: Analyze message content for inappropriate material using Google Cloud Natural Language API.
2. **Client Module** 📡: Handles message creation, interaction with moderators, and storage operations.
3. **Core Module** 🔄: Provides cryptographic functions, configuration management, and blockchain interfaces.
4. **UI Module** 🎨: User interface for interacting with the messaging platform.
5. **Storage Layer** 📂: Distributed storage system for persisting signed messages.

## 🕵️‍ Content Moderation

LIBR supports customizable content moderation:

- Currently implemented using Google Cloud Natural Language API
- User-dependent moderation that can be customized based on needs
- Multiple moderation categories supported (toxic content, insults, profanity, etc.)
- Users can implement their own moderation services or use the provided API

To modify moderation settings, users can edit the configuration in the file:
```
modules/mod/config.go
```

## 📚 Contributing

1. Fork the repository to start working on your changes.
2. Create a feature branch.
3. Commit your changes.
4. Push to the branch.
5. Create a Pull Request to merge your changes.

## 🌐 License

This project is licensed under the terms of the MIT license. See [LICENSE](LICENSE) for more details.

## 👨‍💼 Team

LIBR is developed by the DALDA team at IIT Jodhpur. 🎓🌟

