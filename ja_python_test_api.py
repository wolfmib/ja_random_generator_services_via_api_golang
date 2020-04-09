import requests
import json


response = requests.get("http://localhost:12345/get_random/1/7")
print("Response the random number is (1,7) =   ", response.json())

my_json = response.json()
print("Value is ",my_json['value'])