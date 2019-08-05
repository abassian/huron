# HURON

## BFT Consensus algorithm for Abassian Shuffle

We use an adaptation of the Hashgraph consensus algorithm, invented by Leemon
Baird. Hashgraph is best described in the
[white-paper](http://www.swirlds.com/downloads/SWIRLDS-TR-2016-01.pdf) and its
[accompanying document](http://www.swirlds.com/downloads/SWIRLDS-TR-2016-02.pdf).


## Design

While initially it was built for Abassian Shuffle, Huron is designed to integrate with applications written in any programming
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
