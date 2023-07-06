# Counter service

Tiny service to fetch aggregated statistics from participants 

The goal of this service is to propose a standardized way to fetch active participants counter to enable counter federation (like legacy Influenzanet website).
This service handles counter for **one instance** but several studies can be fetched using the same service instance.

## Handled metrics (for each study)

- participants_active (*count*): Count of participants of the study with active status
- participants_intake (*count*): Count of participants of the study with intake survey submitted since the `FROM_DATE` env

The type of each the counter is provided between parenthesis see below for types


## Usage

The server is loaded by default on :5021 port

Endpoints :

### $baseURI/

Status page, just to say hello.

### $baseURI/whoami

Fetch the service meta data

```json
{
studies: [
    "grippenet"
],
influenzanet: "",
from: "2022-11-22T00:00:00Z"
}
```

`studies`: list of studies for which a counter is available
`influenzanet`: name of the influenzanet compliant study if it's not 'influenzanet'
`from`: `FROM_DATE` value

### $baseURI/study/`$study`

Fetch stats for the provided study key (replace `$study` by the study key name)

Response: an array of Counter results

Each Counter has 4 possible fields: 

- name: counter name
- value: object value (depend on type see below)
- type: 'count' or 'map'
- time: time field when the counter has been evaluated for the last time

value field is 
- a number if type is 'count'
- an object (key value pair) if type is 'map'

Example

```json
[
    {
        "name": "simple_counter",
        "type": "count",
        "value": 9911
    },
    {
        "name": "map_counter",
        "type": "map",
        "value": {
                "1": 292,
                "2": 861,
                "3": 934,
                "4": 938,
                "5": 928,
                "6": 620
        }
    }
]
```

## Configuration

Environments:

- Db Connection vars: see in [Study Service](https://github.com/grippenet/study-service/blob/master/build/docker/example/study-service-env.list), accept User Db variables and general db client settings

- `STUDIES`: comma separated list of studies
- `INSTANCE_ID` instance id of the target instance of the platform
- `FROM_DATE`: Date (ISO format YYYY-MM-DD) from which count intake submissions
- `UPDATE_DELAY`: Delay in Minutes to update the counter internally
- `INFLUENZANET_STUDY`: Name of the Influenzanet compliant study 

To configure Http Server:

- `PORT`: change the listening port
- `GIN_MODE`: can be configured ('release' will be less verbose)
