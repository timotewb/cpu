import requests
import time

locations = list()

def write_list_to_file(filename, lst):
    with open(filename, 'w') as f:
        for item in lst:
            f.write(str(item) + '\n')

for i in range(1,3000):
    url:str = f"https://www.gassy.co.nz/api/marker/{i}"
    print(url)

    payload = {}
    headers = {}

    response = requests.request("GET", url, headers=headers, data=payload)
    if response.status_code == 200:
        print(" - found")
        locations.append(i)
    time.sleep(1)


write_list_to_file("gaspy_locations.txt", locations)