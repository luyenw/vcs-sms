import pandas as pd
import random
import string
import ipaddress
def generate_random_name(length=8):
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for i in range(length))
def generate_random_ipv4():
    return str(ipaddress.IPv4Address(random.randint(0, 2**32 - 1)))

N = 10000
data = {
    'name': [generate_random_name() for _ in range(N)],
    'ipv4': [generate_random_ipv4() for _ in range(N)]
}
df = pd.DataFrame(data)
df.to_excel(f'output_{N}.xlsx', index=False)
