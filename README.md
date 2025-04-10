# norman

Realtime distributed OLAP datastore, designed to answer OLAP queries with low latency written in Go

## cURLs requests for testing

To create a new Schema

```bash
curl -X POST -H "Content-Type: application/json" -d @test/commander/create_schema_flights.json http://localhost:8080/commander/v1/tenants/default/schemas
```

## Kafka initialization

Create a topic for testing the ingestion job

```bash
$ docker exec -it kafka bash
[appuser@kafka]# cd /
[appuser@kafka]# cd bin
[appuser@kafka]# ./kafka-topics --create --topic test-events --bootstrap-server localhost:9092
[appuser@kafka]# ./kafka-topics --describe --topic test-events --bootstrap-server localhost:9092
```

Now, let's write some events into the topic

```bash
[appuser@kafka]# ./kafka-console-producer --topic test-events --bootstrap-server localhost:9092
```

Copy the following JSON lines

```JSON
{"studentID":205,"firstName":"Natalie","lastName":"Jones","gender":"Female","subject":"Maths","score":3.8,"timestampInEpoch":1571900400000}
{"studentID":205,"firstName":"Natalie","lastName":"Jones","gender":"Female","subject":"History","score":3.5,"timestampInEpoch":1571900400000}
{"studentID":207,"firstName":"Bob","lastName":"Lewis","gender":"Male","subject":"Maths","score":3.2,"timestampInEpoch":1571900400000}
{"studentID":207,"firstName":"Bob","lastName":"Lewis","gender":"Male","subject":"Chemistry","score":3.6,"timestampInEpoch":1572418800000}
{"studentID":209,"firstName":"Jane","lastName":"Doe","gender":"Female","subject":"Geography","score":3.8,"timestampInEpoch":1572505200000}
{"studentID":209,"firstName":"Jane","lastName":"Doe","gender":"Female","subject":"English","score":3.5,"timestampInEpoch":1572505200000}
{"studentID":209,"firstName":"Jane","lastName":"Doe","gender":"Female","subject":"Maths","score":3.2,"timestampInEpoch":1572678000000}
{"studentID":209,"firstName":"Jane","lastName":"Doe","gender":"Female","subject":"Physics","score":3.6,"timestampInEpoch":1572678000000}
{"studentID":211,"firstName":"John","lastName":"Doe","gender":"Male","subject":"Maths","score":3.8,"timestampInEpoch":1572678000000}
{"studentID":211,"firstName":"John","lastName":"Doe","gender":"Male","subject":"English","score":3.5,"timestampInEpoch":1572678000000}
{"studentID":211,"firstName":"John","lastName":"Doe","gender":"Male","subject":"History","score":3.2,"timestampInEpoch":1572854400000}
{"studentID":212,"firstName":"Nick","lastName":"Young","gender":"Male","subject":"History","score":3.6,"timestampInEpoch":1572854400000}
```

Read the events for testing

```bash
[appuser@kafka]# ./kafka-console-consumer --topic test-events --from-beginning --bootstrap-server localhost:9092
```

If you see the events, you are done!

### Ingest few thousands lines

Copy the file to the kafka running container

```bash
docker cp test/data/transcript.jsonl kafka:/home/appuser/transcript.json
```

Execute the ingestion job to create the Kafka consumer and prepare Norman for ingesting events (see the steps below) and then push these messages to the kafka topic with the following command

```bash
[appuser@kafka]# ./kafka-console-producer --topic test-events --bootstrap-server localhost:9092 < /home/appuser/transcript.json
```

### Run the ingestion job

Create the schema in Norman

```bash
curl -X POST -H "Content-Type: application/json" -d @test/commander/create_schema_transcript.json http://localhost:8080/commander/v1/tenants/default/tables
```

Create the ingestion job

```bash
curl -X POST -H "Content-Type: application/json" -d @test/commander/create_job_transcript.json http://localhost:8080/commander/v1/tenants/default/jobs
```
