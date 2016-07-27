# ec2FleetCompare
Small / Fast command line tool to determine the cheapest instance types based on fleet size, and minimum cpu, memory and network requirements. It provides filter options and outputs a list of instance types that fulfil your criteria based on price (both demand, RI and spot prices supported). 

ec2FleetCompare can be used to look at cost comparisons on individual instance costs but can also be used to find the cheapest options for a "fleet" of ec2 instances. As long as you know the total number of VCPS's or GB's or ram required across an entire fleet the tool will provide you the cheapest option to achieve this.

# Download

Pre compiled binaries are available for Mac, Linux and Windows. Please see download links below:

[Linux Download] (https://s3-us-west-2.amazonaws.com/andy-gen/ec2FleetCompare/linux/ec2FleetCompare)  - (md5 ba862dc75d474c9d6c24ba86382e5a5e)

[Mac Download] (https://s3-us-west-2.amazonaws.com/andy-gen/ec2FleetCompare/osx/ec2FleetCompare) - (md5 57b5ef13a0ab7263ecaf0656996fa155) 

[Windows Download] (https://s3-us-west-2.amazonaws.com/andy-gen/ec2FleetCompare/win/ec2FleetCompare.exe) - (md5 375657fec1b332c2343dc41668e4aa65)

# Usage

To view help / option information. All options have defaults, the various options over-ride them.

```
./ec2FleetCompare --help
```

# Examples

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

