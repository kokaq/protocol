# Wire Protocol

## Common Header

Following is the 8-byte header protocol. Depending on the type, the appropriate message payload follows the fixed header section:

| offset | name                | size (bytes) |         | meaning                                                                         |
|--------|---------------------|--------------|---------|---------------------------------------------------------------------------------|
| 0      | Magic               | 2            |         | Magic number, used to identify kokaq message. '0x420'                           |
| 2      | Version             | 1            |         | Protocol version, current version is 1.                                         |
| 3      | Message Type flag   | 1            | bit 0-5 | Message Type, 0: Operational Message, 1: Admin Message, 2: Control Message      |
|        | RQ                  |              | bit 6-7 | RQ flags, 0: response, 1: two way request, 3: one way request                   |
| 4      | Opaque              | 4            |         | The Opaque data set in the request will be copied back in the response          |

```bash
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       |msg type   |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+

  messge type:
    0x01    Operational
    0x02    Admin
    0x03    Control

  request type:
    0x01    Response
    0x02    TwoWay
    0x03    OneWay

```

## Operation Header

### Operation Request Header

| offset | header name         | size (bytes) |
|--------|---------------------|--------------|
| 0      | OpCode              | 1            |
| 1      | Client              | 1            |
| 2      | Opaque              | 1            |
| 3      | Tag/ID              | 1            |

```bash
operational request header
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | opcode        | clientId      | opaque        | Tag/ID        |
------+---------------+---------------+---------------+---------------+

  opcode:
    0x00    Nop
    0x01    Create
    0x02    Delete
    0x03    Get
    0x04    Peek
    0x05    Pop
    0x07    Push
    0x08    AcquirePeekLock 
    0x05    ReleasePeekLock 

  clientId:
    0x00    ProxyHttp
    0x01    ProxyAmqp
    0x02    QueueService
    0x03    StorageService
    0x04    HealthService
```

### Operation Response Header

| offset | header name         | size (bytes) |
|--------|---------------------|--------------|
| 0      | OpCode              | 1            |
| 1      | Status+Reason       | 1            |
| 2      | Opaque              | 1            |
| 3      | Tag/ID              | 1            |

```bash
operational response header
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+-------+-------+---------------+---------------+
   12 | opcode        | status|reason | opaque        | Tag/ID        |
------+---------------+-------+-------+---------------+---------------+

  opcode:
    0x00    Nop
    0x01    Create
    0x02    Delete
    0x03    Get
    0x04    Peek
    0x05    Pop
    0x07    Push
    0x08    AcquirePeekLock 
    0x05    ReleasePeekLock 

  status:
    0x00    Success
    0x01    Fail
    0x02    PartialSucess
    0x03    Unknown

  reason:
    0x01    Ok
    0x02    Bad
    0x03    Exists
    0x03    NotAllowed
    0x03    Infra

  tag:
    0x01    Metadata
    0x02    Payload
```

## Message Body

### Payload Component

```bash

Payload
------+---------------+---------------+---------------+---------------+
    0 | Tag/ID (0x01) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
    4 | namespace len | queue length  | payload len                   |
------+---------------+---------------+---------------+---------------+
    8 | nsname?                                                       |
------+---------------+---------------+---------------+---------------+
   12 | queuename?                                                    |
------+---------------+---------------+---------------+---------------+
   16 | payload?                                                      |
------+---------------+---------------+---------------+---------------+
------+---------------+---------------+---------------+---------------+


  tag:
    0x01    Metadata
    0x02    Payload
```

### MetaData Component

```bash
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
      |              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | Tag/ID (0x02) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
    4 | fieldcount len                                                |
------+---------------+---------------+---------------+---------------+
    8 | fieldtag      |  size         |               |               |
------+---------------+---------------+---------------+---------------+
------+---------------+---------------+---------------+---------------+



  tag:
    0x01    Metadata
    0x02    Payload
```

#### Pre-defined metadata

| Metadata               | Tag  |
|------------------------|------|
| TimeToLive             | 0x01 |
| Version                | 0x02 |
| Creation Time          | 0x03 |
| Expiration Time        | 0x04 |
| UUID                   | 0x05 |
| Source Info            | 0x06 |
| Last Modification time | 0x07 |
| Originator RequestID   | 0x08 |
| Correlation ID         | 0x09 |
| Request Handling Time  | 0x0a |

## Request

### Simple Request

```bash

      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       | msg type  |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
    8 | opcode        | clientId      | opaque                        |
------+---------------+---------------+---------------+---------------+
```

### Request With Payload

```bash

      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       | msg type  |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
    8 | opcode        | clientId      | opaque        | Tag/ID        |
------+---------------+---------------+---------------+---------------+
   12 | Tag/ID (0x01) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
   16 | namespace len | queue length  | payload len                   |
------+---------------+---------------+---------------+---------------+
   20 | nsname?                                                       | 
------+---------------+---------------+---------------+---------------+
   24 | queuename?                                                    |
------+---------------+---------------+---------------+---------------+
   28 | payload?                                                      |
------+---------------+---------------+---------------+---------------+
------+---------------+---------------+---------------+---------------+

```

### Request With Metadata

```bash
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       | msg type  |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
    8 | opcode        | clientId      | opaque        | Tag/ID        |
------+---------------+---------------+---------------+---------------+
   12 | Tag/ID (0x01) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
   16 | fieldcount len                                                |
------+---------------+---------------+---------------+---------------+
   20 | fieldtag      |  size         |               |               |
------+---------------+---------------+---------------+---------------+
------+---------------+---------------+---------------+---------------+

```

## Response

### Simple Response

```bash
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       | msg type  |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
    8 | opcode        | status|reason | opaque        | Tag/ID        |
------+---------------+---------------+---------------+---------------+
```

### Response With Payload

```bash
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       |msg type   |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
    8 | opcode        | status|reason | opaque        | Tag/ID        |
------+---------------+---------------+---------------+---------------+
   12 | Tag/ID (0x02) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
   16 | namespace len | queue length  | payload len                   |
------+---------------+---------------+---------------+---------------+
   20 | nsname?                                                       | 
------+---------------+---------------+---------------+---------------+
   24 | queuename?                                                    |
------+---------------+---------------+---------------+---------------+
   28 | payload?                                                      |
------+---------------+---------------+---------------+---------------+
------+---------------+---------------+---------------+---------------+

```

### Response With Metadata

```bash
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       |msg type   |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
    8 | opcode        | status|reason | opaque        | Tag/ID        |
------+---------------+---------------+---------------+---------------+
   12 | Tag/ID (0x02) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
   16 | fieldcount len                                                |
------+---------------+---------------+---------------+---------------+
   20 | fieldtag      |  size         |               |               |
------+---------------+---------------+---------------+---------------+
------+---------------+---------------+---------------+---------------+

```
