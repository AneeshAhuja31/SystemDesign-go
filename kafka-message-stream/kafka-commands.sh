#configure topics w partitions
docker exec -it kafka kafka-topics.sh --create --topic orders --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092
#verify
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --describe --topic orders --bootstrap-server localhost:9092
#inc partitions
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --alter --topic orders --partitions 6 --bootstrap-server localhost:9092
#del topic
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh --delete --topic orders --bootstrap-server localhost:9092