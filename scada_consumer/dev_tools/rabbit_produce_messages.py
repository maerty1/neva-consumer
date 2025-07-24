import json
import os
import pika


def send_message(channel, queue_name: str, data, service_name: str):
    if isinstance(data, dict):
        serialized_data = json.dumps(data).encode("utf-8")
    elif isinstance(data, list):
        serialized_data = json.dumps(data).encode("utf-8")
    else:
        serialized_data = data.SerializeToString()

    channel.basic_publish(exchange="", routing_key=queue_name, body=serialized_data)

    print(f"Сообщение из {service_name} отправлено в очередь {queue_name}")


def main():
    # Setting up the RabbitMQ connection and channel
    # rabbitmq_url = "amqp://admin:pYEDBqnMLoWWE@89.169.162.110:5672/"
    rabbitmq_url = "amqp://admin:pYEDBqnMLoWWE@localhost:5672/"
    connection_params = pika.URLParameters(rabbitmq_url)
    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()

    queue_name = "ASAP"
    service_name = "scada_producer"
    data = [
        {
            "value": "100",
            "dataType": "DINT",
            "lastChanged": "2024-04-06T07:05:15.5548405Z",
            "statusCodes": 12,
            "nodeId": "f80f7085-1dbf-4f49-9e63-a279a4a3227e",
            "nodeName": "RabbitMQ super",
            "ownerId": 1,
            "hash": "V2_CTP6_5",
            # "hash": "V1_CTP6_DAMN",
            "dataPointClassEnum": "Input",
        }
    ]

    # Declare the queue in case it doesn't exist
    channel.queue_declare(queue=queue_name, durable=True)

    # Send the message
    send_message(channel, queue_name, data, service_name)

    # Close the connection
    connection.close()


if __name__ == "__main__":
    main()
