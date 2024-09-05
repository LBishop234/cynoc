# CyNoC - Cycle Accurate NoC Simulation

A transmission level, cycle accurate, Network-on-Chip (NoC) simulator for NoCs implementing wormhole switching [[4]](#4), priority pre-emptive arbitration [[5]](#5), virtual channels [[3]](#3) & *Inq-n* [[2]](#2) routers.
In addition to simulating packet transmission and network latency, CyNoC also implements Shi & Burns' analysis model [[1]](#1).

## Configuration

CyNoC is invoked via the terminal and requires three configuration files to run. Examples can be found in `taskfile.yaml` and the `examples` directory.

E.g.: `./simulator -c example/basic/config.yaml -t example/basic/3-3-square.xml -tr example/basic/traffic.csv -a -log`

### CLI Flags

| Flag | Shorthand | Operation |
| :--- | :-------- | :-------- |
| `-config FILE` | `-c FILE` | Specify simulation characteristics configuration file (*yaml*) |
| `-topology FILE` | `-t FILE` | Specify topology configuration file  (*GraphML*) |
| `-traffic FILE` | `-tr FILE` | Specify traffic flows configuration file (*csv*) |
| `-cycle_limit VAL` | `-cy VAL` | Override the number of simulation cycles specified in the configuration file |
| `-max_priority VAL` | `-mp VAL` | Override the maximum traffic flow priority value specified in the configuration file |
| `-buffer_size VAL` | `-bs VAL` | Override the buffer size specified in the configuration file |
| `-processing_delay VAL` | `-pd VAL` | Override the header flit processing delay specified in the configuration file |
| `-analysis` | `-a` | Enables calculation Shi & Burns analysis model [[1]](#1) |
| `-no-console-output` | `-nco` | Disables results output to the terminal, does not affect logging messages |
| `-results-csv FILE` | `-csv FILE` | Specifies the *csv* filepath where simulator results will be written to |
| `-log` | | Enables $\geq$ LOG level messages |
| `-debug` | | Enables $\geq$ DEBUG level messages |
| `-trace` | | Enables $\geq$ TRACE level messages |

### Simulation Configuration File

Simulation & hardware characteristics are configured in a *yaml* file.

E.g. `config.yaml`:
``` yaml
# Number of network cycles simulated
cycle_limit: 16000
# Maximum priority value a traffic flow may possess (used to calculate virtual channel size)
max_priority: 4
# Total size of a buffer in flits (divided by max_priority to calculate virtual channel size)
buffer_size: 16
# Header flit processing delay experienced at each router, in network cycles.
processing_delay: 1
```

### Topology Configuration File

Network topology is configured by a [*GraphML*](http://graphml.graphdrawing.org/) file, 
Every edge definition will create a full duplex connection between specified nodes, i.e. only `n1 - n2` needs be defined as opposed to `n1 -> n2` & `n2 -> n1`.

E.g. `topology.xml`:
``` xml
<?xml version="1.0" encoding="UTF-8"?>
<graphml xmlns="http://graphml.graphdrawing.org/xmlns"  
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd">
  <graph id="G" edgedefault="undirected">
    <node id="n1"/>
    <node id="n2"/>
    <node id="n3"/>
    <edge id="e1" source="n1" target="n2"/>
    <edge id="e2" source="n2" target="n3"/>
  </graph>
</graphml>
```

### Traffic Flow Configuration File

Traffic flows are configured in a *.csv* file.

E.g. `traffic.csv`:
``` csv
id,priority,period,deadline,jitter,packet_size,route
t1,1,50,100,0,4,"[n1,n2,n2]"
t2,1,60,100,0,4,"[n1,n1]"
t3,2,50,100,0,4,"[n2,n2]"
```
- `id`: the traffic flow's unique id.
- `priority`: the traffic flow's unique priority level in the range [1, max_priority], inherited by all created packets.
- `period`: the regular interval, in network cycles, defining when the traffic flow creates a new packet.
- `deadline`: the maximum tolerated latency for packets created by the traffic flow, i.e. packets must arrive at their destination router within $x$ cycles of creation.
    - Requires $deadline \leq period + jitter$ [[2]](#2).
- `jitter`: the maximum jitter the traffic flow's packets may experience, i.e. how long, in cycles, after creation may a packet be released to the network for transmission.
    - E.g. for a traffic flow with period $p$ and jitter $j$, a packet created on cycle $np$ will be released $x$ cycles after the packet's creation where $np \leq x < np+j$.
- `packet_size`: the packet's size defining the number of flits it produces (including header and tail flits).
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
- `No. pkts`: the total number of packets created by the traffic flow.
- `No. > D_i`: the number of packets which exceeded their deadline.
- `min`: minimum simulated packet latency, from creation to arrival at destination.
- `mean`: mean simulated packet latency, from creation to arrival at destination.
- `max`: maximum simulated packet latency, from creation to arrival at destination.
- `D_i`: the traffic flow's packet deadline.
- `J^R_i + C_i` *(requires analysis)*: the traffic flow's release jitter added to maximum basic network latency, giving the maximum packet latency without interference.
- `J^R_i + R_i` *(requires analysis)*: the traffic flow's release jitter added to Shi & Burns worst case network latency [[1]](#1), giving the traffic flow's latency upper bound according to Shi & Burns.

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
- `Direct_Interference_Count`: the number of traffic flows which impose direct interference [[1]](#1) on this traffic flow.
- `Indirect_Interference_Count`: the number of traffic flows which impose indirect interference [[1]](#1) on this traffic flow.
- `Num_Packets`: the number of packets created by the traffic flow.
- `Num_Packets_Exceeded_Deadline`: the number of packets which exceeded their deadline.
- `Min_Latency`: minimum simulated packet latency, from creation to arrival at destination.
- `Mean_Latency`: mean simulated packet latency, from creation to arrival at destination.
- `Max_Latency`: maximum simulated packet latency, from creation to arrival at destination.
- `Deadline`: the traffic flow's packet deadline.
- `Schedulable`: the traffic flow's schedulability according to simulation results.
- `Jitter`: the traffic flow's release jitter.
- `Jitter_Plus_Basic` *(requires analysis)*: the traffic flow's release jitter added to maximum basic network latency, giving the maximum packet latency without interference.
- `Jitter_Plus_Shi_Burns` *(requires analysis)*: the traffic flow's release jitter added to Shi & Burns worst case network latency [[1]](#1), giving the traffic flow's latency upper bound according to Shi & Burns.
- `Shi_Burns_Schedulable` *(requires analysis)*: the traffic flow's schedulability according to Shi and Burns [[1]](#1).

# Important Notes

## Notes

Please be aware Shi & Burns analysis model is not correct and has been shown to produce optimistic latency upper bounds under specific routing combinations [[6]](#6).
More recent paper have provided fixes for this flaw and should be used as an alternative [[7]](#7).

## Usage & Acknowledgements

**TODO**

This project is available for use under the **INSERT** license.

We do however request that any academic publications which utilize this simulator cite the repository in their paper.

E.g. for BibLaTeX citations:
```bibtex
@online{CyNoC,
  author = {Leo Bishop},
  title = {The CyNoC Simulator},
  year = 2024,
  url = {https://github.com/LBishop234/cynoc},
  urldate = {2024-13-09}
}
```

## Contributions

This project is maintained in my spare time for my research on NoCs.
As such any contributions are very welcome though please be aware I will respond as best my schedule allows.

Contributions implementing more recent NoC analysis models would be particularly welcome.

## References
- <a id='1'>[1]</a>
Shi, Z. and Burns, A., 2008, April. Real-time communication analysis for on-chip networks with wormhole switching. In *Second ACM/IEEE International Symposium on Networks-on-Chip (nocs 2008)* (pp. 161-170). IEEE.
- <a id='2'>[2]</a>
Xiong, Q., Wu, F., Lu, Z. and Xie, C., 2017. Extending real-time analysis for wormhole NoCs. *IEEE Transactions on Computers*, 66(9), pp.1532-1546.
- <a id='3'>[3]</a>
Dally, W.J., 1992. Virtual-channel flow control. *IEEE Transactions on Parallel and Distributed systems*, 3(1), pp.194-205.
- <a id='4'>[4]</a>
Ni, L.M. and McKinley, P.K., 1993. A survey of wormhole routing techniques in direct networks. Computer, 26(1), pp.62-76.
- <a id='5'>[5]</a>
Song, H., Kwon, B. and Yoon, H., 1999. Throttle and preempt: a flow control policy for real-time traffic in wormhole networks. *Journal of systems architecture*, 45(8), pp.633-649.
- <a id='6'>[6]</a>
Xiong, H., Lu, Z., Wu, F., Xie, C., 2016. Real-time analysis for wormhole NoC: revisited and revised. *Proceeding of the 26th edition on Great Lakes Symposium on VLSI*.
- <a id='7'>[7]</a>
Indrusiak, L. S., Nikolic, B., Burns, A., 2016. Analysis of buffering effects on hard real-time priority-preemptive wormhole networks, *arXiv:1606.02942*.