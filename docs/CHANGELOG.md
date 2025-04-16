# Changelog

## [0.1.0] - 2025-04-17

### Initial Release

This is the initial release of the Libr project, a decentralized messaging system with built-in content moderation.

### Added

#### Core Module

- Cryptographic utilities for message signing and verification
- Key pair generation and management system
- Environment configuration utilities
- Blockchain transaction model and interaction utilities

#### Node Module

- Database node implementation with PostgreSQL backend
- Message storage and retrieval API
- Timestamp-based message bucketing
- Docker and docker-compose configuration for easy deployment
- Worker pool for asynchronous message processing
- Node health check endpoints

#### Moderator Module

- Content moderation service based on Google Cloud NLP
- Configurable content filtering with category-based weights
- Message signing for verification of moderation approval
- Support for distributed moderation with multiple moderators

#### Client Module

- Message sending with automated moderator routing
- Timestamp-based message retrieval
- Database node selection algorithm based on blockchain state
- Support for k-of-n database node fault tolerance
- Majority-vote message resolution from multiple database nodes

#### Testing

- Test server for moderator and database API endpoints
- Express-based test implementation

### Technical Details

- Elliptic curve cryptography (secp256k1) for secure message signing
- SHA-256 hash-based verification
- Timestamp bucketing for message grouping
- JSON-based message certificates with multi-party signatures
- Fault-tolerant design with configurable replication factors

## Note

This is an early release. More features and improvements will be added in future versions.
