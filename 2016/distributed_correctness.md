title: "Distributed Something"
author:
  name: Jonathan Harlap
  twitter: jharlap
  github: jharlap
  email: "jonathan.harlap@wattpad.com"
theme: sjaakvandenberg/cleaver-light

--

# Something <br> Something <br> Distributed Something

--

### Distributed Programs

* execute concurrently
* have no shared memory
* and no shared clock

--

### CAP Theorem

* Consistency: All nodes see the same data at the same time
* Availability: Every request receives a response indicating success or failure
* Partition tolerance: The system operates normally even when network failures cause system partitioning

Eric Brewer stated that it is *impossible* for a *distributed* computer system to simultaneously guarantee all three

CAP prohibits only a tiny part of the design space: perfect availability and consistency in the presence of partitions, which are rare.

[_CAP Twelve Years Later: How the "Rules" Have Changed_](http://www.infoq.com/articles/cap-twelve-years-later-how-the-rules-have-changed)

--

### Consensus Problem

<div class="left">
<p>Multiple nodes achieve consensus if they agree on a value.</p>
<ul>
 <li> CP via Paxos or Raft with strong consistency guarantee is slow, high traffic </li>
 <li> AP via gossip - each node broadcasts updates to a random set of neighbours </li>
 <ul>
  <li> n = N: O(n²) traffic, O(1) propagation delay </li>
  <li> n ⊂ N: O(log n) traffic, O(log n) delay on average </li>
 </ul>
</ul>
</div>

<div class="right">
<img src="causality/CAP.svg" />
</div>

--

### Stateless vs Stateful

Stateless services

* aggregate/process responses from other services
* maintain **zero** local state
* scale horizontally effortlessly

Stateful services

* generate responses by processing state from a persistent store
* scales vertically easily, but without addressing availability
* scales horizontally by sharding for R/W loads

--

### Challenge

<img src="causality/overview.svg" height="600">

Which components are stateless? Stateful?

--

### One Service

Text cache:

* Every node maintains peer list
* Hash(key) → owner (node)
* Concurrent cache misses → single cache fill

AP system: When peers are unreachable, local cache miss → local fill

--

### Subsystem

PHP → Text Cache → Text Service → S3

* Updates always accepted
* Reads may return 404 until available (eventually consistent, LWW. one-write-wins per revision)
* PHP Memcache version ID can be stale
* MySQL read replica lag

--

### Interlude: Data Hazards

Multi-threaded CPUs must detect and resolve data hazards

* Read After Write (*true dependency*)
    1. **R2** ← R1 + R3
    2. R4 ← **R2** + R3
* Write After Read (*anti-depdendency*)
    1. R4 ← R1 + **R5**
    2. **R5** ← R1 + R2
* Write After Write (*output dependency*)
    1. **R2** ← R4 + R7
    2. **R2** ← R1 + R3

--

### Global

Mobile Client (*offline + local cache*) → PHP (*APC + memcache + MySQL replica lag*) → Text Cache → Text Service → S3

<img src="causality/text-overview.svg" height="500" />

--

### Potential Bugs

Delayed service-comment message consumer → username/avatar in comments lags behind story header

Home cached story title stale → Story title on home lags behind story info page

Mobile cached story info stale → Story title on info page lags behind home

...

--

### Thoughts

Delays, retries, out of order messages, concurrent updates, network partitions, stale caches... These all contribute to inconsistencies.

Users experience inconsistencies as **BUGS**, losing faith 😱

Abstractions at each level simplify design, but abstractions always leak 😱

CRDTs, Operational Transforms, ACID 2.0, CALM: Strategies to balance availability and consistency, minimize user pain.


--

# Conflict-free Replicated Data Types

--

### Conflict-free Replicated Data Types

State-based CRDTs (a.k.a. convergent RDTs, CvRDTs) transfer **complete state**
* Defines 3 functions
  * Update modifies local state
  * Query reads a value
  * Merge accepts remote updates
* High bandwidth usage but minimal network layer guarantees (cannot lose messages)

Merge functions must be:
* Associative: *m(A, B) = m(B, A)*
* Commutative: *m(m(A, B), C) = m(A, m(B, C))*
* Idempotent: *m(A, A) = A*

Merge function example: MAX

--

Operation-based CRDTs (a.k.a. commutative RDTs, CmRDTs) transfer **update operations**
* Merge function only needs to be commutative
* Requires network layer to guarantee exactly once delivery
* ∴ less bandwidth, more complex networking
* Merge function example: SUM

Proven possible to convert any CRDT design into either type

--
### Some Known CRDTs

* LWW Register
* G-counter
* PN-counter
* G-set (a.k.a. 1P-set)
* 2P-set
* OR-set
* ... more sets, graphs, sequences

--

### LWW Register

* Register is tuple (value, timestamp)
* Update: *r' = (new value, next(r.timestamp))*
* Merge: *m(A, B) = if A.timestamp > B.timestamp then A else B*

<img src="causality/crdt-register-lww.svg" height="500" />

--

# Who was last?

--

### Time
"Time is a component quantity of various measurements used to sequence events, to compare the duration of events or the intervals between them, and to quantify rates of change of quantities in material reality or in the conscious experience."
[Wikipedia: Sequence of Events]( https://en.m.wikipedia.org/wiki/Sequence_of_events )

Three properties of time we care about:

* Ordering on events
* Measurable
* Correlation with external events

Distributed simulations of physical clocks do *NOT* guarantee these properties for events originating on more than one node.

[Falsehoods programmers believe about time]( http://infiniteundo.com/post/25326999628/falsehoods-programmers-believe-about-time )

--

### Lamport Timestamp

| Sending a message:                 |   Receiving a message:                |
|------------------------------------|---------------------------------------|
|   time = time + 1                  |                                       |
|   timestamp = time                 |     (message, timestamp) = receive()  |
|   send(Message, timestamp)         |     time = max(timestamp, time) + 1   |

<img src="causality/lamport-clock.svg" height="400" />

Partial ordering of events. Concurrent events can have overlapping timestamps.

--

### Vector Clocks

<img src="causality/vector-clock.svg" height="400" />

* Partial ordering of events
* Dynamic set of actors → vector grows
* Vector pruning introduces errors

--

### Others

* [How to Have your Causality and Wall Clocks, too]( https://youtu.be/YqNGbvFHoKM )
* [Interval Tree Clocks: A Logical Clock for Dynamic Systems]( http://gsd.di.uminho.pt/members/cbm/ps/itc2008.pdf )
* [Time, clocks, and the ordering of events in a distributed system]( http://research.microsoft.com/en-us/um/people/lamport/pubs/time-clocks.pdf )
* [Dotted Version Vectors: Efficient Causality Tracking for Distributed Key-Value Stores](http://gsd.di.uminho.pt/members/vff/dotted-version-vectors-2012.pdf)

--

# Questions

--

### Resources

* [Distributed Systems for fun and profit]( http://book.mixu.net/distsys/ )
* Conflict-free Replicated Data Types
  * [Strong Eventual Consistency and CRDTs (video + slides)](http://research.microsoft.com/apps/video/default.aspx?id=153540)
  * [A comprehensive study of CRDTs (paper)](https://hal.inria.fr/inria-00555588)
* Databases
  * [Jensen: Partition tolerance analysis of DBs](https://aphyr.com/tags/jepsen)
  * [The Declarative Imperative](http://www.eecs.berkeley.edu/Pubs/TechRpts/2010/EECS-2010-90.pdf)
  * [Probabilistically Bounded Staleness for Practical Partial Quorums](http://pbs.cs.berkeley.edu/pbs-vldb2012.pdf)
* Time
  * [Time, clocks, and the ordering of events in a distributed system]( http://research.microsoft.com/en-us/um/people/lamport/pubs/time-clocks.pdf )
  * [Dotted Version Vectors: Efficient Causality Tracking for Distributed Key-Value Stores](http://gsd.di.uminho.pt/members/vff/dotted-version-vectors-2012.pdf)
  * [How to Have your Causality and Wall Clocks, too]( https://youtu.be/YqNGbvFHoKM )
  * [Interval Tree Clocks: A Logical Clock for Dynamic Systems]( http://gsd.di.uminho.pt/members/cbm/ps/itc2008.pdf )

