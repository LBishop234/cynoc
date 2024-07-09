# CYNoC - Cycle Accurate NoC Simulation

## Configuration

### CLI Flags

| Flag | Shorthand | Operation |
| :--- | :-------- | :-------- |
| `-config FILE` | `-c FILE` | Specify simulation *yaml* configuration file |
| `-topology FILE` | `-t FILE` | Specify topology *GraphML* configuration file |
| `-traffic FILE` | `-tr FILE` | Specify traffic flows *csv* configuration file |
| `-cycle_limit VAL` | `-cy VAL` | Override the number of simulation cycles specified in the configuration file |
| `-max_priority VAL` | `-mp VAL` | Override the maximum traffic flow priority value specified in the configuration file |
| `-buffer_size VAL` | `-bs VAL` | Override the buffer size specified in the configuration file |
| `-processing_delay VAL` | `-pd VAL` | Override the header flit router processing delay specified in the configuration file |
| `-link_bandwidth VAL` | `-lb VAL` | Override the link bandwidth specified in the configuration file |
| `-analysis` | `-a` | Enables calculation of maximum basic network latency [[1]](#1) and Shi & Burns worst case network latency [[2]](#2) analyses for the configured simulation case |
| `-no-console-output` | `-nco` | Disables terminal results output, does not affect logging messages |
| `-results-csv FILE` | `-csv FILE` | Specifies the *csv* file to write full results to, creates the file if it does not exist |
| `-log` | | Enables LOG level messages |
| `-debug` | | Enables DEBUG level messages |
| `-trace` | | Enables TRACE level messages |

### Simulation Configuration File

Simulation configuration is configured via *yaml* file, e.g. `config.yaml`:
``` yaml
# Maximum number of network cycles simulated
cycle_limit: 16000
# Maximum priority value a traffic flow may possess
max_priority: 4
# Total size of a buffer in flits (divided by max_priority to calculate virtual channel size)
buffer_size: 16
# Header processing delay experienced at each router.
processing_delay: 6
# Link bandwidth in flits per cycle
link_bandwidth: 8
```

### Topology Configuration File

Network topology is configured via a *GraphML* file, 
with each node specifying `(x,y)` position attributes.
Every edge definition will create a full duplex connection between specified nodes,
i.e. only `n1 - n2` needs be defined as opposed to `n1 -> n2` & `n2 -> n1`.

E.g. `topology.xml`:
``` xml
<?xml version="1.0" encoding="UTF-8"?>
<graphml xmlns="http://graphml.graphdrawing.org/xmlns"  
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd">
  <graph id="G" edgedefault="undirected">
    <node id="n1">
        <data key="x">0</data>
        <data key="y">0</data>
    </node>
    <node id="n2">
        <data key="x">1</data>
        <data key="y">0</data>
    </node>
    <node id="n3">
        <data key="x">2</data>
        <data key="y">0</data>
    </node>
    <edge id="e1" source="n1" target="n2"/>
    <edge id="e2" source="n2" target="n3"/>
  </graph>
</graphml>
```

### Traffic Flow Configuration File

Traffic flows are configured using a *.csv* file, e.g. `traffic.csv`:
``` csv
id,priority,period,deadline,jitter,packet_size,route
t1,1,50,100,0,4,"[n1,n2,n3]"
t2,1,60,100,0,4,"[n1,n2]"
t3,2,50,100,0,4,"[n2,n3]"
```
- `id`: the traffic flow's unique id.
- `priority`: the traffic flow's unique priority level in the range [1, max_priority], inherited by all created packets.
- `period`: the regular interval, in cycles, at which point the traffic flow creates a new packet.
- `deadline`: the maximum tolerated latency for packets created by the traffic flow, i.e. packets must arrive at their destination router within $x$ cycles of creation.
    - Requires $deadline \leq period + jitter$ [[3]](#3).
- `jitter`: the maximum jitter the traffic flow's packets may experience, i.e. how long, in cycles, after creation may a packet be released to the network for transmission.
    - E.g. for a traffic flow with period $p$ and jitter $j$, a packet created on cycle $np$ will be released $x$ cycles after the packet's creation where $np \leq x < np+j$.
- `packet_size`: the number of body flits produced by the packet.
- `route`: the fixed route the traffic flow's packets traverse across the network.

## Results

### Terminal Output

| T_i | No. pkts | No. > D_i | min | mean  | max | D_i | J^R_i + C_i | J^R_i + R_i |
| --- | -------- | --------- | --- | ----- | --- | --- | ----------- | ----------- |
| t1  | 6400     | 0         | 21  | 25.54 | 30  | 100 | 33          | 33          |
| t2  | 5334     | 0         | 27  | 28.01 | 29  | 50  | 33          | 56          |
| t3  | 4000     | 0         | 23  | 25.53 | 28  | 150 | 33          | 33          |
| t4  | 8000     | 0         | 27  | 29.00 | 31  | 100 | 34          | 34          |
| t5  | 5334     | 0         | 25  | 25.00 | 25  | 25  | 27          | 27          |

- `T_i`: the traffic flow's unique id.
- `No. pkts`: the number of packets created by the traffic flow.
- `No. > D_i`: the number of packets exceeded their deadline.
- `min`: minimum simulated packet latency, from creation to arrival at destination.
- `mean`: mean simulated packet latency, from creation to arrival at destination.
- `max`: maximum simulated packet latency, from creation to arrival at destination.
- `D_i`: the traffic flow's deadline.
- `J^R_i + C_i` *(requires analysis)*: the traffic flow's release jitter added to maximum basic network latency [[1]](#1), giving the maximum packet latency without interference.
- `J^R_i + R_i` *(requires analysis)*: the traffic flow's release jitter added to Shi & Burns worst case network latency [[2]](#2), giving the worst case packet latency according to Shi & Burns.

### CSV File Output

```csv
TF_ID,Direct_Interference_Count,Indirect_Interference_Count,Num_Packets_Routed,Num_Packets_Exceeded_Deadline,Min_Latency,Mean_Latency,Max_Latency,Deadline,Schedulable,Jitter,Jitter_Plus_Basic,Jitter_Plus_Shi_And_Burns,Shi_Burns_Schedulable
t1,0,0,6400,0,21,25.55,30,100,true,10,33,33,true
t2,1,0,5334,0,27,28.02,29,50,true,3,33,56,false
t3,0,0,4000,0,23,25.51,28,150,true,6,33,33,true
t4,0,0,8000,0,27,28.97,31,100,true,5,34,34,true
t5,0,0,5334,0,25,25.00,25,25,true,1,27,27,false
```

- `TF_ID`: the traffic flow's unique id.
- `Direct_Interference_Count`: the number of traffic flows which impose direct interference as defined by Shi & Burns [[2]](#2).
- `Indirect_Interference_Count`: the number of traffic flows which impose indirect interference as defined by Shi & Burns [[2]](#2).
- `Num_Packets`: the number of packets created by the traffic flow.
- `Num_Packets_Exceeded_Deadline`: the number of packets exceeded their deadline.
- `Min_Latency`: minimum simulated packet latency, from creation to arrival at destination.
- `Mean_Latency`: mean simulated packet latency, from creation to arrival at destination.
- `Max_Latency`: maximum simulated packet latency, from creation to arrival at destination.
- `Deadline`: the traffic flow's deadline.
- `Schedulable`: the traffic flow's schedulability according to simulation results.
- `Jitter`: the traffic flow's release jitter.
- `Jitter_Plus_Basic` *(requires analysis)*: the traffic flow's release jitter added to maximum basic network latency [[1]](#1), giving the maximum packet latency without interference.
- `Jitter_Plus_Shi_Burns` *(requires analysis)*: the traffic flow's release jitter added to Shi & Burns worst case network latency [[2]](#2), giving the worst case packet latency according to Shi & Burns.
- `Shi_Burns_Schedulable` *(requires analysis)*: the traffic flow's schedulability according to Shi and Burns worst case network latency analysis [[2]](#2).

## References
- <a id='1'>[1]</a>
Duato, J., Yalamanchili, S., 1997. *Interconnection Networks: An Engineering Approach*. IEEE.
- <a id='2'>[2]</a>
Shi, Z. and Burns, A., 2008. Real-time communication analysis for on-chip networks with wormhole switching. In: *Second ACM/IEEE International Symposium on Networks-on-Chip (NoCs 2008)*, pp.161-170.
- <a id='3'>[3]</a>
Xiong, Q., Wu, F., Lu, Z. and Xie, C., 2017. Extending real-time analysis for wormhole NoCs. *IEEE Transactions on Computers*, **66**(9), pp.1532-1546.