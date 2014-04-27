import requests
import struct
import json

div = ((2 ** 64) - 1) / 100
n = 100
coordinator_put_link = ""
coordinator_get_link = ""

for i in range(0, n):
	form_wrapper = {"key": struct.pack("I", div * i), "value" : str(i)}
	request_object = requests.post(coordinator_put_link, data=form_wrapper)

#check the stored values
for i in range(0, n):
	request_object = requests.get(coordinator_put_link, params = {"key" : struct.pack("I", div * i)})
	value_object = json.loads(request_object.text)
	assert value_object["value"] == str(i)
