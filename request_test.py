import requests
import struct
import json
import sys
import time

div = ((2 ** 64) - 1) / 100
n = int( sys.argv[3] )
coordinator_put_link = sys.argv[1]
coordinator_get_link = sys.argv[2]

shortestRequestTime= 9999999999.0
longestRequestTime= 0.0
averageRequestTime= 0.0
for i in range(0, n):
  #form_wrapper = {"key": struct.pack("L", div * i), "value" : str(i)}
  form_wrapper = {"key": str(div*i) , "value" : str(i)}
  #print form_wrapper
  start= time.clock()
  request_object = requests.get(coordinator_put_link, params=form_wrapper)
  end= time.clock()
  duration= end - start
  if duration < shortestRequestTime:
      shortestRequestTime= duration
  if duration > longestRequestTime:
      longestRequestTime= duration
  averageRequestTime+= duration
  
  #print request_object.url
  time.sleep(.01)
  
averageRequestTime= averageRequestTime / n


#check the stored values
oldValues= {}
numErrors= 0
shortestGetTime= 9999999999.0
longestGetTime= 0.0
averageGetTime= 0.0
for i in range(0, n):
  #request_object = requests.get(coordinator_put_link, params = {"key" : struct.pack("L", div * i)})
  start= time.clock()
  request_object = requests.get(coordinator_get_link, params = {"key" : str(div*i) })
  end= time.clock()
  duration= end - start
  if duration < shortestGetTime:
      shortestGetTime= duration
  if duration > longestGetTime:
      longestGetTime= duration
  averageGetTime+= duration
  
  value_object = json.loads(request_object.text)
  try:
    assert value_object["value"] == str(i)
    oldValues[str(div*i)]= str(i)
  except AssertionError:
      numErrors+=1
      print "expected value: " + str(i) 
      print "obtained value: " + value_object["value"]
      
averageGetTime= averageGetTime / n

#now, update values
for key in oldValues:
    form_wrapper= { "key": key, "value": oldValues[key] + "||" + key }
    request_object= requests.get(coordinator_put_link, params=form_wrapper )
    time.sleep(.01)

#check if new values returned
numNewErrors= 0
for key in oldValues:
    request_object= requests.get(coordinator_get_link, params= {"key" : key} )
    value_object= json.loads(request_object.text) 
    try:
        assert value_object["value"] == oldValues[key] + "||" + key
    except AssertionError:
        numNewErrors+=1
        print "expected new value: " + oldValues[key] + "||" + key
        print "obtained value: " + value_object["value"]
        
#print statistics
print ("shortest put request time: %.2f ms" % (1000 * shortestRequestTime)) #+ str(1000 * shortestRequestTime)
print ("longest put request time: %.2f ms" % (1000 * longestRequestTime) )
print ("average put request time: %.2f ms" % (1000 * averageRequestTime) )
print (" ")
print ("shortest get request time: %.2f ms" % (1000 * shortestGetTime) )
print ("longest get request time: %.2f ms" % (1000 * longestGetTime) )
print ("average get request time: %.2f ms" % (1000 * averageGetTime) )
print (" ")
print "assertion errors on initial gets: " + str(100 * numErrors / float(n) ) + "%"
print "assertion errors on updates: " + str(100 * numNewErrors / float(n) ) + "%"
