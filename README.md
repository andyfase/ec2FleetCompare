# ec2FleetCompare
Small / Fast command line tool to determine the cheapest instance types based on fleet size, and minimum cpu, memory and network requirements. It provides filter options and outputs a list of instance types that fulfil your criteria based on price (both demand, RI and spot prices supported). 

ec2FleetCompare can be used to look at cost comparisons on individual instance costs but can also be used to find the cheapest options for a "fleet" of ec2 instances. As long as you know the total number of VCPS's or GB's or ram required across an entire fleet the tool will provide you the cheapest option to achieve this.

# Download

Pre compiled binaries are available for Mac, Linux and Windows. Please see download links below:

[Linux Download] (https://s3-us-west-2.amazonaws.com/andy-gen/ec2FleetCompare/linux/ec2FleetCompare)  - (md5 a5547c2fbd070a653f643c7adc8487b8)

[Mac Download] (https://s3-us-west-2.amazonaws.com/andy-gen/ec2FleetCompare/osx/ec2FleetCompare) - (md5 aede3be60ec614188a8e84c35a508e17) 

[Windows Download] (https://s3-us-west-2.amazonaws.com/andy-gen/ec2FleetCompare/win/ec2FleetCompare.exe) - (md5 eeeb31a2a537dfb89d6e42dbd81de2f2)

Note: you may need to make executable i.e ```chmod 500 ./ec2FleetCompare``` or similiar for windows.

# Usage

To view help / option information. All options have defaults, the various options over-ride them.

```
./ec2FleetCompare --help
```
# Output

The command line tool will output a ASCII based table describing the instance types (and number of them if using for a fleet) that fulfil your criteria. This is sorted by default by on-demand pricing but can also be sorted via spot pricing too

```
+--------+-------------+-------+-----------+----------+------------+---------+---------+--------------------+--------------------+------------+
| # INST |    TYPE     | VCPUS | VCPU FREQ | MEM/INST |  NETWORK   | IS TYPE | IS SIZE |   DEMAND $/HOUR    |    SPOT $/HOUR     | SPOT % SAV |
+--------+-------------+-------+-----------+----------+------------+---------+---------+--------------------+--------------------+------------+
|      1 | c4.8xlarge  |    36 | 2.9 GHz   |     60.0 | 10 Gigabit | N/A     | N/A     | $1.68 ($1.68 each) | $0.34 ($0.34 each) | 80%        |
|      1 | c3.8xlarge  |    32 | 2.8 GHz   |     60.0 | 10 Gigabit | SSD     | 640 GB  | $1.68 ($1.68 each) | $0.36 ($0.36 each) | 79%        |
|      1 | cc2.8xlarge |    32 | 2.6 GHz   |     60.5 | 10 Gigabit | HDD     | 3360 GB | $2.00 ($2.00 each) | $0.28 ($0.28 each) | 86%        |
|      1 | m4.10xlarge |    40 | 2.4 GHz   |    160.0 | 10 Gigabit | N/A     | N/A     | $2.39 ($2.39 each) | $0.45 ($0.45 each) | 81%        |
|      1 | g2.8xlarge  |    32 | 2.6 GHz   |     60.0 | 10 Gigabit | SSD     | 240 GB  | $2.60 ($2.60 each) | $2.80 ($2.80 each) | -8%        |
|      1 | r3.8xlarge  |    32 | 2.5 GHz   |    244.0 | 10 Gigabit | SSD     | 640 GB  | $2.66 ($2.66 each) | $0.60 ($0.60 each) | 78%        |
|      1 | cr1.8xlarge |    32 |           |    244.0 | 10 Gigabit | SSD     | 240 GB  | $3.50 ($3.50 each) | $0.57 ($0.57 each) | 84%        |
|      1 | d2.8xlarge  |    36 | 2.4 GHz   |    244.0 | 10 Gigabit | HDD     | 8000 GB | $5.52 ($5.52 each) | $0.76 ($0.76 each) | 86%        |
|      1 | 2.8xlarge  |    32 | 2.5 GHz   |    244.0 | 10 Gigabit | SSD     | 6400 GB | $6.82 ($6.82 each) | $0.75 ($0.75 each) | 89%        |
+--------+-------------+-------+-----------+----------+------------+---------+---------+--------------------+--------------------+------------+
```

# Example Usage

Find linux based ec2 instances with over 8GB ram and at least 4 VCPU's
```
./ec2FleetCompare -c 4 -m 8
```

Find Windows based ec2 instances that have at least 1TB of SSD instance store disk available.
```
./ec2FleetCompare -os win -dt ssd -d 1024
```

Find Windows based ec2 instances that have Gigabit network interfaces, sorted by Spot pricing
```
./ec2FleetCompare -os win -nw gbit -s spot
```

Find cheapest fleet of upto 1000 instance that has a total VCPU capacity of 10000, with each instance at least having 32 VCPUS. Total memory capacity of at least 24TB with all nodes having Gigabit networking. Sorted by spot pricing.
```
./ec2FleetCompare -n 1000 -fc 10000 -c 32 -fm 24576 -nw gbit -s spot
```

Find cheapest fleet of upto 500 i2 type instances with a total memory cpacity of 24TB with each node having at least 3.2TB of SSD instance store disk available. Sorted by spot pricing.
```
./ec2FleetCompare -n 500 -fm 24576 -dt SSD -d 3200 -i i2 -s spot
```

# Developing

This is written in [golang] (https://golang.org/). So you will need to download the GO compiler, set your ```GOPATH``` environment variable correctly and then install all the pre-req modules listed in the source file (```go get <package>```). 

