# Dust Tail

Cutting edge Mesos resource usage grabber written in GO.
Can grab usage statistics from remote nodes and store then locally.

## Building Dust Tail

Tested on GO 1.4.2 (min probably 1.3)

```
go get github.com/skonefal/dust-tail
go build
./dust-tail
```

## Mesos Compatibility

Serenity uses raw ResourceUsage message from resource monitor. I added an endopoint /monitor/serenity_statistics.json to Mesos.

I posted this on my branch: ```git clone -b serenity_resource_usage_for_dust_tail https://github.com/skonefal/mesos```

## Usage

Until external configuration will be completed, you need to hardcode your needs.

```
const (
	SAMPLING_TIME             = 1 * time.Second        // interval between sampling
	HTTP_TIMEOUT              = 200 * time.Millisecond // endpoint timeout
	EXPERIMENT_TIME           = 60 * time.Second       //
	EXPERIMENT_RESULTS_FOLDER = "results"              // folder with results

	STATISTICS_ENDPOINT = "/monitor/serenity_statistics.json" // mesos worker statistics endpoint
)

var mesosAgents = [...]string{
	"http://localhost:5051"} //list of mesos workers that will be sampled

```

Results will be stored in _EXPERIMENT_RESULTS_FOLDER_ in form _endpoint-node-name_timestamp_ eg.

```
localhost_2015-06-03 00:37:18.755666406 +0200 CEST
```

Results are stored as JSON array of arrays of executor resource usage.

Example result: http://pastebin.com/YDhu9LZa
