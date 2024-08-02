import requests
import random
import time

# Define the list of 10 IPv4 addresses
ip_addresses = [f"192.168.1.{i}" for i in range(1, 101)]


# Function to send the POST request
def send_post_request(url):
    # Choose a random IP address from the list
    ip = random.choice(ip_addresses)

    # Create the payload
    payload = {"payload": {"ip": ip, "duration": 15}}

    # Send the POST request
    response = requests.post(url, json=payload)

    # Print the response
    print(f"Status Code: {response.status_code}")
    print(f"Response Body: {response.text}")


# URL to send the request to
url = "http://localhost:5000/healthcheck"

while True:
    send_post_request(url)
    time.sleep(15)
