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
| `-buffer_size VAL` | `-bs VAL` | Override the buffer size (b) specified in the configuration file |
| `-flit_size VAL` | `-fs VAL` | Override the flit size (b) specified in the configuration file |
| `-processing_delay VAL` | `-pd VAL` | Override the header flit router processing delay (cycles) specified in the configuration file |
| `-analysis` | `-a` | Enables calculation of Shi & Burns worst case network latency analysis for the configured simulation case |
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
# Total size of a buffer in bytes
buffer_size: 16
# Flit size in bytes
flit_size: 4
# Header processing delay experienced at each router.
processing_delay: 6
# Link bandwidth in bytes
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
- `deadline`: the maximum tolerated network latency for packets created by the traffic flow, i.e. packets must arrive at their destination router within $x$ cycles of creation.
- `jitter`: the maximum jitter the traffic flow's packets may experience, i.e. how long, in cycles, after creation may a packet be released to the network for transmission.
    - E.g. for a traffic flow with period $p$ and jitter $j$, a packet created on cycle $np$ will be released $x$ cycles after the packet's creation where $np \leq x < np+j$.
- `packet_size`: the size of a packet's payload in bytes.
- `route`: the fixed route the traffic flow's packets traverse across the network.