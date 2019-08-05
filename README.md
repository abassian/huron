# HURON

## BFT Consensus platform for distributed applications

Huron enables multiple computers to behave as one. It uses Peer to Peer (P2P)
networking and a consensus algorithm to guarantee that a group of connected
computers process the same commands in the same order; a technique known as
state-machine replication. This makes for secure systems that can tolerate
arbitrary failures, including malicious behavior.

For guidance on how to install and use Huron we are releasing soon extensive documentation.

**NOTE**:
This is alpha software. Please contact us if you intend to run it in production.

## Consensus Algorithm and Blockchain

We use an adaptation of the Hashgraph consensus algorithm, invented by Leemon
Baird. Hashgraph is best described in the
[white-paper](http://www.swirlds.com/downloads/SWIRLDS-TR-2016-01.pdf) and its
[accompanying document](http://www.swirlds.com/downloads/SWIRLDS-TR-2016-02.pdf).


## Design

Huron is designed to integrate with applications written in any programming
language.

### Overview

```text
    +--------------------------------------+
    | APP                                  |
    |                                      |
    |  +-------------+     +------------+  |
    |  | Service     | <-- | State      |  |
    |  |             |     |            |  |
    |  +-------------+     +------------+  |
    |          |                ^          |
    |          |                |          |
    +----------|----------------|----------+
               |                |
--------- SubmitTx(tx) ---- CommitBlock(Block) ------- JSON-RPC/TCP or in-memory
               |                |
 +-------------|----------------|------------------------------+
 | HURON      |                |                              |
 |             v                |                              |
 |          +----------------------+                           |
 |          | App Proxy            |                           |
 |          |                      |                           |
 |          +----------------------+                           |
 |                     |                                       |
 |   +-------------------------------------+                   |
 |   | Core                                |                   |
 |   |                                     |                   |
 |   |  +------------+                     |    +----------+   |
 |   |  | Hashgraph  |       +---------+   |    | Service  |   |
 |   |  +------------+       | Store   |   | -- |          | <----> HTTP
 |   |  +------------+       +----------   |    |          |   |
 |   |  | Blockchain |                     |    +----------+   |
 |   |  +------------+                     |                   |
 |   |                                     |                   |
 |   +-------------------------------------+                   |
 |                     |                                       |
 |   +-------------------------------------+                   |
 |   | Transport                           |                   |
 |   |                                     |                   |
 |   +-------------------------------------+                   |
 |                     ^                                       |
 +---------------------|---------------------------------------+
                       |
                       v
                  P2P Network
```
## USAGE

At the end of the current sprint we are going to release documentation on usage and testing.
