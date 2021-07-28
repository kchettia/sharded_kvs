# Use this script to add N key value pairs to the network
# Installation happens automatically
import requests
import random
import string

letters = list(string.ascii_lowercase)


def generate_keys(num_keys: int = 10000):
    keys = []
    for i in range(num_keys):
        keys.append(
            "{}-{}-{}".format(
                random.choice(letters), random.choice(letters), random.choice(letters)
            )
        )
    return keys


def generate_official_keys(num_keys: int = 10000):
    keys = []
    for i in range(num_keys):
        keys.append("test_2_%d" % i)
    return keys


URI = "http://0.0.0.0:13802/kvs/keys"


def send_key_value(key, value):
    # key, value = get_random_string(), get_random_string()
    url = URI + f"/{key}"
    headers = {"Content-Type": "application/json"}
    requests.put(url, headers=headers, json={"value": value})


for index, key in enumerate(generate_official_keys()):
    print(key,"\n")
    send_key_value(key, index)
