import requests
import struct
import json
import sys
import time

div = ((2 ** 64) - 1) / 100
n = 100
coordinator_put_link = sys.argv[1]
coordinator_get_link = sys.argv[2]

for i in range(0, n):
  form_wrapper = {"key": struct.pack("L", div * i), "value" : str(i)}
  print form_wrapper
  request_object = requests.post(coordinator_put_link, data=form_wrapper)
  time.sleep(.01)

#check the stored values
for i in range(0, n):
  request_object = requests.get(coordinator_put_link, params = {"key" : struct.pack("L", div * i)})
  value_object = json.loads(request_object.text)
  assert value_object["value"] == str(i)
