import requests
import struct
import json
import sys
import time

div = ((2 ** 64) - 1) / 100
n = 10000
coordinator_put_link = sys.argv[1]
coordinator_get_link = sys.argv[2]

shortestRequestTime= 9999999999.0
longestRequestTime= 0.0
averageRequestTime= 0.0
for i in range(0, n):
  #form_wrapper = {"key": struct.pack("L", div * i), "value" : str(i)}
  form_wrapper = {"key": str(div*i) , "value" : str(i)}
  print form_wrapper
  start= time.clock()
  request_object = requests.get(coordinator_put_link, params=form_wrapper)
  end= time.clock()
  duration= end - start
  if duration < shortestRequestTime:
      shortestRequestTime= duration
  if duration > longestRequestTime:
      longestRequestTime= duration
  averageRequestTime+= duration
  
  print request_object.url
  time.sleep(.01)
  
averageRequestTime= averageRequestTime / n
print "shortest request time: " + str(shortestRequestTime)
print "longest request time: " + str(longestRequestTime)
print "average request time: " + str(averageRequestTime)

#check the stored values
oldValues= {}
for i in range(0, n):
  #request_object = requests.get(coordinator_put_link, params = {"key" : struct.pack("L", div * i)})
  request_object = requests.get(coordinator_get_link, params = {"key" : str(div*i) })
  value_object = json.loads(request_object.text)
  assert value_object["value"] == str(i)
  oldValues[str(div*i)]= str(i)

#now, update values
for key in oldValues:
    form_wrapper= { "key": key, "value": oldValues[key] + "||" + key }
    request_object= requests.get(coordinator_put_link, params=form_wrapper )
    time.sleep(.01)

#check if new values returned
for key in oldValues:
    request_object= requests.get(coordinator_get_link, params= {"key" : key} )
    value_object= json.loads(request_object.text)
    assert value_object["value"] == oldValues[key] + "||" + key
