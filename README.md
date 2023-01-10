# vlsm
---

### IP summarization

    > cat test/file
    1.2.3.4
    1.2.3.5
    1.2.3.8
    10.1.1.1
    10.1.0.1

    > vlsm summarize -mask 23 test/file
       3  -  1.2.2.0/23
       2  -  10.1.0.0/23

### CIDR to range

    > vlsm to_range 10.1.1.1/24 1.1.1.1/23

     Submitted : 10.1.1.1/24
         Block : 10.1.1.0/24
       Netmask : 255.255.255.0
    Compliment : 0.0.0.255
      First IP : 10.1.1.0
       Last IP : 10.1.1.255
         Count : 256

     Submitted : 1.1.1.1/23
         Block : 1.1.0.0/23
       Netmask : 255.255.254.0
    Compliment : 0.0.1.255
      First IP : 1.1.0.0
       Last IP : 1.1.1.255
         Count : 512


