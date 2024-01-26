# Counter service

Tiny service to fetch aggregated statistics from participants 

The goal of this service is to propose a standardized way to fetch active participants counter to enable counter federation (like legacy Influenzanet website).
This service handles counter for **one instance** but several studies can be fetched using the same service instance.

> [!CAUTION]
> This service is not ready to be used in production. It's not stable yet. Any part of this service is susceptible to change without compatibility with previous It should be released in start of year 2024. 

## Influenzanet standard installation

To enable simple discovery of platforms, this counter service is expected to be available on the same path of the influenzanet platform web domain.

The URL must be : https://[platform-domain]/.well-known/influenzanet/counter

For example, if the platform base domain  is example.influenzanet.com then the expected URL should be:
    https://example.influenzanet.com/.well-known/influenzanet/counter


## Handled metrics (for each study)

- participants_enrolled (*count*): Count of participants of the study with active status
- participants_active (*count*): Count of participants of the study with at least one survey submitted in the active survey list (if the delay of the last submission is less than the active_delay)

The type of each the counter is provided between parenthesis see below for types.

Metrics are not updated on each request, but periodically (defined by `update_delay` duration of each counter). Values exposed by the server are the values
evaluated at the update.

## Usage

The server is loaded by default on :5021 port

Endpoints :

### $baseURI/
Returns the default counters and reference to other public counters

```
{
    "platform": "",
    "extra": ["other_counter"],
    "counters": {
        "influenzanet": {
            "period": 86400,
            "data": [
                {
                    "name": "participants_active",
                    "value": 5891,
                    "type": "count",
                    "time": 1703783790
                },
                {
                    "name": "participants_enrolled",
                    "value": 11585,
                    "type": "count",
                    "time": 1703783790
                }
            ]
        }
    }
}
```

- `platform` : code of the platform (usually country code)
- `extra`: list of other available counter names (if public=true in counter definition)
- `counters`: list of root counters (structure of each entry is the same as the individual counter, see below)

### $baseURI/counter/`$name`

Fetch stats for a specific counter with name provided in $name (e.g. 'influenzanet').

Response a counter stat object:

- `period`: the update period of the counter (time between 2 updates) in seconds
- `data`: array of metrics stat

Each metric has 4 possible fields: 

- name: metric name
- value: number if type='count' or an object if type='map'
- type: 'count' or 'map'
- time: timestamp field when the counter has been evaluated for the last time (Unix timestamp)

Example

```json
{
    "period": 86400,
    "data": [
        {
            "name": "participants_enrolled",
            "value": 11585,
            "type": "count",
            "time": 1703783790
        },
        {
            "name": "participants_active",
            "value": 5891,
            "type": "count",
            "time": 1703783790
        }
    ]
}
```

### $baseURI/meta.json (optional)

Show configuration of counters after the config is loaded. It's intended to be used only for debug and disabled by default.

## Configuration

Configuration is done using environment variables:

- Db Connection vars to Study Db: see in [Study Service](https://github.com/grippenet/study-service/blob/master/build/docker/example/study-service-env.list), accept User Db variables and general db client settings.

- `INSTANCE_ID` instance id of the target instance of the platform
- `PLATFORM` : Code of the influenzanet Platform (Optional, attributed by Influenzanet)
- `INFLUENZANET_COUNTER`: Specification for Influenzanet counter. It can be either the study name of the Influenzanet Surveillance compliant study or a JSON counter definition, if json object is incomplete, default value is used for the missing keys (so you can provide only values where you want to override the default).
- `EXTRA_COUNTERS` : Optional, A Json array with extra counter definitions (each entry should be a Counter Definition json object )
- `EXTRA_COUNTERS_FILE` : Same as  `EXTRA_COUNTERS` but with the name of the file containing the json array

The InfluenzaNet default counter uses this configuration:
```json
{
    "active_surveys":["intake", "weekly", "vaccination"],
    "root": true,
    "public": true,
    "name": "influenzanet",
    "active_delay": "13104h",
    "update_delay": "24h",
}
```

Study key is expected (no default is provided). It's not advised to override the counter name as it should be a common value.


### Counter Definition

 a counter can be defined by a JSON 

 ```json
{
    "studykey":"my_study",
    "active_surveys":["my_survey"],
    "root": false,
    "public": true,
    "name": "my_counter",
    "active_delay": "720h",
    "update_delay": "24h",
}

```

- studykey: name of the study (in study service)
- active_surveys: list of survey key for which a submitted response will count its participant as active
- active_delay: Maximum delay (from current time) of survey submission to consider the participant as active (integer or time.Duration format, e.g. '10s','30m','1h')
- update_delay: Delay between two update of the counter
- name: Name of the counter (must be unique)
- root: boolean, if true the counter is shown at the root of the service address
- public: boolean, if true the counter is shown as extra in the root (if root=false), ignored if root=true 

### HTTP Server
To configure Http Server:

- `PORT`: change the listening port
- `GIN_MODE`: can be configured ('release' will be less verbose)
- `META_AUTH_KEY`: if provided to a non empty value, activate the `meta.json` endpoint 
